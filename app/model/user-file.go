package model

// UserFileListReq 用户文件列表请求结构体
type UserFileListReq struct {
	FileName string   `p:"fileName"`
	Status   string      `p:"status" v:"between:0,1#请输入合法的状态值"`
	UploadAt  []string `p:"uploadAt"`
	Page
}

// UserFileUpdateReq 用户文件更新请求结构体
type UserFileUpdateReq struct {
	Id       int    `p:"id" v:"required#id不能为空"`
	FileName string `p:"fileName" v:"required#文件名不能为空"`
	Status   string    `p:"status" v:"between:0,1#请输入合法的状态值"`
}