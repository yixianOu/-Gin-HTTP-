package middleware

//用于编写针对 validator 的语言包翻译的相关功能
import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	"github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func Translations() gin.HandlerFunc {
	return func(c *gin.Context) {
		//New 一组翻译器
		uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
		//看报文头部中的语言信息
		locale := c.GetHeader("locale")
		//根据信息选择翻译器
		trans, _ := uni.GetTranslator(locale)
		//注册校验器
		v, ok := binding.Validator.Engine().(*validator.Validate)
		//给校验器注册翻译器，根据信息选择不同的注册函数
		if ok {
			switch locale {
			case "zh":
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
				break
			case "en":
				_ = en_translations.RegisterDefaultTranslations(v, trans)
				break
			default:
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
				break
			}
			c.Set("trans", trans)
		}

		c.Next()
	}
}
