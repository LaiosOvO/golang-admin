package postgresql

import (
	"fmt"
	"math"
)

// Point 地理坐标点
type Point struct {
	Latitude  float64 `json:"latitude"`  // 纬度
	Longitude float64 `json:"longitude"` // 经度
}

// Geometry 几何类型
type Geometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

// PostGISFunc PostGIS 函数封装
type PostGISFunc struct {
	client *Client
}

// NewPostGISFunc 创建 PostGIS 函数实例
func (c *Client) PostGISFunc() *PostGISFunc {
	return &PostGISFunc{client: c}
}

// ST_Point 创建点
func (p *PostGISFunc) ST_Point(lat, lng float64) string {
	return fmt.Sprintf("ST_Point(%f, %f)", lng, lat)
}

// ST_Distance 计算两点间距离（米）
func (p *PostGISFunc) ST_Distance(lat1, lng1, lat2, lng2 float64) string {
	return fmt.Sprintf("ST_Distance(ST_Point(%f, %f), ST_Point(%f, %f))", lng1, lat1, lng2, lat2)
}

// ST_Contains 判断几何对象包含关系
func (p *PostGISFunc) ST_Contains(geom1, geom2 string) string {
	return fmt.Sprintf("ST_Contains(%s, %s)", geom1, geom2)
}

// ST_Within 判断几何对象是否在另一个几何对象内
func (p *PostGISFunc) ST_Within(geom1, geom2 string) string {
	return fmt.Sprintf("ST_Within(%s, %s)", geom1, geom2)
}

// ST_Buffer 创建缓冲区
func (p *PostGISFunc) ST_Buffer(geom string, radius float64) string {
	return fmt.Sprintf("ST_Buffer(%s, %f)", geom, radius)
}

// ST_Area 计算面积
func (p *PostGISFunc) ST_Area(geom string) string {
	return fmt.Sprintf("ST_Area(%s)", geom)
}

// ST_AsGeoJSON 将几何对象转换为 GeoJSON
func (p *PostGISFunc) ST_AsGeoJSON(geom string) string {
	return fmt.Sprintf("ST_AsGeoJSON(%s)", geom)
}

// ST_AsText 将几何对象转换为 WKT
func (p *PostGISFunc) ST_AsText(geom string) string {
	return fmt.Sprintf("ST_AsText(%s)", geom)
}

// ST_GeomFromText 从 WKT 创建几何对象
func (p *PostGISFunc) ST_GeomFromText(wkt string) string {
	return fmt.Sprintf("ST_GeomFromText('%s')", wkt)
}

// ST_MakePoint 创建点
func (p *PostGISFunc) ST_MakePoint(lng, lat float64) string {
	return fmt.Sprintf("ST_MakePoint(%f, %f)", lng, lat)
}

// DistanceBetween 计算两点间距离（使用 Haversine 公式，米）
func DistanceBetween(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371000 // 地球半径，米

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// CreateBoundingBox 创建边界框
func CreateBoundingBox(lat, lng, radiusKm float64) (minLat, minLng, maxLat, maxLng float64) {
	// 大约每度距离
	const kmPerDegreeLat = 110.574
	kmPerDegreeLng := 111.320 * math.Cos(lat*math.Pi/180)

	// 计算边界框
	deltaLat := radiusKm / kmPerDegreeLat
	deltaLng := radiusKm / kmPerDegreeLng

	minLat = lat - deltaLat
	maxLat = lat + deltaLat
	minLng = lng - deltaLng
	maxLng = lng + deltaLng

	return minLat, minLng, maxLat, maxLng
}

// PointInPolygon 判断点是否在多边形内
func PointInPolygon(point Point, polygon []Point) bool {
	if len(polygon) < 3 {
		return false
	}

	intersections := 0
	n := len(polygon)

	for i := 0; i < n; i++ {
		j := (i + 1) % n

		// 检查射线与边的交点
		if ((polygon[i].Latitude > point.Latitude) != (polygon[j].Latitude > point.Latitude)) &&
			(point.Longitude < (polygon[j].Longitude-polygon[i].Longitude)*(point.Latitude-polygon[i].Latitude)/(polygon[j].Latitude-polygon[i].Latitude)+polygon[i].Longitude) {
			intersections++
		}
	}

	return intersections%2 == 1
}
