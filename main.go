package main

import (
	"avito_tech/internal/app"
	"avito_tech/pkg/config"
	"avito_tech/pkg/database"
	"avito_tech/pkg/logger"
	"database/sql"
	"time"
)

func main() {
	logger.Log.Infoln("Let's wait a bit")
	time.Sleep(10 * time.Second)
	cfg := config.GetConfig()
	db := database.Init(cfg)

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Log.Errorln("Error closing database: " + err.Error())
		}
	}(db)

	app.StartApp(cfg, db)
}
