package app

import (
	"avito_tech/internal/handlers"
	"avito_tech/internal/middleware"
	"avito_tech/internal/structures"
	"avito_tech/internal/updater"
	"avito_tech/pkg/logger"
	"database/sql"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

func StartApp(cfg *structures.Config, db *sql.DB) {
	matrix := make(map[int]map[int]int)
	bannerMap := make(map[int]*structures.Banner)
	ind, err := getLastBannerIndex(db)
	if err != nil {
		logger.Log.Error("Error getting last banner index: ", err.Error())
		return
	}
	handlers.BannerCounter.Store(int64(ind) + 1)

	go func() {
		for {
			err = updater.UpdateInfoOnServer(&matrix, &bannerMap, db, cfg)
			if err != nil {
				logger.Log.Error("Error updating info on server: " + err.Error())
			}

			if !cfg.AppMode.IsTest {
				// это для удобства тестирования, по условиям задания можно заменить на 5 * time.Minute
				time.Sleep(5 * time.Second)
			}
		}
	}()

	logger.Log.Infoln("Starting service...")
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	router.Use(middleware.TokenTypeMiddleware())

	router.GET("/user_banner", func(c *gin.Context) {
		handlers.GetUserBanner(c, db, cfg, &matrix, &bannerMap)
	})
	router.GET("/banner", middleware.AuthAdminMiddleware(), func(c *gin.Context) {
		handlers.GetAllBanners(c, &matrix, &bannerMap)
	})
	router.POST("/banner", middleware.AuthAdminMiddleware(), func(c *gin.Context) {
		handlers.CreateBanner(c, db)
	})
	router.PATCH("/banner/:id", middleware.AuthAdminMiddleware(), func(c *gin.Context) {
		handlers.PatchBanner(c, db)
	})
	router.DELETE("/banner/:id", middleware.AuthAdminMiddleware(), func(c *gin.Context) {
		handlers.DeleteBanner(c, db)
	})
	router.GET("/banners/:id/versions", middleware.AuthAdminMiddleware(), func(c *gin.Context) {
		handlers.GetBannerVersions(c, db)
	})

	logger.Log.Infoln("Serving handlers...")
	logger.Log.Infoln("Starting router...")
	logger.Log.Infoln("On port :" + cfg.Listen.Port)
	logger.Log.Fatal(router.Run(":" + cfg.Listen.Port))
}

func getLastBannerIndex(db *sql.DB) (int, error) {
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
