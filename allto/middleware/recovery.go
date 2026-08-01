package middleware

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic: %v", err)
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code": 500, "msg": "服务器内部错误",})
            }
        }()
        c.Next()
    }
}