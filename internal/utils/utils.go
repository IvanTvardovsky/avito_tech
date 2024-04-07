package utils

import (
	"database/sql"
)

func GetLastBannerIndex(db *sql.DB) (int, error) {
	query := `SELECT banner_id FROM bannermatrix ORDER BY banner_id DESC LIMIT 1`
	rows, err := db.Query(query)
	if err != nil {
		return 0, err
	}
	var ind int
	for rows.Next() {
		err = rows.Scan(&ind)
		if err != nil {
			return 0, err
		}
	}

	return ind, nil
}
