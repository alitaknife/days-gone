package service

import (
	"days-gone/app/model"
	"github.com/gogf/gf/encoding/gbase64"
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

	urlStr := gurl.BuildQuery(u) // parameter splicing
	baseUrl := "http://weathernew.pae.baidu.com/weathernew/pc?"
	res := g.Client().GetContent(baseUrl + urlStr)
	// 从返回的 HTML 模板中提取天气信息
	temp := gstr.StrTillEx(res, ";") // ; The first time it appears is the end of the json message
	final := gstr.Str(temp, "{") // { The first time it appears is the start of the json message
	jsonObj, _ := gjson.DecodeToJson(final)
	return jsonObj
}

func (c *commonService) ToBase64(r *ghttp.Request, url string) string {
	var picType = gstr.SubStr(url, gstr.PosR(url, ".") + 1)
	res := g.Client().GetContent(url)
	base64 := gbase64.EncodeToString([]byte(res))
	return "data:image/" + picType + ";base64," + base64
}
