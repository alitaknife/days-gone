package service

import (
	"days-gone/app/dao"
	"days-gone/app/model"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var UserFile = &userFileService{}

type userFileService struct {}

// List 获取文件列表
func (f *userFileService) List(r *ghttp.Request, req *model.UserFileListReq) (userFileList []*model.UserFile, total int, err error) {
	condition := make(g.Map)
	condition["user_name"] = User.GetCacheUserInfo(r).UserName
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
	userFileList = ([]*model.UserFile)(nil)
	limit, offset := req.Paginate()
	total, err = db.Where(condition).Count()
	err = db.Limit(limit).Offset(offset).Where(condition).Scan(&userFileList)
	return userFileList, total, err
}

// Update 更新用户文件
func (f *userFileService) Update(r *ghttp.Request, req *model.FileUpdateReq) error {
	db := dao.UserFile.Ctx(r.Context())
	_, err := db.Data(g.Map{"file_name": req.FileName, "status": req.Status}).Where("id", req.Id).Update()
	return err
}

// Delete 删除用户文件
func (f *userFileService) Delete(r *ghttp.Request, id int) error {
	db := dao.UserFile.Ctx(r.Context())
	_, err := db.Data(g.Map{"is_delete": 1}).Where("id", id).Update()
	return err
}

// Download 用户文件下载
func (f *userFileService) Download(r *ghttp.Request, id int) (s string) {
	db := dao.UserFile.Ctx(r.Context())
	resFst, err := db.Fields("file_sha1").Where("id = ? and status = ? and is_delete = ?", id, 0, 0).Value()
	if err != nil || resFst.IsEmpty(){
		return ""
	}
	d := dao.File.Ctx(r.Context())
	resSec, err := d.Fields("file_addr").Where("file_sha1", resFst.String()).Value()
	if err != nil || resSec.IsEmpty() {
		return ""
	}
	return resSec.String()
}