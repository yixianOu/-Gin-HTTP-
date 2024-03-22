package app

//根据特定的鉴权场景对jwt-go库进行设计,组合其提供的 API
import (
	"time"

	"github.com/go-programming-tour-book/blog-service/pkg/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-programming-tour-book/blog-service/global"
)

// Claims JWT 的一些基本属性
type Claims struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	jwt.StandardClaims
}

// GetJWTSecret 获取该项目的 JWT Secret
func GetJWTSecret() []byte {
	return []byte(global.JWTSetting.Secret)
}

// GenerateToken 生成 JWT Token
func GenerateToken(appKey, appSecret string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(global.JWTSetting.Expire)
	claims := Claims{
		AppKey:    util.EncodeMD5(appKey),
		AppSecret: util.EncodeMD5(appSecret),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    global.JWTSetting.Issuer,
		},
	}
	//jwt.NewWithClaims根据 Claims 结构体创建 Token 实例
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//SignedString根据所传入 Secret 不同，进行签名并返回标准的 Token。
	token, err := tokenClaims.SignedString(GetJWTSecret())
	return token, err
}

// ParseToken 解析和校验 Token
func ParseToken(token string) (*Claims, error) {
	//解析鉴权的声明，方法内部是具体的解码和校验过程，返回 *Token。
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Claims)
		//验证基于时间的声明
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
