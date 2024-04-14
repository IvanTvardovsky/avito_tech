package updater

import (
	"avito_tech/internal/structures"
	"avito_tech/pkg/logger"
	"database/sql"
	"encoding/json"
	"time"
)

func UpdateInfoOnServer(matrix *map[int]map[int]int, bannerMap *map[int]*structures.Banner, db *sql.DB, cfg *structures.Config) error {
	query := "SELECT tag_id, feature_id, banner_id, content, is_active, created_at, updated_at FROM "
	if cfg.AppMode.IsTest {
		query += "test_bannermatrix"
	} else {
		query += "bannermatrix"
	}
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Log.Error("Error closing SQL rows:" + err.Error())
			return
		}
	}(rows)

	bannerFlag := make(map[int]bool) // для append

	for rows.Next() {
		var tagID, featureID, bannerID int
		var contentJSON string
		var isActive bool
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&tagID, &featureID, &bannerID, &contentJSON, &isActive, &createdAt, &updatedAt); err != nil {
			return err
		}

		if (*matrix)[tagID] == nil {
			(*matrix)[tagID] = make(map[int]int)
		}
		(*matrix)[tagID][featureID] = bannerID

		if _, ok := (*bannerMap)[bannerID]; !ok {
			(*bannerMap)[bannerID] = &structures.Banner{}
			bannerFlag[bannerID] = true
		} else if !bannerFlag[bannerID] {
			(*bannerMap)[bannerID].TagIDs = nil
			bannerFlag[bannerID] = true
		}

		(*bannerMap)[bannerID].ID = bannerID
		(*bannerMap)[bannerID].TagIDs = append((*bannerMap)[bannerID].TagIDs, tagID)
		(*bannerMap)[bannerID].FeatureID = featureID
		(*bannerMap)[bannerID].IsActive = isActive
		(*bannerMap)[bannerID].CreatedAt = createdAt.String()
		(*bannerMap)[bannerID].UpdatedAt = updatedAt.String()

		var content map[string]interface{}
		if err := json.Unmarshal([]byte(contentJSON), &content); err != nil {
			return err
		}
		(*bannerMap)[bannerID].Content = content
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
