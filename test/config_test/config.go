package config_test

import (
	"avito_tech/internal/structures"
	"avito_tech/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

var cfg *structures.Config
var once sync.Once

func GetTestConfig() *structures.Config {
	once.Do(func() {
		logger.Log.Infoln("Reading app configuration...")
		cfg = &structures.Config{}
		if err := cleanenv.ReadConfig("./../config.yml", cfg); err != nil {
			logger.Log.Errorln(err)
			help, err := cleanenv.GetDescription(cfg, nil)
			logger.Log.Errorln(help)
			logger.Log.Fatalln(err)
		}
	})
	return cfg
}
