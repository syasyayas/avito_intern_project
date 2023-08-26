package postgres

import (
	"avito_project/internal/model"
	"avito_project/internal/repository"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

type FeatureRepo struct {
	db *pgxpool.Pool

	log *logrus.Logger
}

func NewFeatureRepo(db *pgxpool.Pool, log *logrus.Logger) *FeatureRepo {
	return &FeatureRepo{db, log}
}

func (r *FeatureRepo) AddFeature(ctx context.Context, slug string) (int, error) {
	r.log.Debugf("Adding feature %s", slug)

	var id int
	row := r.db.QueryRow(ctx, "INSERT INTO features (slug) VALUES ($1) RETURNING id", slug)
	err := row.Scan(&id)

	return id, repository.PgErrorWrapper(err)

}

func (r *FeatureRepo) DeleteFeature(ctx context.Context, slug string) error {
	r.log.Debugf("Deleting feature %s", slug)

	_, err := r.db.Exec(ctx, "UPDATE features SET deleted_at = now() WHERE slug = $1", slug)

	return repository.PgErrorWrapper(err)
}

func (r *FeatureRepo) getIdBySlug(ctx context.Context, slug string) (int, error) {
	var id int

	row := r.db.QueryRow(ctx, "SELECT id form avito_features.features WHERE slug = $1 AND deleted_at IS NULL", slug)
	err := row.Scan(&id)

	return id, repository.PgErrorWrapper(err)
}

// AddFeaturesToUser inserts feature to user relations
// execution fails if provided feature does not exist or the expiration date is invalid
func (r *FeatureRepo) AddFeaturesToUser(ctx context.Context, userId string, features []*model.Feature) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return repository.PgErrorWrapper(err)
	}
	defer tx.Rollback(ctx)

	for _, feature := range features {
		r.log.Debugf("adding feature %v to user %s", feature, userId)
		if !feature.ExpiresAt.IsZero() {

			if feature.ExpiresAt.Before(time.Now()) {
				return repository.ErrInvalidExpiresAt
			}

			_, err := tx.Exec(ctx,
				`INSERT INTO avito_features.user_feature (user_id, feature_id, expires_at) 
					 VALUES ($1, (SELECT id FROM avito_features.features WHERE slug = $2), $3)`,
				userId, feature.Slug, feature.ExpiresAt)
			if err != nil {
				return repository.PgErrorWrapper(err)
			}
		} else {
			_, err := tx.Exec(ctx,
				`INSERT INTO avito_features.user_feature (user_id, feature_id) 
					 VALUES ($1, (SELECT id FROM avito_features.features WHERE slug = $2))`,
				userId, feature.Slug)
			if err != nil {
				return repository.PgErrorWrapper(err)
			}
		}
	}

	return repository.PgErrorWrapper(tx.Commit(ctx))
}

// DeleteFeaturesFromUser deleted feature to user relation
// execution fails if one of features does not exist but does not if user-feature relation does not exist
func (r *FeatureRepo) DeleteFeaturesFromUser(ctx context.Context, userId string, features []*model.Feature) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return repository.PgErrorWrapper(err)
	}
	defer tx.Rollback(ctx)

	for _, feature := range features {
		featureId, err := r.getIdBySlug(ctx, feature.Slug)
		if errors.Is(err, repository.ErrNotFound) {
			return err
		}

		_, err = tx.Exec(ctx, "DELETE FROM avito_features.user_feature WHERE user_id = $1 AND feature_id = $2", userId, featureId)

		if errors.Is(repository.PgErrorWrapper(err), repository.ErrNotFound) {
			continue
		}
	}

	return repository.PgErrorWrapper(tx.Commit(ctx))
}

func (r *FeatureRepo) GetUserWithFeatures(ctx context.Context, id string) (*model.User, error) {
	r.log.Debugf("Retrieving info for user %s", id)

	rows, err := r.db.Query(ctx, "SELECT uf.feature_id, f.slug, uf.expires_at FROM avito_features.user_feature uf JOIN avito_features.features f ON uf.feature_id = f.id WHERE uf.user_id = $1", id)
	if err != nil {
		return nil, repository.PgErrorWrapper(err)
	}
	defer rows.Close()

	var user = &model.User{ID: id}
	for rows.Next() {
		var featureId int
		var slug string
		var expiresAt time.Time

		err = rows.Scan(&featureId, &slug, &expiresAt)
		if err != nil {
			return nil, err
		}

		feature := model.Feature{
			ID:   featureId,
			Slug: slug,
		}

		if !expiresAt.IsZero() {
			if expiresAt.Before(time.Now()) {
				_ = r.DeleteFeatureFromUser(ctx, id, featureId)
				continue
			}
			feature.ExpiresAt = expiresAt
		}
		user.Features = append(user.Features, feature)

	}
	return user, nil
}

func (r *FeatureRepo) GetHistory(ctx context.Context, after time.Time, before time.Time) ([]model.History, error) {
	r.log.Debugf("Getting history between %s and %s", before.String(), after.String())

	rows, err := r.db.Query(ctx, "SELECT user_id, feature_id, opeartion, date FROM avito_features.history WHERE date BETWEEN $1 AND $2", after, before)
	if err != nil {
		return nil, repository.PgErrorWrapper(err)
	}

	var res []model.History

	for rows.Next() {
		hist := model.History{}
		err := rows.Scan(&hist.UserID, &hist.FeatureID, &hist.Operation, &hist.Date)
		if err != nil {
			return nil, repository.PgErrorWrapper(err)
		}
		res = append(res, hist)
	}
	return res, nil
}

func (r *FeatureRepo) DeleteFeatureFromUser(ctx context.Context, userID string, featureID int) error {
	r.log.Debugf("Deleting relation userID: %s, featureID: %d", userID, featureID)

	_, err := r.db.Exec(ctx, "DELETE FROM avito_features.user_feature WHERE user_id = $1 AND feature_id = $2", userID, featureID)
	return repository.PgErrorWrapper(err)
}
