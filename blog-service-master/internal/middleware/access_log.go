package middleware

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/pkg/logger"
)

// AccessLogWriter 针对访问日志的结构体
type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 可以直接通过 AccessLogWriter 的 body 取到值
func (w AccessLogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		//初始化AccessLogWriter，将其赋予给当前的 Writer 写入流
		bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyWriter

		beginTime := time.Now().Unix()
		c.Next()
		endTime := time.Now().Unix()
		//指定方法得到我们所需的日志属性
		fields := logger.Fields{
			"request":  c.Request.PostForm.Encode(),
			"response": bodyWriter.body.String(),
		}
		s := "access log: method: %s, status_code: %d, " +
			"begin_time: %d, end_time: %d"
		//将信息写入日志
		global.Logger.WithFields(fields).Infof(c, s,
			c.Request.Method,
			bodyWriter.Status(),
			beginTime,
			endTime,
		)
	}
}
