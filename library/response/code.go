package response

type Code int

const (
	ServerBusy Code = iota
	ErrorParsePram
	SuccessAdd
	SuccessFirst
	SuccessCreated
	SuccessDeleted
	SuccessUpdated
	SuccessFastUpload
	SuccessGetList
	SuccessOperation
	SuccessBatchDeleted
	ErrorAdd
	ErrorGetOne
	ErrorCreated
	ErrorUpdated
	ErrorDeleted
	ErrorGetList
	ErrorOperation
	ErrorBatchDeleted

	ErrorGetCap
	ErrorGetType
	ErrorGetUpDays
	ErrorNoFileUpload
	ErrorFileArdExist
	ErrorDownload

	ErrorUserArdExist
	ErrorSignIn
	ErrorSignInNoFind
	SuccessSignIn
	SuccessSignUp
	SuccessUserInfo
	SuccessGetCap
	SuccessGetType
	SuccessGetUpDays
	ErrorUserInfo
	ErrorCaptcha
)

// Message code 和 message 一一对应
func (c Code) Message() string {
	if msg, ok := Message[c]; ok {
		return msg
	}
	return Message[ServerBusy]
}
