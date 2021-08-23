package api

import (
	"days-gone/app/model"
	"days-gone/app/service"
	"days-gone/library/response"
	"github.com/gogf/gf/net/ghttp"
	"strconv"
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
		response.JsonErrExit(r, response.ErrorParsePram)
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

func (f *fileApi) Download(r *ghttp.Request) {
	id, ok := r.Get("id").(string)
	if ok {
		idInt, _ := strconv.Atoi(id)
		file, err := service.File.Download(r, idInt)
		if err == nil && file != nil {
			r.Response.ServeFileDownload(file.FileAddr)
		} else {
			response.JsonErrExit(r, response.ErrorDownload)
		}
	} else {
		response.JsonErrExit(r, response.ErrorDownload)
	}
}
