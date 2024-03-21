package limiter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// LimiterIface 接口，用于定义当前限流器所必须要的方法
type LimiterIface interface {
	Key(c *gin.Context) string
	GetBucket(key string) (*ratelimit.Bucket, bool)
	AddBuckets(rules ...LimiterBucketRule) LimiterIface
}

// Limiter 结构体用于存储令牌桶与键值对名称的映射关系
type Limiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
}

// LimiterBucketRule 结构体用于存储令牌桶的一些相应规则属性
type LimiterBucketRule struct {
	Key          string
	FillInterval time.Duration
	Capacity     int64
	Quantum      int64
}
