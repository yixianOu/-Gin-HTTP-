package routers

import (
	swaggerFiles "github.com/swaggo/files"
	"net/http"
	"time"

	"github.com/go-programming-tour-book/blog-service/pkg/limiter"

	"github.com/go-programming-tour-book/blog-service/global"

	"github.com/gin-gonic/gin"
	_ "github.com/go-programming-tour-book/blog-service/docs"
	"github.com/go-programming-tour-book/blog-service/internal/middleware"
	"github.com/go-programming-tour-book/blog-service/internal/routers/api"
	"github.com/go-programming-tour-book/blog-service/internal/routers/api/v1"
	ginSwagger "github.com/swaggo/gin-swagger"
	//"github.com/swaggo/gin-swagger/swaggerFiles"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(
	limiter.LimiterBucketRule{
		Key:          "/auth",
		FillInterval: time.Second,
		Capacity:     10,
		Quantum:      10,
	},
)

// NewRouter 将路由的 Handler 方法注册到对应的路由规则上
func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}

	r.Use(middleware.Tracing())
	r.Use(middleware.RateLimiter(methodLimiters))
	r.Use(middleware.ContextTimeout(global.AppSetting.DefaultContextTimeout))
	r.Use(middleware.Translations())

	article := v1.NewArticle()
	tag := v1.NewTag()
	//new一个上传器
	upload := api.NewUpload()
	r.GET("/debug/vars", api.Expvar)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//将上传器中的上传方法注册到上传路径上
	r.POST("/upload/file", upload.UploadFile)
	r.POST("/auth", api.GetAuth)
	//设置文件服务去提供静态资源的访问,绑定了项目的upload目录
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
	//利用 gin 中的分组路由的概念，只针对 apiv1 的路由分组进行 JWT 中间件的引用
	//只有 apiv1 路由分组里的路由方法会受此中间件的约束
	apiv1 := r.Group("/api/v1")
	apiv1.Use(middleware.JWT()) //middleware.JWT()
	{
		// 创建标签
		apiv1.POST("/tags", tag.Create)
		// 删除指定标签
		apiv1.DELETE("/tags/:id", tag.Delete)
		// 更新指定标签
		apiv1.PUT("/tags/:id", tag.Update)
		// 获取标签列表
		apiv1.GET("/tags", tag.List)

		// 创建文章
		apiv1.POST("/articles", article.Create)
		// 删除指定文章
		apiv1.DELETE("/articles/:id", article.Delete)
		// 更新指定文章
		apiv1.PUT("/articles/:id", article.Update)
		// 获取指定文章
		apiv1.GET("/articles/:id", article.Get)
		// 获取文章列表
		apiv1.GET("/articles", article.List)
	}

	return r
}
