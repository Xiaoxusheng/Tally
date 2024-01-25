package router

import (
	"Tally/controller/admin"
	"Tally/controller/user"
	"Tally/global"
	m "Tally/middleware"
	"Tally/utils"
	"errors"
	"fmt"
	validator "github.com/go-playground/validator/v10"
	casbinmw "github.com/labstack/echo-contrib/casbin"
	"github.com/labstack/echo/v4"
)

func Routers(e *echo.Echo) {
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	//middleware.Recover()
	//登录
	e.POST("/user/login", user.Login)
	//注册
	e.POST("/user/register", user.Register)
	e.GET("/oauth/redirect", user.Token)

	users := e.Group("/user")
	users.GET("/auth2_login", user.OAuthLogin)
	users.Use(m.ParseToken())
	//账单列表
	users.GET("/list", user.TallyList)
	//个人信息
	users.GET("/info", user.Info)
	//修改密码
	users.GET("/change_pwd", user.ChangePwd)
	//修改个人信息
	users.GET("/change_user_info", user.ChangeUserInfo)
	//获取登陆天数
	users.GET("/login_info", user.LoginInfo)
	//注销
	users.GET("/sign_out", user.SignOut)
	//关注
	users.GET("/follow", user.FollowUser)
	//取消关注
	users.GET("/cancel_follow", user.CancelFollow)
	//关注列表
	users.GET("/get_follow_list", user.GetFollowList)
	//共同关注
	users.GET("/together_follow", user.TogetherFollow)
	//退出登录
	users.GET("/logout", user.Logout)
	//添加记账记录
	users.POST("/add", user.AddTallyLog)
	//根据分类获取记账信息
	users.GET("/allot_kind", user.AllotKind)
	//绑定分类
	users.GET("/bind", user.BindKind)
	//记账列表
	users.GET("/date_list", user.DateList)
	//添加分类
	users.Any("/addKind", user.AddKind)
	//分类列表
	users.GET("/kind_list", user.KindList)
	//搜索
	users.GET("/search", user.Search)
	//添加收藏
	users.GET("/add_collect", user.AddCollect)
	//删除收藏
	users.GET("/del_collect", user.DeleteCollect)
	//收藏列表
	users.GET("/collect_list", user.CollectList)
	//分析账单
	users.GET("/analysis", user.Analysis)
	//上传文件
	users.POST("/upload", user.Upload)
	//记账博客
	users.POST("/blog_text", user.BlogText)
	//删除博客
	users.GET("/blog_del", user.DeleteBlog)
	//点赞
	users.GET("/likes", user.Likes)
	//是否点赞
	users.GET("/is_like", user.IsLike)
	//点赞列表
	users.GET("/like_list", user.LikeList)
	//博客列表
	users.GET("/blog_list", user.BlogList)
	//浏览历史
	users.GET("/get_blog_history_list", user.GetBlogHistoryList)
	//修改博客状态
	users.GET("/update_blog_status", user.UpdateBlogStatus)
	//收藏博客
	users.GET("/blog_collect", user.CollectBlog)
	//博客评论
	users.POST("/blog_comment", user.PushComment)
	//评论列表
	users.POST("/blog_comment_list", user.CommentList)
	//删除评论
	users.POST("/blog_del_comment", user.DeleteComment)
	//获取日志文件压缩包
	users.GET("/export_log", user.ExportLog)

	root := e.Group("/admin")
	root.Use(m.ParseToken())
	root.Use(casbinmw.MiddlewareWithConfig(casbinmw.Config{
		Skipper:        nil,
		Enforcer:       global.Global.CasBin,
		EnforceHandler: nil,
		UserGetter: func(c echo.Context) (string, error) {
			id := utils.GetIdentity(c, "identity")
			if id == "" {
				return "", errors.New("identity is not found")
			}
			fmt.Println(id)
			return id, nil
		},
		ErrorHandler: nil,
	}))
	//添加可以访问的资源
	root.POST("/add_resource", admin.AddResource)
	//分配角色
	root.POST("/add_rolesForUser", admin.AddRolesForUser)
	//删除角色
	root.POST("/delete_roleForUser", admin.DeleteRoleForUser)
	//删除资源
	root.POST("/delete_permissionForUser", admin.DeletePermissionForUser)
	//查看能访问的资源
	root.POST("/get_permissionsForUser", admin.GetPermissionsForUser)
	//查看所有权限
	root.POST("/get_allNamedSubjects", admin.GetAllNamedSubjects)
}
