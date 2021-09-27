package router

import (
	"days-gone/app/api"
	"days-gone/library/response"
	"days-gone/utils"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func MiddlewareEnd(r *ghttp.Request) {
	r.Middleware.Next()
	// 后置中间件错误拦截处理
	if err := r.GetError(); err != nil {
		r.Response.ClearBuffer()
		r.Response.Status = http.StatusOK
		response.ErrorResp(r).SetCode(int32(gerror.Code(err))).SetMsg(gerror.Current(err).Error()).JsonExit()
	}
}

func init() {
	s := g.Server()
	baseUrl := g.Config().GetString("BaseUrl")
	s.Group(baseUrl+"/", func(group *ghttp.RouterGroup) {
		// 允许跨域
		group.Middleware(MiddlewareCORS, MiddlewareEnd)
		// 无需鉴权
		group.POST("/user/sign-up", api.User.SignUp)
		group.POST("/common/weather", api.Common.Weather)
		// 需要认证
		group.Group("/", func(group *ghttp.RouterGroup) {
			utils.Auth.LoginBeforeFunc = api.User.SignIn
			utils.Auth.LogoutBeforeFunc = api.User.LogOut
			utils.Auth.Middleware(group)
			group.GET("/user/info", api.User.Info)
			group.POST("/user/upload-avatar", api.User.UploadAvatar)
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
			group.GET("/user-file/used-cap", api.UserFile.UsedCap)
			group.GET("/user-file/files-type", api.UserFile.FilesType)
			group.GET("/user-file/upload-days", api.UserFile.UploadDays)

			group.GET("/common/pic-base64", api.Common.ToBase64)
		})
	})
}
