package middleware

import (
	"allto/database"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var rateLimitScript=redis.NewScript(`
	local count=redis.call("INCR",KEYS[1])
	if count ==1 then
	redis.call("PEXPIRE",KEYS[1],ARGV[1])
	end
	return count
`)


func RateLimit(limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		rdb:=database.GetRedis()
		if rdb==nil{
			//若redis不可用则放行，降级策略
			c.Next()
			return 
		}
		key :="rate_limit:"+ c.ClientIP()
		count, err:=rateLimitScript.Run(
			database.GetCtx(),
			rdb,
			[]string{key},
			60000,
		).Int()
		if err!=nil{
			//出错放行
			c.Next()
			return 
		}
		if count > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁",
			})
			return
		}
		c.Next()
	}
}