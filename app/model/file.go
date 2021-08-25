package model

// FileListReq 文件列表请求结构体
type FileListReq struct {
	FileName string   `p:"fileName"`
	Status   string      `p:"status" v:"between:0,1#请输入合法的状态值"`
	CreatAt  []string `p:"creatAt"`
	Page
}

// FileUpdateReq 文件更新请求结构体
type FileUpdateReq struct {
	Id       int    `p:"id" v:"required#id不能为空"`
	FileName string `p:"fileName" v:"required#文件名不能为空"`
	Status   string    `p:"status" v:"required|between:0,1#请输入合法的状态值"`
}
