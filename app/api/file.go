package api

import (
	"days-gone/app/model"
	"days-gone/app/service"
	"days-gone/library/response"
	"days-gone/utils"
	"github.com/gogf/gf/net/ghttp"
)

var File = &fileApi{}

type fileApi struct{}

func (f *fileApi) Upload(r *ghttp.Request) {
	file := r.GetUploadFile("upload-file")
	if file != nil {
		code, err := service.File.Upload(r, file)
		if err != nil {
			response.JsonErrExit(r, code)
			return
		}
		response.JsonSucExit(r, code)
	}
	response.JsonErrExit(r, response.ErrorNoFileUpload)
}

func (f *fileApi) List(r *ghttp.Request) {
	var fileListReq *model.FileListReq
	if err := r.Parse(&fileListReq); err != nil {
		response.JsonErrStrExit(r, err.Error())
		return
	}

	if fileList, total, err := service.File.List(r, fileListReq); err != nil {
		response.JsonErrExit(r, response.ErrorGetList)
	} else {
		response.JsonSucExit(r, response.SuccessGetList, response.PageResponse{List: fileList, Total: total, Current: fileListReq.Current, Size: fileListReq.Size})
	}
}

func (f *fileApi) Update(r *ghttp.Request) {
	var fileUpdateReq *model.FileUpdateReq
	if err := r.Parse(&fileUpdateReq); err != nil {
		response.JsonErrStrExit(r, err.Error())
		return
	}

	if err := service.File.Update(r, fileUpdateReq); err != nil {
		response.JsonErrExit(r, response.ErrorUpdated)
	} else {
		response.JsonSucExit(r, response.SuccessUpdated)
	}
}

func (f *fileApi) Delete(r *ghttp.Request) {
	id, err := utils.ValidId(r) // 校验 Id 合法性
	if err != nil {
		response.JsonErrStrExit(r, err.Error())
	} else {
		if err := service.File.Delete(r, id); err != nil {
			response.JsonErrExit(r, response.ErrorDeleted)
		} else {
			response.JsonSucExit(r, response.SuccessDeleted)
		}
	}
}

func (f *fileApi) Download(r *ghttp.Request) {
	id, err := utils.ValidId(r) // 校验 Id 合法性
	if err != nil {
		response.JsonErrStrExit(r, err.Error())
	} else {
		res, err := service.File.Download(r, id)
		if err != nil || res.IsEmpty() {
			response.JsonErrExit(r, response.ErrorDownload)
		} else {
			r.Response.ServeFileDownload(res.String())
		}
	}
}
