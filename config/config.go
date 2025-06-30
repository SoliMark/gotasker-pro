package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppPort string
	DBURL   string
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() (*Config, error) {
	var loadErr error

	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using system env only")
		}
		viper.AutomaticEnv()
		// port fallback
		viper.SetDefault("PORT", "8080")

		dbURL := viper.GetString("DB_URL")
		if dbURL == "" {
			loadErr = fmt.Errorf("missing require DB_URL")
			return
		}

		cfg = &Config{
			AppPort: viper.GetString("PORT"),
			DBURL:   dbURL,
		}
	})

	if loadErr != nil {
		return nil, loadErr
	}
	return cfg, nil
}
