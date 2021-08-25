package service

import (
	"days-gone/app/model"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/encoding/gurl"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/text/gstr"
	"net/url"
)

var Common = &commonService{}

type commonService struct {}

func (c *commonService) Weather(r *ghttp.Request, loc *model.Location) *gjson.Json {
	u := url.Values{}
	u.Add("query", loc.City + "天气")
	u.Add("srcid", "4982")
	u.Add("city_name", loc.City)
	u.Add("province_name", loc.Province)

	urlStr := gurl.BuildQuery(u) // 参数拼接
	baseUrl := "http://weathernew.pae.baidu.com/weathernew/pc?"
	res := g.Client().GetContent(baseUrl + urlStr)
	// 从返回的 HTML 模板中提取天气信息
	temp := gstr.StrTillEx(res, ";") // ; 第一次出现是 json 信息结尾
	final := gstr.Str(temp, "{") // { 第一次出现是 json 信息开头
	jsonObj, _ := gjson.DecodeToJson(final)
	return jsonObj
}
