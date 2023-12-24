package global

//自定义错误

const (
	UserCode    = iota + 100011 //用户
	TallyCode                   //账单
	VerifyCode                  //验证失败
	KindCode                    //种类
	SearchCode                  //搜索
	CollectCode                 //收藏
	FileCode                    //文件
	BlogCode                    //博客
	LikesCode                   //点赞
)

// redis 使用前缀
const (
	ListKey         = "kind_list"
	CollectKey      = "collect"
	TallyListKey    = "TallyList"
	BlogLikesKey    = "blogLikes"
	BlogSetLikesKey = "blogLikesSet"
	BlogText        = "blogText"
)

// 错误信息
const (
	LoginErr        = "用户名或密码错误！"
	UserIdentityErr = "获取用户错误！"
	ParseErr        = "解析失败！"
	CollectErr      = "账单不存在！"
	CollectToErr    = "收藏失败！"
	MarshalErr      = "序列化失败！"
	FileErr         = "文件上传错误！"
	BlogErr         = "新增博客失败！"
	LikesErr        = "点赞失败！"
	LikesAlreadyErr = "已经点赞！"

	QueryErr = "获取必要参数失败"
)

// 日=日志颜色
const (
	Red    = 31
	Yellow = 33
	Blue   = 36
	Gray   = 39
)
