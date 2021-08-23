package response

import "github.com/gogf/gf/net/ghttp"

var (
	SUCCESS      = 0    // 正常
	FAIL         = -1   // 失败
	ERROR        = -99  // 异常
	UNAUTHORIZED = -401 // 未认证
)

type JsonResponse struct {
	Code int         `json:"code"` // 错误码((0:成功, 1:失败, >1:错误码))
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// PageResponse 分页返回结构体
type PageResponse struct {
	List    interface{} `json:"list"`
	Total   int         `json:"total"`
	Current int         `json:"current"`
	Size    int         `json:"size"`
}

func Json(r *ghttp.Request, code int, msg string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	} else {
		responseData = ""
	}
	_ = r.Response.WriteJson(JsonResponse{
		Code: code,
		Msg:  msg,
		Data: responseData,
	})
}

func JsonSucExit(r *ghttp.Request, msgCode Code, data ...interface{}) {
	Json(r, SUCCESS, msgCode.Message(), data...)
	r.Exit()
}

func JsonErrExit(r *ghttp.Request, msgCode Code, data ...interface{}) {
	Json(r, FAIL, msgCode.Message(), data...)
	r.Exit()
}

func JsonSucStrExit(r *ghttp.Request, msg string, data ...interface{}) {
	Json(r, SUCCESS, msg, data...)
	r.Exit()
}

func JsonErrStrExit(r *ghttp.Request, msg string, data ...interface{}) {
	Json(r, FAIL, msg, data...)
	r.Exit()
}
