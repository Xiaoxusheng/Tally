package user

import (
	"Tally/common"
	"Tally/dao"
	"Tally/global"
	"Tally/models"
	"Tally/utils"
	"fmt"
	"github.com/labstack/echo/v4"
)

type Blog struct {
	BlogText string `form:"blogText"  validate:"required"`
	Url      string `form:"url" validate:"required"`
}

// Upload 上传文件
func Upload(c echo.Context) error {
	list := make([]global.UrlList, 0, 9)
	n := 0
	form, err := c.MultipartForm()
	if err != nil {
		global.Global.Log.Warn("上传失败" + err.Error())
		return common.Fail(c, global.FileCode, global.FileErr)
	}
	if len(form.File["file"]) > 9 {
		return common.Fail(c, global.FileCode, global.FileErr)
	}
	urlChan := make(chan global.UrlList, 9)
	for i, r := range form.File["file"] {
		go utils.Upload(r, urlChan, i+1)
	}
	for n != len(form.File["file"]) {
		select {
		case s := <-urlChan:
			if s.Url != "" {
				n++
				list = append(list, s)
				continue
			}
			n++
		}
	}
	return common.Ok(c, list)
}

// BlogText 记账博客
func BlogText(c echo.Context) error {
	blog := new(Blog)
	err := c.Bind(blog)
	if err != nil {
		return common.Fail(c, global.BlogCode, global.ParseErr)
	}
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.BlogCode, global.UserIdentityErr)
	}
	blogId := utils.GetUidV4()
	//写进数据库
	err = dao.InsertBlog(&models.Blog{
		Identity:     blogId,
		UserIdentity: id,
		ImgUrl:       blog.Url,
		Text:         blog.BlogText,
		Likes:        0,
	})
	if err != nil {
		return common.Fail(c, global.BlogCode, global.BlogErr)
	}
	go func() {
		//写进
		global.Global.Redis.Set(global.Global.Ctx, global.BlogLikesKey+blogId, 0, 0)
	}()
	return common.Ok(c, nil)
}

// Likes 博客点赞
func Likes(c echo.Context) error {
	blogId := c.QueryParam("blogId")
	if blogId == "" {
		return common.Fail(c, global.LikesCode, global.QueryErr)
	}
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.LikesCode, global.UserIdentityErr)
	}
	//判断是否已经点赞
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.BlogSetLikesKey+blogId, id).Val()
	fmt.Println("vals", val)
	if val {
		//删除
		global.Global.Redis.SRem(global.Global.Ctx, global.BlogSetLikesKey+blogId, id)
		res, err := global.Global.Redis.IncrBy(global.Global.Ctx, global.BlogLikesKey+blogId, -1).Result()
		global.Global.Log.Info(res, err)
		if err != nil {
			return common.Fail(c, global.LikesCode, global.LikesErr)
		}
		return common.Fail(c, global.LikesCode, global.LikesAlreadyErr)
	}
	//添加进集合
	result, err := global.Global.Redis.SAdd(global.Global.Ctx, global.BlogSetLikesKey+blogId, id).Result()
	global.Global.Log.Info("reslu", result)
	if err != nil {
		return err
	}
	//加1
	res, err := global.Global.Redis.IncrBy(global.Global.Ctx, global.BlogLikesKey+blogId, 1).Result()
	global.Global.Log.Info(res, err)
	if err != nil {
		return common.Fail(c, global.LikesCode, global.LikesErr)
	}
	return common.Ok(c, nil)
}

// BlogList 博客列表
func BlogList(c echo.Context) error {
	return common.Ok(c, nil)
}
