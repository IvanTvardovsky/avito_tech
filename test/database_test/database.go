package database_test

import (
	"avito_tech/internal/structures"
	"avito_tech/pkg/logger"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
)

func initTestDB(cfg *structures.Config) *sql.DB {
	logger.Log.Infoln("Connecting to test database...")
	logger.Log.Traceln(fmt.Sprintf("Connecting to host=%s port=%d user=%s dbname=%s",
		cfg.TestStorage.Host, cfg.TestStorage.Port, cfg.TestStorage.Username, cfg.TestStorage.Database))
	psqlconn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.TestStorage.Username, cfg.TestStorage.Password, cfg.TestStorage.Host, cfg.TestStorage.Port, cfg.TestStorage.Database)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		logger.Log.Fatalln("Can not connect to database: " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		logger.Log.Fatalln("Error pinging database: " + err.Error())
	}
	logger.Log.Infoln("Connected to database")
	return db
}

func InitAndFillTestDB(cfg *structures.Config) *sql.DB {
	db := initTestDB(cfg)

	dbInitDataCases := []structures.Banner{
		{
			ID:        0,
			TagIDs:    []int{1, 2, 3},
			FeatureID: 111,
			Content: map[string]interface{}{
				"title": "hello im test1",
				"url":   "google.com",
			},
			IsActive: true,
		},
		{
			ID:        1,
			TagIDs:    []int{2, 4, 5},
			FeatureID: 222,
			Content: map[string]interface{}{
				"title": "hello im inactive test2",
				"url":   "google2.com",
			},
			IsActive: false,
		},
		{
			ID:        2,
			TagIDs:    []int{5, 42, 142},
			FeatureID: 42,
			Content: map[string]interface{}{
				"title": ".",
				"url":   "#?@#@?#@",
			},
			IsActive: true,
		},
		{
			ID:        3,
			TagIDs:    []int{5, 42, 142},
			FeatureID: 43,
			Content: map[string]interface{}{
				"title": "./////.",
				"url":   "#?@#@?#@",
			},
			IsActive: true,
		},
	}

	query := `INSERT INTO test_bannermatrix (tag_id, feature_id, banner_id, content, created_at, updated_at, is_active)
			  SELECT
					unnest($1::integer[]) AS tag_id,
					$2 AS feature_id,
					$3 AS banner_id,
					$4::json AS content,
					CURRENT_TIMESTAMP AS created_at,
					CURRENT_TIMESTAMP AS updated_at,
					$5 AS is_active; `

	for _, banner := range dbInitDataCases {
		tagIDsArray := pq.Array(banner.TagIDs)
		bannerContentJSON, err := json.Marshal(banner.Content)
		if err != nil {
			logger.Log.Fatal(err)
		}

		_, err = db.Exec(query, tagIDsArray, banner.FeatureID, banner.ID, bannerContentJSON, banner.IsActive)
		if err != nil {
			logger.Log.Fatal(err)
		}
	}

	return db
}

func ClearDatabase(db *sql.DB) {
	_, err := db.Exec("DELETE FROM test_bannermatrix")
	if err != nil {
		logger.Log.Fatal(err)
	}
}
