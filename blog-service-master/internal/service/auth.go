package service

import "errors"

// AuthRequest 接口入参的校验
type AuthRequest struct {
	AppKey    string `form:"app_key" binding:"required"`
	AppSecret string `form:"app_secret" binding:"required"`
}

// CheckAuth 使用客户端所传入的认证信息作为筛选条件获取数据行，
// 以此根据是否取到认证信息 ID 来进行是否存在的判定
func (svc *Service) CheckAuth(param *AuthRequest) error {
	auth, err := svc.dao.GetAuth(
		param.AppKey,
		param.AppSecret,
	)
	if err != nil {
		return err
	}

	if auth.ID > 0 {
		return nil
	}

	return errors.New("auth info does not exist.")
}
