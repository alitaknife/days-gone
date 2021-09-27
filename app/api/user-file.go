package api

import (
	"days-gone/app/model"
	"days-gone/app/service"
	"days-gone/library/response"
	"days-gone/utils"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var UserFile = &userFileApi{}

type userFileApi struct {}

func (f *userFileApi) List(r *ghttp.Request) {
	var userFileListReq *model.UserFileListReq
	if err := r.Parse(&userFileListReq); err != nil {
		if v, ok := err.(gvalid.Error); ok{
			panic(gerror.NewCode(52001, v.FirstString()))
		}
		panic(gerror.NewCode(52001, "parse error"))
	}

	if userFileList, total, err := service.UserFile.List(r, userFileListReq); err != nil {
		panic(gerror.WrapCode(52001, err, err.Error()))
	} else {
		response.SucResp(r).SetData(model.PageRes{
			List: userFileList, Total: total, Current: userFileListReq.Current, Size: userFileListReq.Size,
		}).JsonExit()
	}
}

func (f *userFileApi) Update(r *ghttp.Request) {
	var fileUpdateReq *model.FileUpdateReq
	if err := r.Parse(&fileUpdateReq); err != nil {
		if v, ok := err.(gvalid.Error); ok{
			panic(gerror.NewCode(52002, v.FirstString()))
		}
		panic(gerror.NewCode(52002, "parse error"))
	}

	if err := service.UserFile.Update(r, fileUpdateReq); err != nil {
		panic(gerror.WrapCode(52002, err, err.Error()))
	} else {
		response.SucResp(r).JsonExit()
	}
}

func (f *userFileApi) Delete(r *ghttp.Request) {
	id, err := utils.ValidId(r) // valid id
	if err != nil {
		panic(gerror.WrapCode(52003, err, err.Error()))
	} else {
		if err := service.UserFile.Delete(r, id); err != nil {
			panic(gerror.WrapCode(52003, err, err.Error()))
		} else {
			response.SucResp(r).JsonExit()
		}
	}
}

func (f *userFileApi) Download(r *ghttp.Request) {
	id, err := utils.ValidId(r) // valid id
	if err != nil {
		panic(gerror.WrapCode(52004, err, err.Error()))
	} else {
		res, err := service.UserFile.Download(r, id)
		if err != nil {
			panic(gerror.WrapCode(52004, err, err.Error()))
		}
		if res == "" {
			response.ErrorResp(r).JsonExit()
		} else {
			r.Response.ServeFileDownload(res)
		}
	}
}

// UsedCap get the capacity of cloud zone used
func (f *userFileApi) UsedCap(r *ghttp.Request)  {
	sizeP, err := service.UserFile.UsedCap(r)
	if err != nil {
		panic(gerror.WrapCode(52005, err, err.Error()))
	}
	response.SucResp(r).SetData(sizeP).JsonExit()
}

// FilesType get user`s all file`s types
func (f *userFileApi) FilesType(r *ghttp.Request)  {
	res, err := service.UserFile.FilesType(r)
	if err != nil {
		panic(gerror.WrapCode(52006, err, err.Error()))
	}
	response.SucResp(r).SetData(res.List()).JsonExit()
}

// UploadDays get the number of files uploaded in the last month
func (f *userFileApi) UploadDays(r *ghttp.Request)  {
	res, err := service.UserFile.UploadDays(r)
	if err != nil {
		panic(gerror.WrapCode(52007, err, err.Error()))
	}
	response.SucResp(r).SetData(res.List()).JsonExit()

}
