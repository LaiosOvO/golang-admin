// MongoDB初始化脚本

// 创建管理员用户
db = db.getSiblingDB("admin");

// 创建应用数据库用户
db.createUser({
  user: "gin_admin",
  pwd: "gin_admin123",
  roles: [
    {
      role: "readWrite",
      db: "gin_admin"
    }
  ]
});

// 切换到应用数据库
db = db.getSiblingDB("gin_admin");

// 创建基础集合和索引
db.createCollection("users");
db.users.createIndex({ "username": 1 }, { unique: true });
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "mobile": 1 }, { unique: true, sparse: true });

db.createCollection("roles");
db.roles.createIndex({ "code": 1 }, { unique: true });

db.createCollection("permissions");
db.permissions.createIndex({ "code": 1 }, { unique: true });

db.createCollection("audit_logs");
db.audit_logs.createIndex({ "created_at": -1 });
db.audit_logs.createIndex({ "user_id": 1 });
db.audit_logs.createIndex({ "action": 1 });

db.createCollection("file_storage");
db.file_storage.createIndex({ "filename": 1 });
db.file_storage.createIndex({ "uploaded_at": -1 });

// 插入初始数据
db.roles.insertMany([
  {
    code: "admin",
    name: "管理员",
    description: "系统管理员角色",
    permissions: ["*"],
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    code: "user",
    name: "普通用户",
    description: "普通用户角色",
    permissions: ["read", "write"],
    created_at: new Date(),
    updated_at: new Date()
  }
]);

print("MongoDB initialization completed successfully");