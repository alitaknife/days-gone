package api

import (
	"days-gone/app/model"
	"days-gone/app/service"
	"days-gone/library/response"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
)

var Common = &commonApi{}

type commonApi struct {}

func (c *commonApi) Weather(r *ghttp.Request) {
	var location *model.Location
	err := r.Parse(&location)
	if err != nil {
		response.JsonErrStrExit(r, gerror.Current(err).Error())
	}
	weather := service.Common.Weather(r, location)
	if weather != nil {
		response.JsonSucExit(r, response.SuccessFirst, weather)
		return
	}
	response.JsonErrStrExit(r, "获取天气失败!")
}

func (c *commonApi) ToBase64(r *ghttp.Request)  {
	url, ok := r.Get("url").(string)
	if ok {
		base64 := service.Common.ToBase64(r, url)
		response.JsonSucExit(r, response.SuccessFirst, base64)
	} else {
		response.JsonErrStrExit(r, "图片转换失败!")
	}
}
