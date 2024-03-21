package middleware

import (
	"fmt"
	"time"

	"github.com/go-programming-tour-book/blog-service/pkg/email"

	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
)

func Recovery() gin.HandlerFunc {
	// Mailer 是固定的
	defaultMailer := email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	})
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Logger.WithCallersFrames().Errorf(c, "panic recover err: %v", err)
				//捕获到异常后调用 SendMail 方法进行预警邮件发送
				err := defaultMailer.SendMail(
					global.EmailSetting.To,
					fmt.Sprintf("异常抛出，发生时间: %d", time.Now().Unix()),
					fmt.Sprintf("错误信息: %v", err),
				)
				//记录邮件错误的情况
				if err != nil {
					global.Logger.Panicf(c, "mail.SendMail err: %v", err)
				}
				//返回响应
				app.NewResponse(c).ToErrorResponse(errcode.ServerError)
				c.Abort()
			}
		}()
		c.Next()
	}
}
