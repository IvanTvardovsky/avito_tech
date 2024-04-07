package main

import (
	"avito_tech/internal/handlers"
	"avito_tech/internal/middleware"
	"avito_tech/internal/structures"
	"avito_tech/internal/utils"
	"avito_tech/pkg/config"
	"avito_tech/pkg/database"
	"avito_tech/pkg/logger"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

func updateInfoOnServer(matrix *map[int]map[int]int, bannerMap *map[int]*structures.Banner, db *sql.DB) error {
	rows, err := db.Query("SELECT tag_id, feature_id, banner_id, content, is_active, created_at, updated_at FROM BannerMatrix")
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

func hi(c *gin.Context) {
	c.Writer.Write([]byte("Hi"))
}

func printMatrixAndBanners(matrix map[int]map[int]int, bannerMap map[int]*structures.Banner) {
	for tagID, features := range matrix {
		for featureID, bannerID := range features {
			banner, ok := bannerMap[bannerID]
			if !ok {
				continue
			}

			fmt.Printf("Tag ID: %d, Feature ID: %d, Banner ID: %d\n", tagID, featureID, bannerID)
			fmt.Printf("Banner Info: %+v\n", banner)
		}
	}
}

func main() {
	cfg := config.GetConfig()
	db := database.Init(cfg)

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Log.Errorln("Error closing database: " + err.Error())
		}
	}(db)

	matrix := make(map[int]map[int]int)
	bannerMap := make(map[int]*structures.Banner)
	ind, err := utils.GetLastBannerIndex(db)
	if err != nil {
		logger.Log.Error("Error getting last banner index: ", err.Error())
		return
	}
	handlers.BannerCounter.Store(int64(ind) + 1)

	go func() {
		for {
			err = updateInfoOnServer(&matrix, &bannerMap, db)
			if err != nil {
				logger.Log.Error("Error updating info on server: " + err.Error())
			}
			//todo увеличить время + сделать канал/контекст для завершения?
			time.Sleep(10 * time.Second)
		}
	}()

	logger.Log.Infoln("Starting service...")
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	router.GET("/", hi)
	router.GET("/user_banner")                                                 // 200 400 401 403 404 500
	router.GET("/banner", middleware.AuthMiddleware(), handlers.GetAllBanners) // 200 401 403 500
	router.POST("/banner", middleware.AuthMiddleware(), func(c *gin.Context) {
		handlers.CreateBanner(c, db)

	}) // 201 400 401 403 500
	router.PATCH("/banner/{id}")  // 200 400 401 403 404 500
	router.DELETE("/banner/{id}") // 204 400 401 403 404 500

	logger.Log.Infoln("Serving handlers...")
	logger.Log.Info("Starting router...")
	logger.Log.Info("On port :" + cfg.Listen.Port)
	logger.Log.Fatal(router.Run(":" + cfg.Listen.Port))
}
