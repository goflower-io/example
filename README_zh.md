# goflower-io 示例代码

goflower-io 生态系统的全栈工作示例，展示 **MySQL**、**PostgreSQL** 和 **SQLite3** 三种数据库——每个示例都包含由生成代码构建的完整 HTTP + gRPC 服务。

[English](README.md)

---

## goflower-io 生态系统

```
┌─────────────────────────────────────────────────────┐
│                  example（本仓库）                    │
│  MySQL / PostgreSQL / SQLite3 全栈示例               │
└───────────────────────────────┬─────────────────────┘
                                │ 使用
         ┌──────────────────────┼────────────────────┐
         ▼                      ▼                    ▼
  ┌─────────────┐       ┌─────────────┐      ┌────────────┐
  │    crud     │       │    xsql     │      │   golib    │
  │  （代码生成） │       │（SQL 构建器）│      │（HTTP/gRPC）│
  └─────────────┘       └─────────────┘      └────────────┘
```

| 库 | 在本示例中的角色 |
|---|---|
| [crud](https://github.com/goflower-io/crud) | 从 `user.sql` 生成了全部 Model 代码和 gRPC 服务骨架 |
| [xsql](https://github.com/goflower-io/xsql) | 生成的 CRUD 操作内部使用的运行时 SQL 构建器和 DB 客户端 |
| [golib](https://github.com/goflower-io/golib) | HTTP 请求解析、响应助手和 gRPC 服务注册 |

---

## 目录结构

```
example/
├── mysql/
│   ├── main.go
│   ├── crud/
│   │   ├── aa_client.go         # DB 客户端（xsql 支持，读写分离）
│   │   ├── sql/user.sql         # 源 DDL — crud 代码生成的输入
│   │   └── user/user.go         # 生成的 Model + CRUD 操作
│   ├── proto/user.api.proto     # 生成的 gRPC 服务定义
│   ├── api/
│   │   ├── user.api.pb.go
│   │   └── user.api_grpc.pb.go
│   ├── service/
│   │   ├── user.service.go      # gRPC 服务实现
│   │   └── user.http.go         # HTTP 处理器接入（golib）
│   └── views/                   # Templ HTML 组件
├── postgres/                    # PostgreSQL 示例（结构相同）
└── sqlite/                      # SQLite3 示例（结构相同）
```

---

## 环境要求

- Go 1.21+
- `crud` 工具（重新生成代码时需要）：`go install github.com/goflower-io/crud@main`
- `grpcurl`（gRPC 测试）：`go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`
- MySQL 或 PostgreSQL 实例（SQLite3 示例无需外部数据库）

---

## 运行示例

### SQLite3 — 无需外部数据库

```bash
cd sqlite
# 编辑 main.go，将 DSN 设置为你的 SQLite 文件路径
go run main.go
```

### MySQL

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE test CHARACTER SET utf8mb4;"

cd mysql
# 编辑 main.go，设置 DSN / ReadDSN 为你的 MySQL 连接信息
go run main.go
# HTTP 服务监听 :8088
```

### PostgreSQL

```bash
psql -U postgres -c "CREATE DATABASE test;"

cd postgres
# 编辑 main.go，设置 DSN 为你的 PostgreSQL 连接信息
go run main.go
```

---

## SQL 建表语句

三个示例使用相同的逻辑表结构，MySQL 版本如下：

```sql
CREATE TABLE `user` (
  `id`    int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `name`  varchar(100)     NOT NULL COMMENT '名称|text|validate:"max=100,min=1"',
  `age`   int(11)          NOT NULL DEFAULT '0'  COMMENT '年龄|number|validate:"max=140,min=0"',
  `sex`   int(11)          NOT NULL DEFAULT '2'  COMMENT '性别|select|validate:"oneof=0 1 2"|0:女 1:男 2:无',
  `ctime` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `ix_name`  (`name`)  USING BTREE,
  KEY `ix_mtime` (`mtime`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

> 字段注释格式 `'标签|输入类型|校验规则|枚举值'`——crud 自动解析这些元数据，生成校验 tag 和 HTML 表单提示。

---

## HTTP 接口测试

```bash
BASE=http://localhost:8088

# 创建用户
curl -X POST "$BASE/UserService/CreateUser" \
  -d "name=alice&age=18&sex=1"

# 查询用户
curl "$BASE/UserService/GetUser?id=1"

# 查询用户列表（JSON 响应）
curl -H "Accept: application/json" \
  "$BASE/UserService/ListUsers?page=1&page_size=10"

# 更新用户（只更新 name 和 age）
curl -X POST "$BASE/UserService/UpdateUser" \
  -d "user.id=1&user.name=bob&user.age=20&masks=2&masks=3"

# 删除用户
curl "$BASE/UserService/DeleteUser?id=1"
```

---

## 使用 grpcurl 测试 gRPC 接口

通过 [golib](https://github.com/goflower-io/golib) 的 `App` 在 `main.go` 中接入 gRPC：

```go
import "github.com/goflower-io/golib/net/app"

a := app.New(app.WithAddr("0.0.0.0", 8080))
a.RegisteGrpcService(&api.UserService_ServiceDesc, svc)
a.Run()
```

然后使用 grpcurl 进行测试：

```bash
# 列出所有服务（golib 自动注册 gRPC 反射）
grpcurl -plaintext localhost:8080 list

# 查看 UserService 接口描述
grpcurl -plaintext localhost:8080 describe UserService

# 创建用户
grpcurl -plaintext \
  -d '{"name":"alice","age":18,"sex":1}' \
  localhost:8080 UserService/CreateUser

# 按 ID 查询用户
grpcurl -plaintext \
  -d '{"id":1}' \
  localhost:8080 UserService/GetUser

# 更新指定字段（masks: 2=name, 3=age, 4=sex）
grpcurl -plaintext \
  -d '{"user":{"id":1,"name":"bob","age":25,"sex":0},"masks":[2,3]}' \
  localhost:8080 UserService/UpdateUser

# 分页查询用户列表
grpcurl -plaintext \
  -d '{"page":1,"page_size":10}' \
  localhost:8080 UserService/ListUsers

# 带过滤条件（age > 18）和排序（mtime 倒序）
grpcurl -plaintext \
  -d '{
    "page": 1,
    "page_size": 10,
    "filters": [{"field": 3, "op": "GT", "val": "18"}],
    "orderbys": [{"field": 6, "desc": true}]
  }' \
  localhost:8080 UserService/ListUsers

# 基于游标的分页查询（ListUsersMore）
grpcurl -plaintext \
  -d '{"page_size":5,"cursor":{"orderbys":[{"field":1,"desc":false}]}}' \
  localhost:8080 UserService/ListUsersMore

# 删除用户
grpcurl -plaintext \
  -d '{"id":1}' \
  localhost:8080 UserService/DeleteUser

# 健康检查
grpcurl -plaintext localhost:8080 grpc.health.v1.Health/Check
```

### UserField 枚举值（用于 filters 和 orderbys）

| 枚举值 | 对应字段 |
|---|---|
| 1 | id |
| 2 | name |
| 3 | age |
| 4 | sex |
| 5 | ctime |
| 6 | mtime |

### Filter op 操作符

| op 字符串 | SQL 操作符 |
|---|---|
| `EQ` | `=` |
| `NEQ` | `!=` |
| `GT` | `>` |
| `GTE` | `>=` |
| `LT` | `<` |
| `LTE` | `<=` |
| `IN` | `IN (...)` |
| `CONTAINS` | `LIKE '%val%'` |
| `HAS_PREFIX` | `LIKE 'val%'` |

---

## 重新生成代码

修改 `.sql` 文件后，在对应示例目录执行：

```bash
cd mysql  # 或 postgres / sqlite

# 只重新生成 CRUD Model
crud -dialect mysql

# 同时生成 proto 文件和服务骨架
crud -dialect mysql -service -protopkg example
```

---

## 相关仓库

- [crud](https://github.com/goflower-io/crud) — 代码生成工具
- [xsql](https://github.com/goflower-io/xsql) — SQL 构建器和 DB 客户端库
- [golib](https://github.com/goflower-io/golib) — HTTP/gRPC 应用服务器框架
