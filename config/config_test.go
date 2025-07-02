package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func resetConfig() {
	cfg = nil
	once = sync.Once{}
	viper.Reset()
}

func TestLoadConfig_WithEnvVars(t *testing.T) {
	t.Setenv("PORT", "9999")
	t.Setenv("DB_URL", "postgres://localhost:5432/testdb")

	resetConfig()

	c, err := LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "9999", c.AppPort)
	assert.Equal(t, "postgres://localhost:5432/testdb", c.DBURL)
}

func TestLoadConfig_FallbackPort(t *testing.T) {
	t.Setenv("DB_URL", "postgres://localhost:5432/testdb")

	resetConfig()
	c, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "8080", c.AppPort)
	assert.Equal(t, "postgres://localhost:5432/testdb", c.DBURL)
}

func TestLoadConfig_MissingDBURL(t *testing.T) {
	resetConfig()

	c, err := LoadConfig()
	assert.Error(t, err)
	assert.Nil(t, c)

	t.Logf("error:%v", err)
}
