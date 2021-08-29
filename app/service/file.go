package service

import (
	"context"
	"days-gone/app/dao"
	"days-gone/app/model"
	"days-gone/library/response"
	"days-gone/utils"
	"errors"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
)

var File = &fileService{}

type fileService struct{}

// FastUpload 文件快传
func (f *fileService) FastUpload(r *ghttp.Request, sha1 string, name string) (int8, error) {
	db := dao.File.Ctx(r.Context())
	dbUserFile := dao.UserFile.Ctx(r.Context())
	resFile, err := db.Where("file_sha1", sha1).One()
	if err == nil && len(resFile) > 0{
		resUserFile, err := dbUserFile.Where(g.Map{"file_name": name, "file_sha1": sha1}).One()
		if err != nil {
			return 1, errors.New("文件上传失败")
		}
		if len(resUserFile) > 0{
			// 用户文件表中存在该条数据,直接返回
			return 1, errors.New("该文件已经存在")
		}
		var userFile = &model.UserFile{}
		err = resFile.Struct(userFile)
		if err != nil {
			return 1, errors.New("文件上传失败")
		}
		userFile.UserName = User.GetCacheUserInfo(r).UserName
		userFile.UploadAt = gtime.Now()
		userFile.Status = 0
		// 写入用户文件表
		resInsert, err := dbUserFile.Insert(userFile)
		if err != nil {
			return 1, errors.New("文件上传失败")
		}
		rows, err := resInsert.RowsAffected()
		if err != nil || rows == 0{
			return 1, errors.New("文件上传失败")
		} else {
			// 秒传成功
			return 0, nil
		}
	}
	// 文件表中没有记录要去调用普通接口
	return 2, nil
}

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
	condition["is_delete"] = 0
	if req.Status != "" {
		condition["status"] = req.Status
	}
	if req.FileName != "" {
		condition["file_name like ?"] = req.FileName + "%"
	}
	if len(req.CreatAt) > 0 {
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

// Delete 删除文件
func (f *fileService) Delete(r *ghttp.Request, id int) error {
	db := dao.File.Ctx(r.Context())
	_, err := db.Data(g.Map{"is_delete": 1}).Where("id", id).Update()
	return err
}

// Download 文件下载
func (f *fileService) Download(r *ghttp.Request, id int) (res gdb.Value, err error) {
	db := dao.File.Ctx(r.Context())
	res, err = db.Fields("file_addr").Where("id = ? and status = ? and is_delete = ?", id, 0, 0).Value()
	return res, err
}
