package api

import (
	"days-gone/app/model"
	"days-gone/app/service"
	"days-gone/library/response"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var User = userApi{}

type userApi struct{}

func (u *userApi) SignUp(r *ghttp.Request) {
	var userSignUpReq *model.UserSignUpReq
	if err := r.Parse(&userSignUpReq); err != nil {
		if v, ok := err.(gvalid.Error); ok{
			panic(gerror.NewCode(50001, v.FirstString()))
		}
		panic(gerror.NewCode(50001, "parse error"))
	}

	// insert user by service
	if err := service.User.SignUp(r, userSignUpReq); err != nil {
		panic(gerror.WrapCode(50001, err, err.Error()))
	}
	response.SucResp(r).JsonExit()
}

func (u *userApi) SignIn(r *ghttp.Request) (string, interface{}) {
	var userSignInReq *model.UserSignInReq
	if err := r.Parse(&userSignInReq); err != nil {
		if v, ok := err.(gvalid.Error); ok{
			response.ErrorResp(r).SetCode(50002).SetMsg(v.FirstString()).JsonExit()
		}
		response.ErrorResp(r).SetCode(50002).SetMsg("parse error").JsonExit()
	}
	user, err := service.User.SignIn(r, userSignInReq)
	// login failed
	if err != nil || user == nil {
		response.ErrorResp(r).SetCode(50002).SetMsg(gerror.Current(err).Error()).JsonExit()
		// return an empty string to set no token
		return "", nil
	}
	// sign in passed
	return user.UserName, user
}

// Info get user`s information
func (u *userApi) Info(r *ghttp.Request) {
	userCache := service.User.GetCacheUserInfo(r)
	if userCache == nil {
		panic(gerror.NewCode(50003, "cannot find user"))
	}
	response.SucResp(r).SetData(userCache).JsonExit()
}

// UpdateInfo update user`s information
func (u *userApi) UpdateInfo(r *ghttp.Request) {
	var userInfoReq *model.UserInfoReq
	// parse the parameter
	if err := r.Parse(&userInfoReq); err != nil {
		if v, ok := err.(gvalid.Error); ok{
			panic(gerror.NewCode(50002, v.FirstString()))
		}
		panic(gerror.NewCode(50002, "parse error"))
	}

	// query user by service
	if err := service.User.UpdateInfo(r, userInfoReq); err != nil {
		panic(gerror.WrapCode(50002, err, err.Error()))
	}
	response.SucResp(r).JsonExit()
}

func (u *userApi) UploadAvatar(r *ghttp.Request) {
	var avatar *model.Avatar
	if err := r.Parse(&avatar); err != nil {
		if v, ok := err.(gvalid.Error); ok{
			panic(gerror.NewCode(50002, v.FirstString()))
		}
		panic(gerror.NewCode(50002, "parse error"))
	}

	avatarUrl, err := service.User.UploadAvatar(r, avatar)
	if err != nil {
		panic(gerror.WrapCode(50002, err, err.Error()))
	}
	response.SucResp(r).SetData(avatarUrl).JsonExit()
}

// LogOut call before logout
func (u *userApi) LogOut(r *ghttp.Request) bool {
	isSuc := service.User.LogOut(r)
	return isSuc
}
