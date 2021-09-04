package api

import (
	"days-gone/app/model"
	"days-gone/app/service"
	"days-gone/library/response"
	"github.com/gogf/gf/net/ghttp"
)

var User = userApi{}

type userApi struct{}

func (u *userApi) SignUp(r *ghttp.Request) {
	var userSignUpReq *model.UserSignUpReq
	if err := r.Parse(&userSignUpReq); err != nil {
		response.JsonErrStrExit(r, err.Error())
		return
	}

	// 用户已存在
	if err := service.User.SignUp(r, userSignUpReq); err != nil {
		response.JsonErrExit(r, response.ErrorUserArdExist)
		return
	}
	response.JsonSucExit(r, response.SuccessSignUp)
}

func (u *userApi) SignIn(r *ghttp.Request) (string, interface{}) {
	var userSignInReq *model.UserSignInReq
	if err := r.Parse(&userSignInReq); err != nil {
		response.JsonErrStrExit(r, err.Error())
	}
	user, err := service.User.SignIn(r, userSignInReq)
	// 登录出错
	if err != nil {
		response.JsonErrExit(r, response.ErrorSignIn)
	}
	if user == nil {
		// 用户不存在
		response.JsonErrExit(r, response.ErrorSignInNoFind)
		// 返回空字符表示 gToken 登录失败
		return "", nil
	}
	// 通过登录
	return user.UserName, user
}

// Info 获取用户信息
func (u *userApi) Info(r *ghttp.Request) {
	userCache := service.User.GetCacheUserInfo(r)
	if userCache != nil {
		response.JsonSucExit(r, response.SuccessUserInfo, userCache)
	}
	response.JsonErrExit(r, response.ErrorUserInfo)
}

// UpdateInfo 更新用户信息
func (u *userApi) UpdateInfo(r *ghttp.Request) {
	var userInfoReq *model.UserInfoReq
	err := r.Parse(&userInfoReq)
	if err != nil {
		response.JsonErrStrExit(r, err.Error())
		return
	}

	if err := service.User.UpdateInfo(r, userInfoReq); err != nil {
		response.JsonErrExit(r, response.ErrorUpdated)
	} else {
		response.JsonSucExit(r, response.SuccessUpdated)
	}
}

func (u *userApi) UploadAvatar(r *ghttp.Request) {
	var avatar *model.Avatar
	err := r.Parse(&avatar)
	if err != nil {
		response.JsonErrStrExit(r, err.Error())
		return
	}
	avatarUrl, err := service.User.UploadAvatar(r, avatar)
	if err != nil {
		response.JsonErrExit(r, response.ErrorNoFileUpload)
		return
	}
	if avatarUrl == "" {
		response.JsonErrExit(r, response.ErrorNoFileUpload)
	} else {
		response.JsonSucExit(r, response.SuccessUpdated, avatarUrl)
	}
}

// LogOut 登出之前调用
func (u *userApi) LogOut(r *ghttp.Request) bool {
	isSuc := service.User.LogOut(r)
	return isSuc
}
