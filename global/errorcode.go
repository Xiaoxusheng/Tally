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
	Collect                     //收藏博客
	CommentCode                 //评论
	LogCode                     //日志压缩
)

// redis 使用前缀
const (
	SignIn          = "sign-in"
	UserFollow      = "follow"
	BanUser         = "ban_user"
	ListKey         = "kind_list"
	CollectKey      = "collect"
	TallyListKey    = "TallyList"
	BlogLikesKey    = "blogLikes"
	BlogSetLikesKey = "blogLikesSet"
	BlogText        = "blogText"
	BlogCollects    = "blogCollect"
	BlogCollectRem  = "blogCollectRem"
	BlogHistory     = "blogHistory"
)

// 错误信息
const (
	UserNotFound         = "用户不存在！"
	BannedUser           = "用户被封禁！"
	ChangeUserInfo       = "修改用户信息失败！"
	LoginErr             = "用户名或密码错误！"
	PassISNull           = "密码不能为空！"
	PasswordIeErr        = "密码错误！"
	ChangePassword       = "修改密码失败！"
	UserIdentityErr      = "获取用户错误！"
	AlreadyFollow        = "已经关注！"
	FollowFail           = "关注失败！"
	FollowNot            = "不能关注自己！"
	AlreadyCancelFollow  = "已经取消关注！"
	CancelFollowFail     = "取消关注失败！"
	ParseErr             = "解析失败！"
	CollectErr           = "账单不存在！"
	CollectToErr         = "收藏失败！"
	MarshalErr           = "序列化失败！"
	FileErr              = "文件上传错误！"
	BlogErr              = "新增博客失败！"
	LikesErr             = "点赞失败！"
	LikesAlreadyErr      = "已经点赞！"
	GetLikeListErr       = "获取点赞列表失败！"
	QueryErr             = "获取必要参数失败！"
	BlogCollect          = "收藏博客失败！"
	BlogNotFound         = "博客不存在！"
	DeleteBlogFail       = "删除博客失败！"
	UpdateBlogStatusFail = "修改博客状态失败！"
	CreateLogErr         = "生成zip失败！"
)

// 日=日志颜色
const (
	Red    = 31
	Yellow = 33
	Blue   = 36
	Gray   = 39
)

// 常量
const (
	Success = 1
	Fail    = 0
)

//s为单位

const (
	InfoTime = 7 * 24 * 60 * 60
)
