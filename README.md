# goflower-io Example

Full-stack working examples for the goflower-io ecosystem, showing **MySQL**, **PostgreSQL**, and **SQLite3** — each with a complete HTTP + gRPC service built from generated code.

[中文文档](README_zh.md)

---

## The goflower-io Ecosystem

```
┌─────────────────────────────────────────────────────┐
│                      example (this repo)             │
│  MySQL / PostgreSQL / SQLite3 full-stack demos       │
└───────────────────────────────┬─────────────────────┘
                                │ uses
         ┌──────────────────────┼────────────────────┐
         ▼                      ▼                    ▼
  ┌─────────────┐       ┌─────────────┐      ┌────────────┐
  │    crud     │       │    xsql     │      │   golib    │
  │  (codegen)  │       │ (SQL build) │      │(HTTP/gRPC) │
  └─────────────┘       └─────────────┘      └────────────┘
```

| Library | Role in this example |
|---|---|
| [crud](https://github.com/goflower-io/crud) | Generated all model code and gRPC service skeleton from `user.sql` |
| [xsql](https://github.com/goflower-io/xsql) | Runtime SQL builder and DB client used inside generated CRUD operations |
| [golib](https://github.com/goflower-io/golib) | HTTP request parsing, response helpers, and gRPC server wiring |

---

## Repository Layout

```
example/
├── mysql/
│   ├── main.go
│   ├── crud/
│   │   ├── aa_client.go         # DB client (xsql-backed, read-write separation)
│   │   ├── sql/user.sql         # Source DDL — input to crud codegen
│   │   └── user/user.go         # Generated model + CRUD operations
│   ├── proto/user.api.proto     # Generated gRPC service definition
│   ├── api/
│   │   ├── user.api.pb.go
│   │   └── user.api_grpc.pb.go
│   ├── service/
│   │   ├── user.service.go      # gRPC service implementation
│   │   └── user.http.go         # HTTP handler wiring (golib)
│   └── views/                   # Templ HTML components
├── postgres/                    # Same structure for PostgreSQL
└── sqlite/                      # Same structure for SQLite3
```

---

## Prerequisites

- Go 1.21+
- `crud` tool (for regenerating code): `go install github.com/goflower-io/crud@main`
- `grpcurl` (for gRPC testing): `go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`
- MySQL or PostgreSQL running (SQLite example needs nothing external)

---

## Run the Examples

### SQLite3 — No external DB needed

```bash
cd sqlite
# Edit main.go: set DSN to your preferred SQLite file path
go run main.go
```

### MySQL

```bash
# Create database
mysql -u root -p -e "CREATE DATABASE test CHARACTER SET utf8mb4;"

cd mysql
# Edit main.go: set DSN / ReadDSN to your MySQL credentials
go run main.go
# HTTP server starts at :8088
```

### PostgreSQL

```bash
psql -U postgres -c "CREATE DATABASE test;"

cd postgres
# Edit main.go: set DSN to your PostgreSQL credentials
go run main.go
```

---

## The SQL Schema

All three examples share the same logical table structure. The MySQL version:

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

> Column comment format `'label|input_type|validate_rules|enum_values'` — crud parses these to generate validation tags and HTML form hints automatically.

---

## HTTP API Testing

```bash
BASE=http://localhost:8088

# Create user
curl -X POST "$BASE/UserService/CreateUser" \
  -d "name=alice&age=18&sex=1"

# Get user
curl "$BASE/UserService/GetUser?id=1"

# List users (JSON response)
curl -H "Accept: application/json" \
  "$BASE/UserService/ListUsers?page=1&page_size=10"

# Update user (name + age only)
curl -X POST "$BASE/UserService/UpdateUser" \
  -d "user.id=1&user.name=bob&user.age=20&masks=2&masks=3"

# Delete user
curl "$BASE/UserService/DeleteUser?id=1"
```

---

## gRPC API Testing with grpcurl

To use gRPC directly, wire the service through [golib](https://github.com/goflower-io/golib)'s `App` in `main.go`:

```go
import "github.com/goflower-io/golib/net/app"

a := app.New(app.WithAddr("0.0.0.0", 8080))
a.RegisteGrpcService(&api.UserService_ServiceDesc, svc)
a.Run()
```

Then use grpcurl:

```bash
# List all services (gRPC reflection is registered automatically)
grpcurl -plaintext localhost:8080 list

# Describe the UserService interface
grpcurl -plaintext localhost:8080 describe UserService

# Create a user
grpcurl -plaintext \
  -d '{"name":"alice","age":18,"sex":1}' \
  localhost:8080 UserService/CreateUser

# Get user by ID
grpcurl -plaintext \
  -d '{"id":1}' \
  localhost:8080 UserService/GetUser

# Update specific fields (masks: 2=name, 3=age, 4=sex)
grpcurl -plaintext \
  -d '{"user":{"id":1,"name":"bob","age":25,"sex":0},"masks":[2,3]}' \
  localhost:8080 UserService/UpdateUser

# List with pagination
grpcurl -plaintext \
  -d '{"page":1,"page_size":10}' \
  localhost:8080 UserService/ListUsers

# List with filter (age > 18) and order by mtime desc
grpcurl -plaintext \
  -d '{
    "page": 1,
    "page_size": 10,
    "filters": [{"field": 3, "op": "GT", "val": "18"}],
    "orderbys": [{"field": 6, "desc": true}]
  }' \
  localhost:8080 UserService/ListUsers

# List with cursor-based pagination (ListUsersMore)
grpcurl -plaintext \
  -d '{"page_size":5,"cursor":{"orderbys":[{"field":1,"desc":false}]}}' \
  localhost:8080 UserService/ListUsersMore

# Delete a user
grpcurl -plaintext \
  -d '{"id":1}' \
  localhost:8080 UserService/DeleteUser

# Health check
grpcurl -plaintext localhost:8080 grpc.health.v1.Health/Check
```

### UserField enum values (for filters and orderbys)

| Value | Field |
|---|---|
| 1 | id |
| 2 | name |
| 3 | age |
| 4 | sex |
| 5 | ctime |
| 6 | mtime |

### Filter op values

| Op string | SQL operator |
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

## Regenerating Code

After modifying a `.sql` file:

```bash
cd mysql  # or postgres / sqlite

# Regenerate CRUD model only
crud -dialect mysql

# Regenerate CRUD model + proto + service skeleton
crud -dialect mysql -service -protopkg example
```

---

## Related

- [crud](https://github.com/goflower-io/crud) — the code generation tool
- [xsql](https://github.com/goflower-io/xsql) — SQL builder and DB client library
- [golib](https://github.com/goflower-io/golib) — HTTP/gRPC application server framework
