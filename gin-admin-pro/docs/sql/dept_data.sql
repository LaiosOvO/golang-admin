-- 部门管理测试数据
INSERT INTO system_dept (id, parent_id, level, sort, name, path, leader_user_id, phone, email, status, ancestors, create_by, update_by, created_at, updated_at) VALUES
-- 总公司
(1, 0, 1, 1, '总公司', '/1', NULL, '15888888888', 'admin@gin-admin.com', 1, '0', 1, 1, NOW(), NOW()),

-- 技术部
(2, 1, 2, 1, '技术部', '/1/2', 1, '15888888888', 'admin@gin-admin.com', 1, '0,1', 1, 1, NOW(), NOW()),
(3, 2, 3, 1, '研发组', '/1/2/3', 1, '15888888888', 'admin@gin-admin.com', 1, '0,1,2', 1, 1, NOW(), NOW()),
(4, 2, 3, 2, '测试组', '/1/2/4', 1, '15888888888', 'admin@gin-admin.com', 1, '0,1,2', 1, 1, NOW(), NOW()),
(5, 2, 3, 3, '运维组', '/1/2/5', 1, '15888888888', 'admin@gin-admin.com', 1, '0,1,2', 1, 1, NOW(), NOW()),

-- 市场部
(6, 1, 2, 2, '市场部', '/1/6', 1, '15888888888', 'admin@gin-admin.com', 1, '0,1', 1, 1, NOW(), NOW()),
(7, 6, 3, 1, '销售组', '/1/6/7', 1, '15888888888', 'admin@gin-admin.com', 1, '0,1,6', 1, 1, NOW(), NOW()),
(8, 6, 3, 2, '推广组', '/1/6/8', 1, '15888888888', 'admin@gin-admin.com', 1, '0,1,6', 1, 1, NOW(), NOW()),

-- 财务部
(9, 1, 2, 3, '财务部', '/1/9', 1, '15888888888', 'admin@gin-admin.com', 1, '0,1', 1, 1, NOW(), NOW()),

-- 人事部
(10, 1, 2, 4, '人事部', '/1/10', 1, '15888888888', 'admin@gin-admin.com', 1, '0,1', 1, 1, NOW(), NOW()),

-- 分公司
(11, 0, 1, 2, '深圳分公司', '/11', 1, '15888888888', 'admin@gin-admin.com', 1, '0', 1, 1, NOW(), NOW()),
(12, 11, 2, 1, '深圳技术部', '/11/12', 1, '15888888888', 'admin@gin-admin.com', 1, '0,11', 1, 1, NOW(), NOW()),
(13, 11, 2, 2, '深圳市场部', '/11/13', 1, '15888888888', 'admin@gin-admin.com', 1, '0,11', 1, 1, NOW(), NOW()),

-- 人力资源部（深圳）
(14, 11, 2, 3, '深圳人事部', '/11/14', 1, '15888888888', 'admin@gin-admin.com', 1, '0,11', 1, 1, NOW(), NOW());