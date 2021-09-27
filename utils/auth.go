package utils

import (
	"days-gone/library/response"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var Auth = (*gtoken.GfToken)(nil)

func init() {
	Auth =  &gtoken.GfToken{
		CacheMode:        g.Cfg().GetInt8("gToken.CacheMode"),
		CacheKey:         g.Cfg().GetString("gToken.CacheKey"),
		Timeout:          g.Cfg().GetInt("gToken.Timeout"),
		MaxRefresh:       g.Cfg().GetInt("gToken.MaxRefresh"),
		TokenDelimiter:   g.Cfg().GetString("gToken.TokenDelimiter"),
		EncryptKey:       g.Cfg().GetBytes("gToken.EncryptKey"),
		AuthFailMsg:      g.Cfg().GetString("gToken.AuthFailMsg"),
		MultiLogin:       g.Cfg().GetBool("gToken.MultiLogin"),
		LoginPath:        "/user/sign-in",
		LogoutPath:       "/user/log-out",
		LoginAfterFunc: formatResp,
		AuthAfterFunc: formatResp,
		LogoutAfterFunc: formatResp,
	}
}

// Encapsulate the return result of auth
var formatResp = func(r *ghttp.Request, res gtoken.Resp) {
	switch res.Code {
	case gtoken.SUCCESS:
		token := res.GetString("token")
		if token != "" {
			response.SucResp(r).SetCode(20002).SetData(g.Map{"token": token}).JsonExit()
		} else {
			r.Middleware.Next()
		}
	case gtoken.ERROR, gtoken.UNAUTHORIZED:
		response.ErrorResp(r).SetCode(50000).SetMsg("The current token is not authenticated").JsonExit()
	case gtoken.FAIL:
		response.ErrorResp(r).SetCode(50000).SetMsg(res.Msg).JsonExit()
	default:
		response.ErrorResp(r).JsonExit()
	}
}


