# allto 博客

基于 Go + Gin 的博客系统，含账号、文章、评论、点赞、关注，
支持 Redis 缓存、JWT 鉴权、限流、Docker Compose 部署。

## 功能

| 模块 | 功能 |
|------|------|
| 账号 | 注册、登录、JWT 鉴权 |
| 文章 | 发布、列表、详情、分页、Redis 阅读
| 点赞 | 点赞、取消点赞、已赞列表、事务保证一致性 |
| 评论 | 发布、删除、列表 |
| 关注 | 关注、取关、粉丝列表、关注列表 |
| 工程 | Viper 配置管理、.env 环境变量、限流、panic 恢复、请求日志 |

## 技术栈

后端：Go 1.26 + Gin + GORM + go-redis + JWT
前端：原生 JavaScript SPA（单页应用）

## 快速启动

### 前置条件

- Go 1.26+
- MySQL 8.0
- Redis 7

### 本地运行

# 1. 安装依赖
go mod tidy

# 2. 配置文件
#  方式 A：编辑 config.yaml
#  方式 B：设置环境变量后启动

# 3. 运行
go run main.go

# 访问前端：http://localhost:8080/static

## 配置说明

| 参数 | 默认值 | 说明 |
 redis.addr | localhost:6379 | Redis 地址 |
| redis.password | 无 | Redis 密码 |
| jwt.secret | my_secret_key | JWT 签名密钥，生产须改 |

## 接口清单

### 公开接口
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /register | 注册 |
| POST | /login | 登录，返回 JWT |

### 文章 /api/posts（需 JWT）
 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/posts | 发布文章 |
| GET | /api/posts | 文章列表（分页） |
| GET | /api/posts/:id | 文章详情 |
| POST | /api/posts/:id/like | 点赞 |
| DELETE | /api/posts/:id/like | 取消点赞 |
| POST | /api/posts/:id/comments | 发表评论 |
| GET | /api/posts/:id/comments | 评论列表 |
| GET | /api/posts/:id/likes | 文章点赞列表 |
| GET | /api/posts/:id/likers | 点赞用户列表
| GET | /api/posts/likes | 我的点赞 |

### 关注 /api/follow（需 JWT）
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/follow/:id | 关注用户 |
| DELETE | /api/follow/:id | 取消关注 |
| GET | /api/following | 我的关注 |
| GET | /api/followers | 我的粉丝 |

## 目录结构

allto/
├── config/         配置加载（Viper）
├── database/       MySQL + Redis 连接
├── frontend/       前端 SPA
├── handler/        HTTP 处理函数
├── middleware/     JWT 鉴权、限流、CORS、日
├── model/          GORM 数据模型
├── response/       统一响应格式
├── router/         Gin 路由注册
├── util/           JWT 工具
├── config.yaml     配置文件
└── main.go         入口

---