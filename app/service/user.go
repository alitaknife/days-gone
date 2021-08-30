package service

import (
	"days-gone/app/dao"
	"days-gone/app/model"
	"days-gone/utils"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/guid"
)

var User = userService{}

type userService struct{}

func (u *userService) SignUp(r *ghttp.Request, req *model.UserSignUpReq) error {
	db := dao.User.Ctx(r.GetCtx())
	salt := guid.S()
	ps, err := gmd5.EncryptString(req.UserPwd + salt)
	if err != nil {
		return err
	}
	userInfo := &model.User{
		UserName:     req.UserName,
		UserNickname: req.UserName,
		UserPwd:      ps,
		Salt:         salt,
		SignUpAt:     gtime.Now(),
	}
	_, err = db.Insert(userInfo)
	return err
}

func (u *userService) SignIn(r *ghttp.Request, req *model.UserSignInReq) (user *model.User, err error) {
	db := dao.User.Ctx(r.GetCtx())
	err = db.Where(g.Map{"user_name": req.UserName}).Scan(&user)
	if err != nil {
		return user, err
	}
	// 密码判断
	if ps, err := gmd5.EncryptString(req.UserPwd + user.Salt); err != nil || ps != user.UserPwd {
		return nil, err
	} else {
		return user, err
	}
}

func (u *userService) UpdateInfo(r *ghttp.Request, req *model.UserInfoReq) error {
	db := dao.User.Ctx(r.Context())
	userCache := u.GetCacheUserInfo(r)
	_ = gconv.Struct(req, userCache)
	res, err := db.Data(*userCache).Where("user_name", userCache.UserName).Update()
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) LogOut(r *ghttp.Request) bool {
	db := dao.User.Ctx(r.GetCtx())
	userName := u.GetCacheUserInfo(r).UserName
	_, err := db.Data(g.Map{"last_active": gtime.Now()}).Where("user_name", userName).Update()
	if err != nil {
		return false
	}
	return true
}

// GetCacheUserInfo 获取缓存用户信息
func (u *userService) GetCacheUserInfo(r *ghttp.Request) *model.UserSignInRes {
	res := utils.Auth.GetTokenData(r)
	userCache := &model.UserSignInRes{}
	_ = gconv.Struct(res.Get("data"), userCache)
	return userCache
}
