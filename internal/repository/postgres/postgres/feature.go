package postgres

import (
	"avito_project/internal/model"
	"avito_project/internal/repository/repoerr"
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

func (r *FeatureRepo) AddFeature(ctx context.Context, slug string) error {
	r.log.Debugf("Adding feature %s", slug)

	_, err := r.db.Exec(ctx, "INSERT INTO avito_features.features (slug) VALUES ($1)", slug)

	return repoerr.PgErrorWrapper(err)

}

func (r *FeatureRepo) DeleteFeature(ctx context.Context, slug string) error {
	r.log.Debugf("Deleting feature %s", slug)

	_, err := r.db.Exec(ctx, "DELETE FROM avito_features.features WHERE slug = $1", slug)

	return repoerr.PgErrorWrapper(err)
}

// AddFeaturesToUser inserts feature to user relations
// execution fails if provided feature does not exist or the expiration date is invalid
func (r *FeatureRepo) AddFeaturesToUser(ctx context.Context, userId string, features []model.Feature) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return repoerr.PgErrorWrapper(err)
	}
	defer tx.Rollback(ctx)

	for _, feature := range features {
		r.log.Debugf("adding feature %v to user %s", feature, userId)
		if err != nil {
			return repoerr.PgErrorWrapper(err)
		}
		if !feature.ExpiresAt.IsZero() {

			if feature.ExpiresAt.Before(time.Now()) {
				return repoerr.ErrInvalidExpiresAt
			}

			_, err := tx.Exec(ctx,
				`INSERT INTO avito_features.user_feature (user_id, feature_slug, expires_at) 
					 VALUES ($1, $2, $3)`,
				userId, feature.Slug, feature.ExpiresAt)
			if err != nil {
				return repoerr.PgErrorWrapper(err)
			}
		} else {
			_, err := tx.Exec(ctx,
				`INSERT INTO avito_features.user_feature (user_id, feature_slug) 
					 VALUES ($1, $2)`,
				userId, feature.Slug)
			if err != nil {
				return repoerr.PgErrorWrapper(err)
			}
		}
	}

	return repoerr.PgErrorWrapper(tx.Commit(ctx))
}

// DeleteFeaturesFromUser deletes feature to user relation
// execution fails if one of features does not exist but does not if user-feature relation does not exist
func (r *FeatureRepo) DeleteFeaturesFromUser(ctx context.Context, userId string, features []model.Feature) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return repoerr.PgErrorWrapper(err)
	}
	defer tx.Rollback(ctx)

	for _, feature := range features {

		_, err = tx.Exec(ctx, "DELETE FROM avito_features.user_feature WHERE user_id = $1 AND feature_slug = $2", userId, feature.Slug)

		if errors.Is(repoerr.PgErrorWrapper(err), repoerr.ErrNotFound) {
			continue
		}
	}

	return repoerr.PgErrorWrapper(tx.Commit(ctx))
}

func (r *FeatureRepo) DeleteFeatureFromUser(ctx context.Context, userID string, featureSlug string) error {
	r.log.Debugf("Deleting relation userID: %s, featureSLug: %s", userID, featureSlug)

	_, err := r.db.Exec(ctx, "DELETE FROM avito_features.user_feature WHERE user_id = $1 AND feature_slug = $2", userID, featureSlug)
	return repoerr.PgErrorWrapper(err)
}

func (r *FeatureRepo) GetUserWithFeatures(ctx context.Context, userId string) (*model.User, error) {
	r.log.Debugf("Retrieving info for user %s", userId)

	rows, err := r.db.Query(ctx, `SELECT uf.feature_slug, uf.expires_at FROM avito_features.user_feature uf WHERE uf.user_id = $1`,
		userId)
	if err != nil {
		return nil, repoerr.PgErrorWrapper(err)
	}
	defer rows.Close()

	var user = &model.User{ID: userId}

	for rows.Next() {
		var slug string
		var expiresAtPtr *time.Time
		err = rows.Scan(&slug, &expiresAtPtr)
		if err != nil {
			return nil, err
		}
		feature := model.Feature{
			Slug: slug,
		}
		var expiresAt time.Time
		if expiresAtPtr != nil {
			expiresAt = *expiresAtPtr
		}

		if !expiresAt.IsZero() {
			if expiresAt.Before(time.Now()) {
				_ = r.DeleteFeatureFromUser(ctx, userId, slug)
				continue
			}
			feature.ExpiresAt = expiresAt
		}
		user.Features = append(user.Features, feature)

	}
	return user, nil
}
