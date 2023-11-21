package global

//自定义错误

const (
	UserCode    = 100011
	TallyCode   = 100012
	VerifyCode  = 400000
	KindCode    = 100013
	SearchCode  = 100014
	CollectCode = 100015
)

const (
	LoginErr        = "用户名或密码错误！"
	UserIdentityErr = "获取用户错误！"
	ParseErr        = "解析失败！"
	CollectErr      = "账单不存在！"
	CollectToErr    = "收藏失败！"
)

const (
	Red    = 31
	Yellow = 33
	Blue   = 36
	Gray   = 39
)
