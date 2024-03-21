package middleware

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
)

// JWT 通过 GetHeader 方法从 Header 中获取 token 参数，
// 并调用 ParseToken 对其进行解析，再根据返回的错误类型进行断言判定
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token string
			ecode = errcode.Success
		)
		//拿到token
		if s, exist := c.GetQuery("token"); exist {
			token = s
		} else {
			token = c.GetHeader("token")
		}

		if token == "" {
			ecode = errcode.InvalidParams
		} else {
			_, err := app.ParseToken(token)
			if err != nil {
				//判断错误原因：失效或非法
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					ecode = errcode.UnauthorizedTokenTimeout
				default:
					ecode = errcode.UnauthorizedTokenError
				}
			}
		}
		//不成功就回退
		if ecode != errcode.Success {
			response := app.NewResponse(c)
			response.ToErrorResponse(ecode)
			c.Abort()
			return

		}
		//成功就执行下一个模块
		c.Next()
	}
}
