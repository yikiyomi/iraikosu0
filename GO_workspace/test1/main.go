package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var jwtSecret = []byte("phasel_gin_secret")

type User struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Username string    `gorm:"unique;not nil" json:"username"`
	Password string    `gorm:"not null" json:"-"`
	CreateAt time.Time `json:"create_at"`
}
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	var err error
	dsn := "root:123456@tcp(127.0.0.1:3306)/allto?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("mysql连接失败", err)
	}
	db.AutoMigrate(&User{})

	r := gin.Default()

	r.POST("/register", register)
	r.POST("/login", login)
	auth := r.Group("/api")
	auth.Use(authMiddleware())
	{
		auth.GET("/profile", profile)
	}
	log.Println("服务器启动：8080")
	r.Run(":8080")
}
func genToken(userID uint, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if len(tokenStr) < 7 || tokenStr[:7] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "格式错误"})
			return
		}
		tokenStr = tokenStr[7:]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
func register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 检查是否存在
	var exist User
	if err := db.Where("username=?", req.Username).First(&exist).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username exists"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := User{Username: req.Username, Password: string(hashed)}
	db.Create(&user)
	c.JSON(http.StatusCreated, gin.H{"msg": "ok", "user_id": user.ID})
}
func login(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, _ := genToken(user.ID, user.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
func profile(c *gin.Context) {
	userID := c.GetUint("userID")
	username := c.GetString("username")
	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"username": username,
	})
}
