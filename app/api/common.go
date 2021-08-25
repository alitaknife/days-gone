package api

import (
	"days-gone/app/model"
	"days-gone/app/service"
	"days-gone/library/response"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
)

var CommonApi = &common{}

type common struct {}

func (c *common) Weather(r *ghttp.Request) {
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
