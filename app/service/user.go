package service

import (
	"days-gone/app/dao"
	"days-gone/app/model"
	"days-gone/utils"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
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
	if err != nil || user == nil{
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

func (u *userService) UploadAvatar(r *ghttp.Request, avatar *model.Avatar) (string, error) {
	var baseUrl string = g.Config().GetString("gitBed.BASE_URL")
	var owner string = g.Config().GetString("gitBed.OWNER")
	var repo string = g.Config().GetString("gitBed.REPO")
	var path string = g.Config().GetString("gitBed.PATH")
	var picName string = gtime.Now().Format("Y/m/d/His")

	url := baseUrl + owner + "/" + repo + "/contents/" + path + "/" + picName + ".png" // url 构建
	res := g.Client().ContentType("multipart/form-data").PostContent(url, g.Map{
		"access_token": g.Config().GetString("gitBed.ACCESS_TOKEN"),
		"content": gstr.StrEx(avatar.Avatar, ","),
		"message": g.Config().GetString("gitBed.MSG"),
		"branch": g.Config().GetString("gitBed.BRANCH"),
		})

	if j, err := gjson.DecodeToJson(res); err != nil {
		return "", err
	} else {
		avatarUrl := j.GetString("content.download_url")
		if avatarUrl == ""{
			return "", err
		}
		// 头像上传成功就更新到数据库
		db := dao.User.Ctx(r.GetCtx())
		result, err := db.Data(g.Map{"avatar": avatarUrl}).Where("user_name", u.GetCacheUserInfo(r).UserName).Update()
		if err != nil {
			return "", err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return "", err
		}
		if rows == 0 {
			return "", nil
		}
		return avatarUrl, nil
	}
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
