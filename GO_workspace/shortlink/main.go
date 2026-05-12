package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局变量
var (
	db  *gorm.DB
	rdb *redis.Client
	ctx = context.Background()
)
var jwtKey = []byte("shortlink_secret")

func main() {
	// 初始化数据库连接
	dsn := "root:123456@tcp(127.0.0.1:3306)/shortlink?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("mysql连接失败:", err)
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("redis连接失败:", err)
	}

	// 自动迁移数据库
	db.AutoMigrate(&User{}, &ShortLink{})

	r := gin.Default()

	//  公开路由
	r.POST("/register", Register)
	r.POST("/login", Login)

	// 短链接重定向
	r.GET("/:code", Redirect)

	// 受保护路由
	auth := r.Group("/miao")
	auth.Use(AuthMiddleware())
	{
		auth.POST("/shorten", CreateShortLink)  // 生成短链接
		auth.GET("/stats/:code", GetClickStats) // 查看点击统计
		auth.GET("/links", ListMyLinks)         // 列出我创建的链接
	}

	log.Println("短链接服务启动:8080")
	r.Run(":8080")
}

// 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// 短链接模型
type ShortLink struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	ShortCode string    `gorm:"unique;not null" json:"short_code"`
	LongURL   string    `gorm:"not null" json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// jwt
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// 中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeeader := c.GetHeader("Authorization")
		if authHeeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少token"})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token格式错误"})
			c.Abort()
			return
		}
		claims, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效token"})
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// base62字符集和自增id策略
const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func idTobase62(id uint) string {
	if id == 0 {
		return "0"
	}
	var result []byte
	for id > 0 {
		remainder := id % 62
		result = append([]byte{base62Chars[remainder]}, result...)
		id = id / 62
	}
	return string(result)
}

// 业务接口
// 注册
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}
	var exist User
	if err := db.Where("username = ?", req.Username).First(&exist).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := User{Username: req.Username, Password: string(hashed)}
	db.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "注册成功", "user_id": user.ID})
}

// 登录
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user User
	if err := db.Where("username=?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户或密码错误"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户或密码错误"})
		return
	}
	token, _ := GenerateToken(user.ID, user.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// 生成短链接
func CreateShortLink(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		LongURL string `json:"long_url" binding:"required,url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url格式错误"})
		return
	}
	// 先插入记录获取自增ID
	link := ShortLink{
		UserID:    userID,
		LongURL:   req.LongURL,
		ExpiresAt: time.Now().AddDate(100, 0, 0),
	}
	db.Create(&link)

	// 根据ID生成短码
	code := idTobase62(link.ID)
	db.Model(&link).Update("short_code", code)

	// redis缓存短链接
	rdb.Set(ctx, "shortlink:"+code, req.LongURL, 0)
	// 初始化点击统计
	rdb.Set(ctx, "clicks:"+code, 0, 0)
	c.JSON(http.StatusOK, gin.H{
		"short_code": code,
		"short_url":  fmt.Sprintf("http://localhost:8080/%s", code),
		"long_url":   req.LongURL,
	})
}

// 短链接重定向
func Redirect(c *gin.Context) {
	code := c.Param("code")
	// 先查redis缓存
	longURL, err := rdb.Get(ctx, "shortlink:"+code).Result()
	if err == redis.Nil {
		var link ShortLink
		if err := db.Where("short_code = ?", code).First(&link).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "链接不存在"})
			return
		}
		longURL = link.LongURL
		// 更新redis缓存
		rdb.Set(ctx, "shortlink:"+code, longURL, 24*time.Hour)
	}
	rdb.Incr(ctx, "clicks:"+code)
	// 重定向到原始URL  302 Found
	c.Redirect(http.StatusFound, longURL)
}

// 查看点击统计
func GetClickStats(c *gin.Context) {
	code := c.Param("code")
	userID := c.GetUint("userID")
	// 验证链接归属用户
	var link ShortLink
	if err := db.Where("short_code = ? and user_id = ?", code, userID).First(&link).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问该链接统计"})
		return
	}
	clicks, _ := rdb.Get(ctx, "clicks:"+code).Int()
	c.JSON(http.StatusOK, gin.H{
		"short_code": code,
		"long_url":   link.LongURL,
		"clicks":     clicks,
	})
}

// 列出我创建的链接
func ListMyLinks(c *gin.Context) {
	userID := c.GetUint("userID")
	var links []ShortLink
	db.Where("user_id=?", userID).Order("created_at desc").Find(&links)

	type LinkWithClicks struct {
		ShortLink
		Clicks int `json:"clicks"`
	}
	result := make([]LinkWithClicks, len(links))
	for i, link := range links {
		clicks, _ := rdb.Get(ctx, "clicks:"+link.ShortCode).Int()
		result[i] = LinkWithClicks{
			ShortLink: link,
			Clicks:    clicks,
		}
	}
	c.JSON(http.StatusOK, gin.H{"links": result})
}
