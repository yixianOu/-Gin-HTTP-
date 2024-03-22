package service

//定义服务结构体，并用context和数据库engine实例化一个服务
import (
	"context"

	otgorm "github.com/eddycjy/opentracing-gorm"

	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/internal/dao"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	//WithContext根据传入的ctx对DBEngine进行设置，返回设定后的DBEngine
	svc.dao = dao.New(otgorm.WithContext(svc.ctx, global.DBEngine))
	return svc
}
