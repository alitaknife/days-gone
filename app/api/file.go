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

var File = &fileApi{}

type fileApi struct{}

func (f *fileApi) FastUpload(r *ghttp.Request) {
	var fastReq = &model.FastUploadReq{}
	if err := r.Parse(&fastReq); err != nil {
		if v, ok := err.(gvalid.Error); ok{
			panic(gerror.NewCode(51002, v.FirstString()))
		}
		panic(gerror.NewCode(51002, "parse error"))
	}

	err := service.File.FastUpload(r, fastReq.FileSha1, fastReq.FileName)
	if err != nil {
		if gerror.Code(err) > 0 {
			// tell user call the normal upload api
			panic(gerror.WrapCode(31002, err, err.Error()))
		}
		panic(gerror.WrapCode(51002, err, err.Error()))
	}
	response.SucResp(r).JsonExit()
}

func (f *fileApi) Upload(r *ghttp.Request) {
	file := r.GetUploadFile("upload-file")
	if file == nil {
		panic(gerror.NewCode(51003, "please upload a file at least"))
	}

	err := service.File.Upload(r, file)
	if err != nil {
		panic(gerror.NewCode(51003, "upload file failed"))
	}
	response.SucResp(r).JsonExit()
}

func (f *fileApi) List(r *ghttp.Request) {
	var fileListReq *model.FileListReq
	if err := r.Parse(&fileListReq); err != nil {
		if v, ok := err.(gvalid.Error); ok{
			panic(gerror.NewCode(51004, v.FirstString()))
		}
		panic(gerror.NewCode(51004, "parse error"))
	}

	if fileList, total, err := service.File.List(r, fileListReq); err != nil {
		panic(gerror.WrapCode(51004, err, err.Error()))
	} else {
		response.SucResp(r).SetData(model.PageRes{
			List: fileList, Total: total, Current: fileListReq.Current, Size: fileListReq.Size,
		}).JsonExit()
	}
}

func (f *fileApi) Update(r *ghttp.Request) {
	var fileUpdateReq *model.FileUpdateReq
	if err := r.Parse(&fileUpdateReq); err != nil {
		if v, ok := err.(gvalid.Error); ok{
			panic(gerror.NewCode(51005, v.FirstString()))
		}
		panic(gerror.NewCode(51005, "parse error"))
	}

	if err := service.File.Update(r, fileUpdateReq); err != nil {
		panic(gerror.WrapCode(51005, err, err.Error()))
	} else {
		response.SucResp(r).JsonExit()
	}
}

func (f *fileApi) Delete(r *ghttp.Request) {
	id, err := utils.ValidId(r) // valid id
	if err != nil {
		panic(gerror.WrapCode(51006, err, err.Error()))
	} else {
		if err := service.File.Delete(r, id); err != nil {
			panic(gerror.WrapCode(51006, err, err.Error()))
		} else {
			response.SucResp(r).JsonExit()
		}
	}
}

func (f *fileApi) Download(r *ghttp.Request) {
	id, err := utils.ValidId(r) // valid id
	if err != nil {
		panic(gerror.WrapCode(51007, err, err.Error()))
	} else {
		res, err := service.File.Download(r, id)
		if err != nil {
			panic(gerror.WrapCode(51007, err, err.Error()))
		} else {
			r.Response.ServeFileDownload(res.String())
		}
	}
}
