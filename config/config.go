package config

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	// App / DB / Auth
	AppPort   string `mapstructure:"PORT"`       // default: 8080
	DBURL     string `mapstructure:"DB_URL"`     // required
	JWTSecret string `mapstructure:"JWT_SECRET"` // required

	// Redis cache for task list
	RedisAddr     string        `mapstructure:"REDIS_ADDR"`      // default: localhost:6379
	RedisPassword string        `mapstructure:"REDIS_PASSWORD"`  // default: ""
	RedisDB       int           `mapstructure:"REDIS_DB"`        // default: 0
	CacheTTLTasks time.Duration `mapstructure:"CACHE_TTL_TASKS"` // default: 60s (supports "1m", "45s", etc.)
}

var (
	cfg  *Config
	once sync.Once
)

// LoadConfig loads configuration from environment variables (and .env if present).
// It uses mapstructure to unmarshal env values into the Config struct, including time.Duration.
func LoadConfig() (*Config, error) {
	var initErr error

	once.Do(func() {
		// Load .env into process env if present (non-fatal).
		if err := godotenv.Load(); err != nil {
			log.Println("config: no .env file found, proceeding with system env only")
		}

		v := viper.New()
		v.AutomaticEnv() // read from process env
		// -----------------------------------------------------------------------------
		// Viper value priority (highest → lowest):
		// 1. Explicit Set() values (v.Set)
		// 2. Bound CLI flags (BindPFlag)
		// 3. Environment variables (BindEnv / AutomaticEnv)
		// 4. Config file values (ReadInConfig)
		// 5. Remote Key/Value store values (etcd/consul)
		// 6. Default values (SetDefault)
		// -----------------------------------------------------------------------------

		// Defaults — safe fallbacks for local/dev.
		v.SetDefault("PORT", "8080")

		// Redis cache for task list
		v.SetDefault("REDIS_ADDR", "localhost:6379")
		v.SetDefault("REDIS_PASSWORD", "")
		v.SetDefault("REDIS_DB", 0)
		v.SetDefault("CACHE_TTL_TASKS", "60s")

		_ = v.BindEnv("PORT")
		_ = v.BindEnv("DB_URL")
		_ = v.BindEnv("JWT_SECRET")

		// Redis cache for task list
		_ = v.BindEnv("REDIS_ADDR")
		_ = v.BindEnv("REDIS_PASSWORD")
		_ = v.BindEnv("REDIS_DB")
		_ = v.BindEnv("CACHE_TTL_TASKS")
		var c Config
		// Enable time.Duration decoding from strings like "60s", "1m".
		if err := v.Unmarshal(&c, viper.DecodeHook(
			mapstructure.StringToTimeDurationHookFunc(),
		)); err != nil {
			initErr = err
			return
		}

		// Basic validation for required fields.
		if c.DBURL == "" {
			initErr = errors.New("config: DB_URL is required")
			return
		}
		if c.JWTSecret == "" {
			initErr = errors.New("config: JWT_SECRET is required")
			return
		}

		cfg = &c
	})

	if initErr != nil {
		return nil, initErr
	}
	return cfg, nil
}

// Optional helpers (nice-to-have)
func (c *Config) RedisEnabled() bool { return c != nil && c.RedisAddr != "" }
