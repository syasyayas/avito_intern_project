package postgres

import (
	"avito_project/internal/model"
	"avito_project/internal/repository"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

type HistoryRepo struct {
	db *pgxpool.Pool

	log *logrus.Logger
}

func NewHistoryRepo(db *pgxpool.Pool, log *logrus.Logger) *HistoryRepo {
	return &HistoryRepo{
		db,
		log,
	}
}

func (r *HistoryRepo) GetHistory(ctx context.Context, after time.Time, before time.Time) (model.History, error) {
	r.log.Debugf("Getting history between %s and %s", before.String(), after.String())

	rows, err := r.db.Query(ctx, "SELECT user_id, (SELECT slug FROM avito_features.features WHERE id = history.feature_id), opeartion, date FROM avito_features.history WHERE date BETWEEN $1 AND $2", after, before)
	if err != nil {
		return nil, repository.PgErrorWrapper(err)
	}

	var res model.History

	for rows.Next() {
		hist := model.Record{}
		err := rows.Scan(&hist.UserID, &hist.FeatureSlug, &hist.Operation, &hist.Date)
		if err != nil {
			return nil, repository.PgErrorWrapper(err)
		}
		res = append(res, hist)
	}
	return res, nil
}
