package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/internal/service"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
)

// GetAuth 校验及获取入参后，绑定并获取到的 app_key 和 app_secret
// 进行数据库查询，检查认证信息是否存在，存在则进行 Token 的生成并返回
func GetAuth(c *gin.Context) {
	//拿到入参
	param := service.AuthRequest{}
	response := app.NewResponse(c)
	//绑定
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//根据request创建服务
	svc := service.New(c.Request.Context())
	//通过服务坚持入参是否包含权限
	err := svc.CheckAuth(&param)
	if err != nil {
		global.Logger.Errorf(c, "svc.CheckAuth err: %v", err)
		response.ToErrorResponse(errcode.UnauthorizedAuthNotExist)
		return
	}
	//生成令牌
	token, err := app.GenerateToken(param.AppKey, param.AppSecret)
	if err != nil {
		global.Logger.Errorf(c, "app.GenerateToken err: %v", err)
		response.ToErrorResponse(errcode.UnauthorizedTokenGenerate)
		return
	}
	//将令牌存入响应中
	response.ToResponse(gin.H{
		"token": token,
	})
}
