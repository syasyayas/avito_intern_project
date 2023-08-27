package model

import "time"

type Record struct {
	UserID      string    `json:"user_id" csv:"user_id"`
	FeatureSlug string    `json:"feature_slug" csv:"feature_id"`
	Operation   string    `json:"operation" csv:"operation"`
	Date        time.Time `json:"date" csv:"date"`
}

type History []Record

func (h History) ParseToCSV() [][]string {
	var res [][]string
	headers := []string{"user_id", "feature_slug", "operation", "date"}
	res = append(res, headers)
	for _, record := range h {
		res = append(res, record.String())
	}
	return res
}

func (r Record) String() []string {
	return []string{
		r.UserID,
		r.FeatureSlug,
		r.Operation,
		r.Date.String(),
	}
}
