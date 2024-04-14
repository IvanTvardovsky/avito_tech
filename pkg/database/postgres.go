package database

import (
	"avito_tech/internal/structures"
	"avito_tech/pkg/logger"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func Init(cfg *structures.Config) *sql.DB {
	logger.Log.Infoln("Connecting to database...")
	logger.Log.Traceln(fmt.Sprintf("Connecting to host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Username, cfg.Storage.Password, cfg.Storage.Database))
	psqlconn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.Storage.Username, cfg.Storage.Password, cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Database)

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
