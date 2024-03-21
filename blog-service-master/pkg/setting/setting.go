package setting

//针对读取配置的行为进行封装，便于应用程序的使用

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Setting 为了完成文件配置的读取，借助第三方开源库 viper
type Setting struct {
	vp *viper.Viper
}

// NewSetting 用于初始化本项目的配置的基础属性如：文件名，类型，相对路径
func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	for _, config := range configs {
		if config != "" {
			vp.AddConfigPath(config)
		}
	}
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	s := &Setting{vp}
	s.WatchSettingChange()
	return s, nil
}

// WatchSettingChange 运行时监视配置文件，若改变则重新加载配置属性
func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.ReloadAllSection()
		})
	}()
}
