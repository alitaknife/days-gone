package api

import (
	"days-gone/app/model"
	"days-gone/app/service"
	"days-gone/library/response"
	"days-gone/utils"
	"github.com/gogf/gf/net/ghttp"
)

var UserFile = &userFileApi{}

type userFileApi struct {}

func (f *userFileApi) List(r *ghttp.Request) {
	var userFileListReq *model.UserFileListReq
	if err := r.Parse(&userFileListReq); err != nil {
		response.JsonErrStrExit(r, err.Error())
		return
	}

	if userFileList, total, err := service.UserFile.List(r, userFileListReq); err != nil {
		response.JsonErrExit(r, response.ErrorGetList)
	} else {
		response.JsonSucExit(r, response.SuccessGetList, response.PageResponse{List: userFileList, Total: total, Current: userFileListReq.Current, Size: userFileListReq.Size})
	}
}

func (f *userFileApi) Update(r *ghttp.Request) {
	var fileUpdateReq *model.FileUpdateReq
	if err := r.Parse(&fileUpdateReq); err != nil {
		response.JsonErrStrExit(r, err.Error())
		return
	}

	if err := service.UserFile.Update(r, fileUpdateReq); err != nil {
		response.JsonErrExit(r, response.ErrorUpdated)
	} else {
		response.JsonSucExit(r, response.SuccessUpdated)
	}
}

func (f *userFileApi) Delete(r *ghttp.Request) {
	id, err := utils.ValidId(r) // 校验 Id 合法性
	if err != nil {
		response.JsonErrStrExit(r, err.Error())
	} else {
		if err := service.UserFile.Delete(r, id); err != nil {
			response.JsonErrExit(r, response.ErrorDeleted)
		} else {
			response.JsonSucExit(r, response.SuccessDeleted)
		}
	}
}

func (f *userFileApi) Download(r *ghttp.Request) {
	id, err := utils.ValidId(r) // 校验 Id 合法性
	if err != nil {
		response.JsonErrStrExit(r, err.Error())
	} else {
		res := service.UserFile.Download(r, id)
		if res != "" {
			r.Response.ServeFileDownload(res)
		} else {
			response.JsonErrExit(r, response.ErrorDownload)
		}
	}
}

// UsedCap 获取用户网盘已使用容量
func (f *userFileApi) UsedCap(r *ghttp.Request)  {
	sizeP, err := service.UserFile.UsedCap(r)
	if err != nil {
		response.JsonErrExit(r, response.ErrorGetCap)
	}
	if sizeP == 0 {
		response.JsonErrExit(r, response.ErrorGetCap)
	}
	response.JsonSucExit(r, response.SuccessGetCap, sizeP)
}

// FilesType 获取用户所有的文件类型以及数目
func (f *userFileApi) FilesType(r *ghttp.Request)  {
	res, err := service.UserFile.FilesType(r)
	if err != nil || res == nil{
		response.JsonErrExit(r, response.ErrorGetType)
	} else {
		response.JsonSucExit(r, response.SuccessGetType, res.List())
	}
}

// UploadDays 获取最近一个月上传的文件数
func (f *userFileApi) UploadDays(r *ghttp.Request)  {
	res, err := service.UserFile.UploadDays(r)
	if err != nil || res == nil{
		response.JsonErrExit(r, response.ErrorGetUpDays)
	} else {
		response.JsonSucExit(r, response.SuccessGetUpDays, res.List())
	}

}
