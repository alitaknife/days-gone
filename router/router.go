package router

import (
	"days-gone/app/api"
	"days-gone/utils"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()
	baseUrl := g.Config().GetString("BaseUrl")
	s.Group("/", func(group *ghttp.RouterGroup) {
		// 无需鉴权
		group.POST("/user/sign-up", api.User.SignUp)
		group.POST("/common/weather", api.CommonApi.Weather)
		// 启动 auth
		utils.Auth = &gtoken.GfToken{
			CacheMode:        g.Cfg().GetInt8("gToken.CacheMode"),
			CacheKey:         g.Cfg().GetString("gToken.CacheKey"),
			Timeout:          g.Cfg().GetInt("gToken.Timeout"),
			MaxRefresh:       g.Cfg().GetInt("gToken.MaxRefresh"),
			TokenDelimiter:   g.Cfg().GetString("gToken.TokenDelimiter"),
			EncryptKey:       g.Cfg().GetBytes("gToken.EncryptKey"),
			AuthFailMsg:      g.Cfg().GetString("gToken.AuthFailMsg"),
			MultiLogin:       g.Cfg().GetBool("gToken.MultiLogin"),
			LoginPath:        "/user/sign-in",
			LoginBeforeFunc:  api.User.SignIn,
			LogoutPath:       "/user/log-out",
			LogoutBeforeFunc: api.User.LogOut,
		}
		// 需要认证
		group.Group(baseUrl+"/", func(group *ghttp.RouterGroup) {
			utils.Auth.Middleware(group)
			group.GET("/user/info", api.User.Info)
			group.POST("/user/update-info", api.User.UpdateInfo)

			group.POST("/file/fast-upload", api.File.FastUpload)
			group.POST("/file/upload", api.File.Upload)
			group.POST("/file/list", api.File.List)
			group.POST("/file/update", api.File.Update)
			group.GET("/file/delete", api.File.Delete)
			group.GET("/file/download", api.File.Download)

			group.POST("/user-file/list", api.UserFile.List)
			group.POST("/user-file/update", api.UserFile.Update)
			group.GET("/user-file/delete", api.UserFile.Delete)
			group.GET("/user-file/download", api.UserFile.Download)
		})
	})
}
