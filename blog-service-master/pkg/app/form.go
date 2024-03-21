package app

//针对入参校验的方法进行了二次封装
import (
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
)

type ValidError struct {
	Key     string
	Message string
}

type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Message
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

// BindAndValid 方法中，通过 ShouldBind 进行参数绑定和入参校验，
// 发生错误后，使用Translator 来对错误消息体进行具体的翻译行为
func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := c.ShouldBind(v)
	//如果有错误，则将错误进行翻译
	if err != nil {
		//从翻译器组拿到特定翻译器
		v := c.Value("trans")
		trans, _ := v.(ut.Translator)
		//判断可否翻译
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			return false, errs
		}
		//使用翻译器，逐个翻译错误信息
		for key, value := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}

		return false, errs
	}

	return true, nil
}
