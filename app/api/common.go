package api

import (
	"days-gone/app/model"
	"days-gone/app/service"
	"days-gone/library/response"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var Common = &commonApi{}

type commonApi struct {}

func (c *commonApi) Weather(r *ghttp.Request) {
	var location *model.Location
	if err := r.Parse(&location); err != nil {
		if v, ok := err.(gvalid.Error); ok{
			panic(gerror.NewCode(53001, v.FirstString()))
		}
		panic(gerror.NewCode(53001, "parse error"))
	}

	weather := service.Common.Weather(r, location)
	if weather == nil{
		panic(gerror.NewCode(53001, "get the weather failed"))
	}
	response.SucResp(r).SetData(weather).JsonExit()
}

func (c *commonApi) ToBase64(r *ghttp.Request)  {
	url, ok := r.Get("url").(string)
	if ok {
		base64 := service.Common.ToBase64(r, url)
		response.SucResp(r).SetData(base64).JsonExit()
	} else {
		panic(gerror.NewCode(53002, "conversion failed"))
	}
}
