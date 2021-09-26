package response

import (
	"days-gone/app/model"
	"github.com/gogf/gf/net/ghttp"
)

type ApiResp struct {
	r *ghttp.Request
	c *model.CommonRes
}

func SucResp(r *ghttp.Request) *ApiResp {
	c := model.CommonRes{
		Code: 20000,
		Msg: "operation succeeded！",
		Data: "",
	}
	return &ApiResp{
		r: r,
		c: &c,
	}
}

func ErrorResp(r *ghttp.Request) *ApiResp {
	c := model.CommonRes{
		Code:  50000,
		Msg:   "operation failed！",
		Data: "",
	}
	return &ApiResp{
		r: r,
		c: &c,
	}
}

func (a *ApiResp) SetCode(code int32) *ApiResp {
	a.c.Code = code
	return a
}

func (a *ApiResp) SetMsg(msg string) *ApiResp {
	a.c.Msg = msg
	return a
}

func (a *ApiResp) SetData(data interface{}) *ApiResp {
	a.c.Data = data
	return a
}

func (a *ApiResp) JsonExit() {
	a.r.Response.WriteJsonExit(a.c)
}



