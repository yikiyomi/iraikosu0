package middleware

import (
    "net/http"
    "strings"  
    "github.com/gin-gonic/gin"
	"allto/util"
)

// jwt验证
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少token"})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token格式错误"})
			c.Abort()
			return
		}
		claims, err := util.ParseToken(parts[1])
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
func SoftAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context){
		authHeader:=c.GetHeader("Authorization")
		if authHeader==""{
			// 没token当游客，userid=0
			c.Set("userID",uint(0))
			c.Set("username","")
			c.Next()
			return 
		}
		parts:=strings.SplitN(authHeader," ",2)
		if len(parts) !=2||parts[0]!="Bearer"{
			// token格式不对当游客
			c.Set("userID",uint(0))
			c.Set("username","")
			c.Next()
			return 
		}
		claims,err:=util.ParseToken(parts[1])
		if err!=nil{
			// token过期当游客
			c.Set("userID",uint(0))
			c.Set("username","")
			c.Next()
			return 
		}
		c.Set("userID",claims.UserID)
		c.Set("username",claims.Username)
		c.Next()
	}
}