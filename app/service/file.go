package service

import (
	"context"
	"days-gone/app/dao"
	"days-gone/app/model"
	"days-gone/utils"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
)

var File = &fileService{}

type fileService struct{}

// FastUpload fast upload by sha1
func (f *fileService) FastUpload(r *ghttp.Request, sha1 string, name string) error {
	db := dao.File.Ctx(r.Context())
	// query file in file db
	resp, err := db.Where("file_sha1", sha1).One()
	if err != nil{
		return gerror.New("fast upload failed")
	} else if resp.IsEmpty(){
		// if there is no this file in file store, should call normal api
		return gerror.NewCode(302, "call the normal api")
	} else {
		// query file in user file db
		dbUserFile := dao.UserFile.Ctx(r.Context())
		reps, err := dbUserFile.Where(g.Map{"file_name": name, "file_sha1": sha1}).One()
		if err != nil {
			return gerror.New("fast upload failed")
		} else if !reps.IsEmpty() {
			return gerror.New("this file maybe exist")
		}
		// if resp is empty, insert this file into user file db
		userFile := &model.UserFile{}
		err = resp.Struct(userFile)
		if err != nil {
			return gerror.Wrap(err, "fast upload failed")
		}
		userFile.UserName = User.GetCacheUserInfo(r).UserName
		userFile.UploadAt = gtime.Now()
		userFile.Status = 0
		// write info into user file db
		resInsert, err := dbUserFile.Insert(userFile)
		if err != nil {
			return gerror.New("fast upload failed")
		}
		rows, err := resInsert.RowsAffected()
		if err != nil || rows == 0{
			return gerror.New("fast upload failed")
		}
		return nil
	}
}

// Upload upload file
func (f *fileService) Upload(r *ghttp.Request, file *ghttp.UploadFile) (err error) {
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
	// open transaction
	err = db.Transaction(r.Context(), func(ctx context.Context, tx *gdb.TX) error {
		_, err := tx.Ctx(ctx).Insert(dao.File.Table, fileInfo) // save file info into file table
		if err != nil {
			return err
		}
		_, err = tx.Ctx(ctx).Insert(dao.UserFile.Table, userFileInfo) // save file info into user file table
		if err != nil {
			return err
		}
		_, err = file.Save(path) // save file into local storage
		return err
	})
	if err != nil {
		return gerror.New("upload file failed")
	}
	return nil
}

// List query file list
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
	if err != nil {
		return fileList, total, gerror.New("query file list failed")
	}
	if total == 0 {
		return []*model.File{}, total, err
	}
	err = db.Limit(limit).Offset(offset).Where(condition).Scan(&fileList)
	if err != nil {
		return fileList, total, gerror.New("query file list failed")
	}
	return fileList, total, err
}

// Update update file info
func (f *fileService) Update(r *ghttp.Request, req *model.FileUpdateReq) error {
	db := dao.File.Ctx(r.Context())
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

// Delete delete file
func (f *fileService) Delete(r *ghttp.Request, id int) error {
	db := dao.File.Ctx(r.Context())
	res, err := db.Data(g.Map{"is_delete": 1}).Where("id", id).Update()
	if err != nil {
		return gerror.New("delete failed")
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return gerror.New("delete failed")
	}
	return err
}

// Download 文件下载
func (f *fileService) Download(r *ghttp.Request, id int) (res gdb.Value, err error) {
	db := dao.File.Ctx(r.Context())
	res, err = db.Fields("file_addr").Where("id = ? and status = ? and is_delete = ?", id, 0, 0).Value()
	if err != nil || res.IsEmpty(){
		return nil, gerror.New("download failed")
	}
	return res, err
}
