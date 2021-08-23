package model

// FileListReq 文件列表请求结构体
type FileListReq struct {
	FileName string   `p:"fileName"`
	Status   int      `p:"status"`
	CreatAt  []string `p:"creatAt"`
	Page
}

// FileUpdateReq 文件更新请求结构体
type FileUpdateReq struct {
	Id       int    `json:"id" v:"required#id不能为空"`
	FileName string `p:"required#文件名不能为空"`
	Status   int    `p:"status" v:"required#状态不能为空"`
}
