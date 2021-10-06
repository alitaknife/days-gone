package service

import (
	"days-gone/app/dao"
	"days-gone/app/model"
	"days-gone/utils"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/guid"
)

var User = userService{}

type userService struct{}

func (u *userService) SignUp(r *ghttp.Request, req *model.UserSignUpReq) error {
	db := dao.User.Ctx(r.GetCtx())
	salt := guid.S()
	ps, err := gmd5.EncryptString(req.UserPwd + salt)
	if err != nil {
		return gerror.Wrap(err, "generate salt error")
	}
	userInfo := &model.User{
		UserName:     req.UserName,
		UserNickname: req.UserName,
		UserPwd:      ps,
		Salt:         salt,
		SignUpAt:     gtime.Now(),
	}
	ra, err := db.Insert(userInfo)
	if err != nil {
		return gerror.New("insert error: please check your name and password")
	}
	if rows, err := ra.RowsAffected(); err != nil {
		return gerror.New("insert error")
	} else if rows == 0 {
		return gerror.New("insert error: no rows founded")
	}
	return nil
}

func (u *userService) SignIn(r *ghttp.Request, req *model.UserSignInReq) (user *model.User, err error) {
	db := dao.User.Ctx(r.GetCtx())
	err = db.Where(g.Map{"user_name": req.UserName}).Scan(&user)
	if err != nil {
		return nil, gerror.New("search error")
	}
	if user == nil {
		return user, gerror.New("cannot find this user")
	}

	// check password
	if ps, err := gmd5.EncryptString(req.UserPwd + user.Salt); err != nil {
		return nil, gerror.New("login failed")
	} else if ps != user.UserPwd {
		return nil, gerror.New("wrong password")
	} else {
		return user, nil
	}
}

func (u *userService) UpdateInfo(r *ghttp.Request, req *model.UserInfoReq) error {
	db := dao.User.Ctx(r.Context())
	res, err := db.Data(*req).Where("user_name", u.GetCacheUserName(r)).Update()
	if err != nil {
		return gerror.New("update failed")
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return gerror.New("update failed")
	}
	return nil
}

func (u *userService) UploadAvatar(r *ghttp.Request, avatar *model.Avatar) (string, error) {
	var baseUrl string = g.Config().GetString("gitBed.BASE_URL")
	var owner string = g.Config().GetString("gitBed.OWNER")
	var repo string = g.Config().GetString("gitBed.REPO")
	var path string = g.Config().GetString("gitBed.PATH")
	var picName string = gtime.Now().Format("Y/m/d/His")

	url := baseUrl + owner + "/" + repo + "/contents/" + path + "/" + picName + ".png" // build url
	res := g.Client().ContentType("multipart/form-data").PostContent(url, g.Map{
		"access_token": g.Config().GetString("gitBed.ACCESS_TOKEN"),
		"content": gstr.StrEx(avatar.Avatar, ","),
		"message": g.Config().GetString("gitBed.MSG"),
		"branch": g.Config().GetString("gitBed.BRANCH"),
		})

	if j, err := gjson.DecodeToJson(res); err != nil {
		return "", gerror.Wrap(err, "upload avatar failed")
	} else {
		avatarUrl := j.GetString("content.download_url")
		if avatarUrl == ""{
			return "", gerror.Wrap(err, "upload avatar failed")
		}
		// if update avatar to pic bed succeed then update it to db
		db := dao.User.Ctx(r.GetCtx())
		result, err := db.Data(g.Map{"avatar": avatarUrl}).Where("user_name", u.GetCacheUserName(r)).Update()
		if err != nil {
			return "", gerror.New("upload avatar failed")
		}
		rows, err := result.RowsAffected()
		if err != nil || rows == 0 {
			return "", gerror.New("upload avatar failed")
		}
		return avatarUrl, nil
	}
}

func (u *userService) LogOut(r *ghttp.Request) bool {
	db := dao.User.Ctx(r.GetCtx())
	userName := u.GetCacheUserName(r)
	_, err := db.Data(g.Map{"last_active": gtime.Now()}).Where("user_name", userName).Update()
	if err != nil {
		return false
	}
	return true
}

// GetCacheUserName get user name in cache
func (u *userService) GetCacheUserName(r *ghttp.Request) string {
	res := utils.Auth.GetTokenData(r)
	return res.GetString("userKey")
}

// GetUserInfo get user info from db
func (u *userService) GetUserInfo(r *ghttp.Request) (*model.UserSignInRes, error) {
	db := dao.User.Ctx(r.GetCtx())
	userInfo := &model.UserSignInRes{}
	err := db.Where("user_name", u.GetCacheUserName(r)).Scan(&userInfo)
	if err != nil {
		return nil, gerror.New("get user info failed")
	}
	if userInfo == nil {
		return userInfo, gerror.New("cannot find this user")
	}
	return userInfo, nil
}
