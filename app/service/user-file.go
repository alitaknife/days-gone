package service

import (
	"days-gone/app/dao"
	"days-gone/app/model"
	"days-gone/utils"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"math"
)

var UserFile = &userFileService{}

type userFileService struct {}

// List 获取文件列表
func (f *userFileService) List(r *ghttp.Request, req *model.UserFileListReq) (userFileList []*model.UserFile, total int, err error) {
	condition := make(g.Map)
	condition["user_name"] = utils.Auth.GetTokenData(r).GetString("userKey")
	condition["is_delete"] = 0

	if req.FileName != "" {
		condition["file_name like ?"] = req.FileName + "%"
	}
	if req.Status != "" {
		condition["status"] = req.Status
	} else {
		condition["status != 2"] = nil
	}
	if len(req.UploadAt) > 0 {
		condition["upload_at between ? and ?"] = req.UploadAt
	}

	db := dao.UserFile.Ctx(r.Context())
	limit, offset := req.Paginate()
	total, err = db.Where(condition).Count()
	if err != nil {
		return userFileList, total, gerror.New("query file list failed")
	}
	if total == 0 {
		return []*model.UserFile{}, total, err
	}
	err = db.Limit(limit).Offset(offset).Where(condition).Scan(&userFileList)
	if err != nil {
		return userFileList, total, gerror.New("query user file list failed")
	}
	return userFileList, total, err
}

// Update update user file
func (f *userFileService) Update(r *ghttp.Request, req *model.FileUpdateReq) error {
	db := dao.UserFile.Ctx(r.Context())
	res, err := db.Data(g.Map{"file_name": req.FileName, "status": req.Status}).Where("id", req.Id).Update()
	if err != nil {
		return gerror.New("update failed")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return gerror.New("update failed")
	}
	if rows == 0 {
		return gerror.New("no file updated")
	}
	return err
}

// Delete delete user file
func (f *userFileService) Delete(r *ghttp.Request, id int) error {
	db := dao.UserFile.Ctx(r.Context())
	res, err := db.Data(g.Map{"is_delete": 1}).Where("id", id).Update()
	if err != nil {
		return gerror.New("delete failed")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return gerror.New("delete failed")
	}
	if rows == 0 {
		return gerror.New("no file deleted")
	}
	return err
}

// Download user file
func (f *userFileService) Download(r *ghttp.Request, id int) (s string, err error) {
	db := dao.UserFile.Ctx(r.Context())
	resFst, err := db.Fields("file_sha1").Where("id = ? and status = ? and is_delete = ?", id, 0, 0).Value()
	if err != nil {
		return "", gerror.New("download failed")
	}
	if resFst.IsEmpty() {
		return "", gerror.New("no file to download")
	}
	d := dao.File.Ctx(r.Context())
	resSec, err := d.Fields("file_addr").Where("file_sha1", resFst.String()).Value()
	if err != nil {
		return "", gerror.New("download failed")
	}
	if resSec.IsEmpty() {
		return "", gerror.New("no file to download")
	}
	return resSec.String(), err
}

// UsedCap get the capacity of cloud zone used
func (f *userFileService) UsedCap(r *ghttp.Request) (float64, error) {
	db := dao.UserFile.Ctx(r.GetCtx())
	userName := User.GetCacheUserName(r)
	size, err := db.Fields("SUM(file_size)").Where(g.Map{"user_name": userName, "status": 0, "is_delete": 0}).Value()
	if err != nil {
		return 0, gerror.New("get the capacity failed")
	}
	// here we assume that the total capacity is 0.1M
	return size.Float64()/(math.Pow(1024, 2))/100, err
}

// FilesType get user`s all file`s types
func (f *userFileService) FilesType(r *ghttp.Request) (gdb.Result, error)  {
	db := dao.UserFile.Ctx(r.GetCtx())
	userName := User.GetCacheUserName(r)
	res, err := db.Fields("count(1) as value", "substring_index(file_name, \".\", -1) as name").Where(g.Map{"user_name": userName, "status": 0, "is_delete": 0}).Group("name").OrderAsc("value").All()
	if err != nil {
		return nil, gerror.New("get the types failed")
	}
	return res, err
}

// UploadDays Get the number of files uploaded in the last month
func (f *userFileService) UploadDays(r *ghttp.Request) (gdb.Result, error) {
	userName := User.GetCacheUserName(r)
	res, err := dao.UserFile.UploadFileDays(userName)
	if err != nil {
		return nil, gerror.New("get the files num failed")
	}
	return res, err
}