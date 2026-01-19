package postgresql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 5432, cfg.Port)
	assert.Equal(t, "gin_admin", cfg.Database)
	assert.Equal(t, "postgres", cfg.Username)
	assert.Equal(t, "password", cfg.Password)
	assert.Equal(t, "disable", cfg.SSLMode)
	assert.Equal(t, "Asia/Shanghai", cfg.Timezone)
	assert.Equal(t, 10, cfg.MaxIdleConns)
	assert.Equal(t, 100, cfg.MaxOpenConns)
	assert.Equal(t, time.Hour, cfg.MaxLifetime)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, 200*time.Millisecond, cfg.SlowThreshold)
	assert.Len(t, cfg.Extensions, 4)
}

func TestExtensionConfig(t *testing.T) {
	extensions := []ExtensionConfig{
		{Name: "postgis", Version: "3.3", Enabled: true},
		{Name: "vector", Version: "0.5.1", Enabled: true},
	}

	assert.Equal(t, "postgis", extensions[0].Name)
	assert.Equal(t, "3.3", extensions[0].Version)
	assert.True(t, extensions[0].Enabled)
}

func TestVectorOperations(t *testing.T) {
	t.Run("Vector creation", func(t *testing.T) {
		vec := Vector{1.0, 2.0, 3.0}
		assert.Equal(t, 3, vec.Size())
		assert.Contains(t, vec.String(), "1.000000")
		assert.Contains(t, vec.String(), "2.000000")
		assert.Contains(t, vec.String(), "3.000000")
	})

	t.Run("Vector normalization", func(t *testing.T) {
		vec := Vector{3.0, 4.0, 0.0}
		normalized := vec.Normalize()
		assert.InDelta(t, 0.6, normalized[0], 0.001)
		assert.InDelta(t, 0.8, normalized[1], 0.001)
		assert.InDelta(t, 0.0, normalized[2], 0.001)
	})

	t.Run("Cosine similarity", func(t *testing.T) {
		vec1 := Vector{1.0, 0.0, 0.0}
		vec2 := Vector{1.0, 0.0, 0.0}
		similarity, err := CosineSimilarity(vec1, vec2)
		assert.NoError(t, err)
		assert.Equal(t, float32(1.0), similarity)
	})

	t.Run("Euclidean distance", func(t *testing.T) {
		vec1 := Vector{0.0, 0.0}
		vec2 := Vector{3.0, 4.0}
		distance, err := EuclideanDistance(vec1, vec2)
		assert.NoError(t, err)
		assert.InDelta(t, 5.0, distance, 0.001)
	})
}

func TestGeographicOperations(t *testing.T) {
	t.Run("Distance between points", func(t *testing.T) {
		// 北京到上海大约 1068 公里
		beijing := Point{39.9042, 116.4074}
		shanghai := Point{31.2304, 121.4737}

		distance := DistanceBetween(beijing.Latitude, beijing.Longitude, shanghai.Latitude, shanghai.Longitude)
		assert.InDelta(t, 1068000, distance, 100000) // 允许100km误差
	})

	t.Run("Bounding box creation", func(t *testing.T) {
		minLat, minLng, maxLat, maxLng := CreateBoundingBox(39.9042, 116.4074, 10.0)

		assert.True(t, minLat < 39.9042)
		assert.True(t, maxLat > 39.9042)
		assert.True(t, minLng < 116.4074)
		assert.True(t, maxLng > 116.4074)
	})

	t.Run("Point in polygon", func(t *testing.T) {
		point := Point{1.0, 1.0}
		polygon := []Point{
			{0.0, 0.0},
			{2.0, 0.0},
			{2.0, 2.0},
			{0.0, 2.0},
		}

		assert.True(t, PointInPolygon(point, polygon))

		outsidePoint := Point{3.0, 3.0}
		assert.False(t, PointInPolygon(outsidePoint, polygon))
	})
}

func TestPostGISFunc(t *testing.T) {
	// 测试 PostGIS 函数字符串生成
	client := &Client{}
	pgFunc := client.PostGISFunc()

	assert.Contains(t, pgFunc.ST_Point(39.9042, 116.4074), "ST_Point")
	assert.Contains(t, pgFunc.ST_Distance(39.9042, 116.4074, 31.2304, 121.4737), "ST_Distance")
	assert.Contains(t, pgFunc.ST_Buffer("geom", 1000.0), "ST_Buffer")
	assert.Contains(t, pgFunc.ST_AsGeoJSON("geom"), "ST_AsGeoJSON")
}

func TestVectorFunc(t *testing.T) {
	client := &Client{}
	vecFunc := client.VectorFunc()

	vec := Vector{1.0, 2.0, 3.0}

	assert.Contains(t, vecFunc.VectorToString(vec), "ARRAY")
	assert.Contains(t, vecFunc.L2Distance("col1", "col2"), "<->")
	assert.Contains(t, vecFunc.CosineDistance("col1", "col2"), "<=>")
}

func TestGetLogger(t *testing.T) {
	levels := []string{"silent", "error", "warn", "info", "debug"}

	for _, level := range levels {
		logger := getLogger(level, 200*time.Millisecond)
		assert.NotNil(t, logger)
	}
}
