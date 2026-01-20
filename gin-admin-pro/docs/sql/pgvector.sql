-- PostgreSQL pgvector扩展初始化
CREATE EXTENSION IF NOT EXISTS vector;

-- 验证扩展安装
SELECT version FROM pg_extension WHERE extname = 'vector';

-- 创建向量数据类型示例表（可选）
CREATE TABLE IF NOT EXISTS vector_example (
    id SERIAL PRIMARY KEY,
    content TEXT,
    embedding VECTOR(128)  -- 128维向量
);

-- 创建向量索引示例（可选）
-- CREATE INDEX ON vector_example USING ivfflat (embedding vector_l2_ops) WITH (lists = 100);