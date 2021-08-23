package service

import (
	"context"
	"days-gone/app/dao"
	"days-gone/app/model"
	"days-gone/library/response"
	"days-gone/utils"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
)

var File = fileService{}

type fileService struct{}

// Upload 文件上传
func (f *fileService) Upload(r *ghttp.Request, file *ghttp.UploadFile) (code response.Code, err error) {
	path := g.Cfg().GetString("upload.Path")
	sha1, _ := utils.Sha1Encrypt(file)
	fileInfo := &model.File{
		FileName: file.Filename,
		FileAddr: gfile.Join(path, file.Filename),
		FileSize: file.Size,
		FileSha1: sha1,
		CreateAt: gtime.Now(),
	}

	userFileInfo := &model.UserFile{
		UserName: User.GetCacheUserInfo(r).UserName,
		FileSha1: sha1,
		FileSize: file.Size,
		FileName: file.Filename,
		UploadAt: gtime.Now(),
	}

	db := dao.File.Ctx(r.Context())
	// 开启事务
	err = db.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
		_, err := tx.Ctx(ctx).Insert(dao.File.Table, fileInfo) // 保存文件信息到表
		if err != nil {
			return err
		}
		_, err = tx.Ctx(ctx).Insert(dao.UserFile.Table, userFileInfo) // 保存文件信息到用户文件表
		if err != nil {
			return err
		}
		_, err = file.Save(path) // 保存文件到本地
		return err
	})
	if err != nil {
		return response.ErrorAdd, err
	}
	return response.SuccessAdd, err
}

// List 获取文件列表
func (f *fileService) List(r *ghttp.Request, req *model.FileListReq) (fileList []*model.File, total int, err error) {
	condition := make(g.Map)
	condition["status"] = req.Status
	if req.FileName != "" {
		condition["file_name like ?"] = req.FileName + "%"
	} else if len(req.CreatAt) > 0 {
		condition["create_at between ? and ?"] = req.CreatAt
	}

	db := dao.File.Ctx(r.Context())
	fileList = ([]*model.File)(nil)
	limit, offset := req.Paginate()
	total, err = db.Where(condition).Count()
	err = db.Limit(limit).Offset(offset).Where(condition).Scan(&fileList)
	return fileList, total, err
}

// Update 更新文件
func (f *fileService) Update(r *ghttp.Request, req *model.FileUpdateReq) error {
	db := dao.File.Ctx(r.Context())
	_, err := db.Data(g.Map{"file_name": req.FileName, "status": req.Status}).Where("id", req.Id).Update()
	return err
}

// Download 文件下载
func (f *fileService) Download(r *ghttp.Request, id int) (file *model.File, err error) {
	db := dao.File.Ctx(r.Context())
	file = (*model.File)(nil)
	err = db.Where("id", id).Scan(&file)
	return file, err
}
