-- 菜单管理测试数据
INSERT INTO system_menu (id, parent_id, level, sort, name, path, component, component_name, icon, type, perms, status, visible, keep_alive, always_show, ancestors, create_by, update_by, created_at, updated_at) VALUES
-- 系统管理模块
(1, 0, 1, 10, '系统管理', '/system', '', '', 'system', 1, '', 1, 1, 1, 1, '0', 1, 1, NOW(), NOW()),

-- 用户管理
(2, 1, 2, 10, '用户管理', 'user', 'system/user/index', 'User', 'user', 2, 'system:user:list', 1, 1, 1, 1, '0,1', 1, 1, NOW(), NOW()),
(3, 2, 3, 10, '用户查询', '', '', '', '', 3, 'system:user:query', 1, 1, 1, 1, '0,1,2', 1, 1, NOW(), NOW()),
(4, 2, 3, 20, '新增用户', '', '', '', '', 3, 'system:user:create', 1, 1, 1, 1, '0,1,2', 1, 1, NOW(), NOW()),
(5, 2, 3, 30, '修改用户', '', '', '', '', 3, 'system:user:update', 1, 1, 1, 1, '0,1,2', 1, 1, NOW(), NOW()),
(6, 2, 3, 40, '删除用户', '', '', '', '', 3, 'system:user:delete', 1, 1, 1, 1, '0,1,2', 1, 1, NOW(), NOW()),
(7, 2, 3, 50, '导出用户', '', '', '', '', 3, 'system:user:export', 1, 1, 1, 1, '0,1,2', 1, 1, NOW(), NOW()),

-- 角色管理
(8, 1, 2, 20, '角色管理', 'role', 'system/role/index', 'Role', 'peoples', 2, 'system:role:list', 1, 1, 1, 1, '0,1', 1, 1, NOW(), NOW()),
(9, 8, 3, 10, '角色查询', '', '', '', '', 3, 'system:role:query', 1, 1, 1, 1, '0,1,8', 1, 1, NOW(), NOW()),
(10, 8, 3, 20, '新增角色', '', '', '', '', 3, 'system:role:create', 1, 1, 1, 1, '0,1,8', 1, 1, NOW(), NOW()),
(11, 8, 3, 30, '修改角色', '', '', '', '', 3, 'system:role:update', 1, 1, 1, 1, '0,1,8', 1, 1, NOW(), NOW()),
(12, 8, 3, 40, '删除角色', '', '', '', '', 3, 'system:role:delete', 1, 1, 1, 1, '0,1,8', 1, 1, NOW(), NOW()),
(13, 8, 3, 50, '角色权限', '', '', '', '', 3, 'system:role:permission', 1, 1, 1, 1, '0,1,8', 1, 1, NOW(), NOW()),

-- 菜单管理
(14, 1, 2, 30, '菜单管理', 'menu', 'system/menu/index', 'Menu', 'tree-table', 2, 'system:menu:list', 1, 1, 1, 1, '0,1', 1, 1, NOW(), NOW()),
(15, 14, 3, 10, '菜单查询', '', '', '', '', 3, 'system:menu:query', 1, 1, 1, 1, '0,1,14', 1, 1, NOW(), NOW()),
(16, 14, 3, 20, '新增菜单', '', '', '', '', 3, 'system:menu:create', 1, 1, 1, 1, '0,1,14', 1, 1, NOW(), NOW()),
(17, 14, 3, 30, '修改菜单', '', '', '', '', 3, 'system:menu:update', 1, 1, 1, 1, '0,1,14', 1, 1, NOW(), NOW()),
(18, 14, 3, 40, '删除菜单', '', '', '', '', 3, 'system:menu:delete', 1, 1, 1, 1, '0,1,14', 1, 1, NOW(), NOW()),

-- 部门管理
(19, 1, 2, 40, '部门管理', 'dept', 'system/dept/index', 'Dept', 'tree', 2, 'system:dept:list', 1, 1, 1, 1, '0,1', 1, 1, NOW(), NOW()),
(20, 19, 3, 10, '部门查询', '', '', '', '', 3, 'system:dept:query', 1, 1, 1, 1, '0,1,19', 1, 1, NOW(), NOW()),
(21, 19, 3, 20, '新增部门', '', '', '', '', 3, 'system:dept:create', 1, 1, 1, 1, '0,1,19', 1, 1, NOW(), NOW()),
(22, 19, 3, 30, '修改部门', '', '', '', '', 3, 'system:dept:update', 1, 1, 1, 1, '0,1,19', 1, 1, NOW(), NOW()),
(23, 19, 3, 40, '删除部门', '', '', '', '', 3, 'system:dept:delete', 1, 1, 1, 1, '0,1,19', 1, 1, NOW(), NOW()),

-- 基础设施模块
(24, 0, 1, 20, '基础设施', '/infra', '', '', 'build', 1, '', 1, 1, 1, 1, '0', 1, 1, NOW(), NOW()),

-- 文件管理
(25, 24, 2, 10, '文件管理', 'file', 'infra/file/index', 'File', 'upload', 2, 'infra:file:list', 1, 1, 1, 1, '0,24', 1, 1, NOW(), NOW()),
(26, 25, 3, 10, '文件查询', '', '', '', '', 3, 'infra:file:query', 1, 1, 1, 1, '0,24,25', 1, 1, NOW(), NOW()),
(27, 25, 3, 20, '文件上传', '', '', '', '', 3, 'infra:file:upload', 1, 1, 1, 1, '0,24,25', 1, 1, NOW(), NOW()),
(28, 25, 3, 30, '文件删除', '', '', '', '', 3, 'infra:file:delete', 1, 1, 1, 1, '0,24,25', 1, 1, NOW(), NOW()),

-- AI模块
(29, 0, 1, 30, 'AI对话', '/ai', '', '', 'chat-dot-round', 1, '', 1, 1, 1, 1, '0', 1, 1, NOW(), NOW()),
(30, 29, 2, 10, '对话聊天', 'chat', 'ai/chat/index', 'Chat', 'message', 2, 'ai:chat', 1, 1, 1, 1, '0,29', 1, 1, NOW(), NOW());

-- 创建角色菜单关联
INSERT INTO system_role_menu (role_id, menu_id) VALUES
-- 管理员拥有所有菜单权限
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7),
(1, 8), (1, 9), (1, 10), (1, 11), (1, 12), (1, 13),
(1, 14), (1, 15), (1, 16), (1, 17), (1, 18),
(1, 19), (1, 20), (1, 21), (1, 22), (1, 23),
(1, 24), (1, 25), (1, 26), (1, 27), (1, 28),
(1, 29), (1, 30);