package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// ContextTimeout 上下文超时时间控制
func ContextTimeout(t time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		// context.WithTimeout 设置请求 context 的超时时间，放回内容
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		defer cancel()
		//并重新赋予给了 gin.Context,在当前请求运行到指定的时间后，使用 context 的流程会进行处理
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
