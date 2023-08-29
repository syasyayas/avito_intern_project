package postgres

import (
	"avito_project/internal/model"
	"avito_project/internal/repository/repoerr"
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

	rows, err := r.db.Query(ctx, "SELECT user_id, feature_slug, operation, date FROM avito_features.history WHERE date >= $1 AND date <= $2", after, before)
	if err != nil {
		return nil, repoerr.PgErrorWrapper(err)
	}

	var res model.History

	for rows.Next() {
		record := model.Record{}
		err := rows.Scan(&record.UserID, &record.FeatureSlug, &record.Operation, &record.Date)
		if err != nil {
			return nil, repoerr.PgErrorWrapper(err)
		}
		res = append(res, record)
	}
	return res, nil
}
