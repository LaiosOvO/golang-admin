package postgresql

import (
	"fmt"
	"math"
	"strings"
)

// Vector 向量类型
type Vector []float32

// String 返回向量的字符串表示
func (v Vector) String() string {
	if len(v) == 0 {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[")
	for i, val := range v {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%.6f", val))
	}
	sb.WriteString("]")
	return sb.String()
}

// Size 返回向量维度
func (v Vector) Size() int {
	return len(v)
}

// Normalize 向量归一化
func (v Vector) Normalize() Vector {
	if len(v) == 0 {
		return v
	}

	var norm float32
	for _, val := range v {
		norm += val * val
	}
	norm = float32(1.0) / float32(math.Sqrt(float64(norm)))

	normalized := make(Vector, len(v))
	for i, val := range v {
		normalized[i] = val * norm
	}
	return normalized
}

// CosineSimilarity 计算余弦相似度
func CosineSimilarity(a, b Vector) (float32, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vectors must have same dimension")
	}
	if len(a) == 0 {
		return 0, nil
	}

	var dotProduct, normA, normB float32

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0, nil
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB)))), nil
}

// EuclideanDistance 计算欧几里得距离
func EuclideanDistance(a, b Vector) (float32, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vectors must have same dimension")
	}
	if len(a) == 0 {
		return 0, nil
	}

	var distance float32
	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		distance += diff * diff
	}

	return float32(math.Sqrt(float64(distance))), nil
}

// VectorFunc pgvector 函数封装
type VectorFunc struct {
	client *Client
}

// NewVectorFunc 创建 Vector 函数实例
func (c *Client) VectorFunc() *VectorFunc {
	return &VectorFunc{client: c}
}

// ToVector 将向量转换为数据库格式
func (v *VectorFunc) ToVector(vec Vector) string {
	return fmt.Sprintf("ARRAY[%s]", strings.Join(v.toFloatStrings(vec), ","))
}

// toFloatStrings 将向量转换为字符串数组
func (v *VectorFunc) toFloatStrings(vec Vector) []string {
	strs := make([]string, len(vec))
	for i, val := range vec {
		strs[i] = fmt.Sprintf("%.6f", val)
	}
	return strs
}

// VectorToString 将向量转换为字符串
func (v *VectorFunc) VectorToString(vec Vector) string {
	return fmt.Sprintf("'%s'", v.ToVector(vec))
}

// L2Distance 计算L2距离（欧几里得距离）
func (v *VectorFunc) L2Distance(col1, col2 string) string {
	return fmt.Sprintf("(%s <-> %s)", col1, col2)
}

// CosineDistance 计算余弦距离
func (v *VectorFunc) CosineDistance(col1, col2 string) string {
	return fmt.Sprintf("(1 - (%s <=> %s))", col1, col2)
}

// L2DistanceToValue 计算与固定向量的L2距离
func (v *VectorFunc) L2DistanceToValue(column string, vec Vector) string {
	return fmt.Sprintf("(%s <-> %s)", column, v.VectorToString(vec))
}

// CosineDistanceToValue 计算与固定向量的余弦距离
func (v *VectorFunc) CosineDistanceToValue(column string, vec Vector) string {
	return fmt.Sprintf("(1 - (%s <=> %s))", column, v.VectorToString(vec))
}

// CreateVectorIndex 创建向量索引
func (v *VectorFunc) CreateVectorIndex(table, column string, dimension int, indexType string) error {
	var indexSQL string
	indexName := fmt.Sprintf("idx_%s_%s", table, column)

	switch strings.ToLower(indexType) {
	case "ivfflat":
		// IVFFlat 索引
		indexSQL = fmt.Sprintf("CREATE INDEX %s ON %s USING ivfflat (%s vector_cosine_ops) WITH (lists = 100)",
			indexName, table, column)
	case "hnsw":
		// HNSW 索引
		indexSQL = fmt.Sprintf("CREATE INDEX %s ON %s USING hnsw (%s vector_cosine_ops)",
			indexName, table, column)
	default:
		// 默认使用 IVFFlat
		indexSQL = fmt.Sprintf("CREATE INDEX %s ON %s USING ivfflat (%s vector_cosine_ops) WITH (lists = 100)",
			indexName, table, column)
	}

	db := v.client.GetDB()
	return db.Exec(indexSQL).Error
}

// DropVectorIndex 删除向量索引
func (v *VectorFunc) DropVectorIndex(table, column string) error {
	indexName := fmt.Sprintf("idx_%s_%s", table, column)
	db := v.client.GetDB()
	return db.Exec(fmt.Sprintf("DROP INDEX IF EXISTS %s", indexName)).Error
}

// SimilaritySearch 相似性搜索
func (v *VectorFunc) SimilaritySearch(table, column string, queryVec Vector, limit int, orderBy string) string {
	distanceCol := v.L2DistanceToValue(column, queryVec)

	if orderBy == "" {
		orderBy = "ASC"
	}

	return fmt.Sprintf("SELECT *, %s as distance FROM %s ORDER BY %s %s LIMIT %d",
		distanceCol, table, distanceCol, orderBy, limit)
}
