package mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "mongodb://localhost:27017", cfg.URI)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 27017, cfg.Port)
	assert.Equal(t, "gin_admin", cfg.Database)
	assert.Equal(t, "", cfg.Username)
	assert.Equal(t, "", cfg.Password)
	assert.Equal(t, "admin", cfg.AuthSource)
	assert.Equal(t, uint64(100), cfg.MaxPoolSize)
	assert.Equal(t, uint64(10), cfg.MinPoolSize)
	assert.Equal(t, time.Minute*5, cfg.MaxConnIdle)
	assert.Equal(t, time.Second*10, cfg.ConnectTimeout)
	assert.Equal(t, time.Second*30, cfg.ServerTimeout)
	assert.Equal(t, time.Second*30, cfg.Timeout)
	assert.Equal(t, 6, cfg.CompressLevel)
	assert.Equal(t, "", cfg.ReplicaSet)
}

func TestConfigGetURI(t *testing.T) {
	t.Run("Default URI", func(t *testing.T) {
		cfg := &Config{
			URI: "mongodb://localhost:27017",
		}
		assert.Equal(t, "mongodb://localhost:27017", cfg.GetURI())
	})

	t.Run("Custom URI", func(t *testing.T) {
		cfg := &Config{
			URI: "mongodb://user:pass@cluster.mongodb.net/mydb",
		}
		assert.Equal(t, "mongodb://user:pass@cluster.mongodb.net/mydb", cfg.GetURI())
	})

	t.Run("Built URI with credentials", func(t *testing.T) {
		cfg := &Config{
			URI:        "mongodb://localhost:27017",
			Host:       "localhost",
			Port:       27017,
			Username:   "admin",
			Password:   "password",
			AuthSource: "mydb",
		}
		uri := cfg.GetURI()
		assert.Contains(t, uri, "admin:password@")
		assert.Contains(t, uri, "localhost:27017")
		assert.Contains(t, uri, "/mydb")
	})

	t.Run("Built URI without credentials", func(t *testing.T) {
		cfg := &Config{
			URI:  "mongodb://localhost:27017",
			Host: "localhost",
			Port: 27017,
		}
		uri := cfg.GetURI()
		assert.Equal(t, "mongodb://localhost:27017", uri)
	})
}
