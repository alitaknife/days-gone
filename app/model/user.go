package model

import "github.com/gogf/gf/os/gtime"

// UserSignUpReq 用户注册请求结构体
type UserSignUpReq struct {
	UserName string `p:"userName" v:"required|length:4,16#请输入账号|账号长度为:min到:max位"`
	UserPwd  string `p:"userPwd" v:"required|length:6,20#请输入密码|密码长度不够"`
	UserPwd2 string `p:"userPwd2" v:"required|length:6,20|same:userPwd#请输入密码|密码长度不够|两次密码不一致"`
}

// UserSignInReq 用户登录请求结构体
type UserSignInReq struct {
	UserName string `p:"userName" v:"required|length:4,16#请输入账号|账号长度为:min到:max位"`
	UserPwd  string `p:"userPwd" v:"required|length:6,20#请输入密码|密码长度不够"`
}

// UserInfoReq 用户更新信息请求结构体
type UserInfoReq struct {
	UserNickname string `p:"userNickname" v:"required|length:1,16#请输入用户昵称|昵称长度为:min到:max位"`
	Sex int `p:"sex" v:"required|between:0,1#请输入合法的状态值"`
	Email string `p:"email" v:"email#请输入合法的邮箱"`
	Phone string `p:"phone" v:"phone#请输入合法的手机号"`
}

// UserSignInRes 用户登录返回结构体
type UserSignInRes struct {
	UserName       string      `json:"userName"`
	UserNickname   string      `json:"userNickname"`
	Avatar         string      `json:"avatar"`
	Sex            int         `json:"sex"`
	Email          string      `json:"email"`
	Phone          string      `json:"phone"`
	EmailValidated int         `json:"emailValidated"`
	PhoneValidated int         `json:"phoneValidated"`
	SignupAt       *gtime.Time `json:"signupAt"`
	LastActive     *gtime.Time `json:"lastActive"`
	Profile        string      `json:"profile"`
	Status         int         `json:"status"`
}
