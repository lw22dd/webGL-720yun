package middleware

import (
	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
// 推荐的中间件顺序
// 1. Recovery (最外层,捕获panic)
// 2. Logger (记录所有请求)
// 3. CORS (跨域处理)
// 4. RateLimiter (限流)
// 5. Auth (认证,仅部分路由需要)
// 6. Permission (权限校验,仅部分路由需要)

// CORS CORS中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
