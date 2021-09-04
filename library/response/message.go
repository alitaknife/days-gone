package response

var Message = map[Code]string{
	ServerBusy: "服务器忙ing...",
	// 解析
	ErrorParsePram: "参数解析错误!",

	// 通用
	SuccessAdd:          "添加成功!",
	SuccessFirst:        "获取一条数据成功!",
	SuccessCreated:      "创建成功!",
	SuccessUpdated:      "更新成功!",
	SuccessFastUpload: "秒传成功",
	SuccessDeleted:      "删除成功!",
	SuccessGetList:      "获取列表数据成功!",
	SuccessOperation:    "操作成功!",
	SuccessBatchDeleted: "批量删除成功!",

	ErrorAdd:          "添加失败!",
	ErrorGetOne:       "获取数据信息失败!",
	ErrorCreated:      "创建失败!",
	ErrorUpdated:      "更新失败!",
	ErrorDeleted:      "删除失败!",
	ErrorGetList:      "获取列表数据失败!",
	ErrorOperation:    "操作失败!",
	ErrorBatchDeleted: "批量删除失败!",

	//File
	ErrorNoFileUpload: "没有上传文件!",
	ErrorFileArdExist: "该文件可能已经存在!",
	ErrorDownload:     "下载失败!",

	// User
	ErrorUserArdExist: "该用户名已被占用!",
	ErrorSignIn:       "登录失败!",
	ErrorSignInNoFind: "该用户可能不存在",
	SuccessSignIn:     "登录成功!",
	SuccessSignUp:     "注册成功!",
	SuccessUserInfo:   "获取用户信息成功!",
	ErrorUserInfo:     "获取用户信息失败!",
	ErrorCaptcha:      "验证码获取失败!",
}
