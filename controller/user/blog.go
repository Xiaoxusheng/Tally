package user

import (
	"Tally/common"
	"Tally/dao"
	"Tally/global"
	"Tally/models"
	"Tally/utils"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type Blog struct {
	BlogText string `form:"blogText"  validate:"required"`
	Url      string `form:"url" validate:"required"`
}
type Collects struct {
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

// BlogText 发布记账博客
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
	blogs := &models.Blog{
		Identity:     blogId,
		UserIdentity: id,
		ImgUrl:       blog.Url,
		Text:         blog.BlogText,
		Likes:        0,
	}
	//写进数据库
	err = dao.InsertBlog(blogs)
	if err != nil {
		return common.Fail(c, global.BlogCode, global.BlogErr)
	}
	go func() {
		text, err := json.Marshal(blogs)
		if err != nil {
			return
		}
		//写进数据
		global.Global.Redis.ZAdd(global.Global.Ctx, global.BlogText+blogId, redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: text,
		})
		//点赞
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

// IsLike 查询是否点赞
func IsLike(c echo.Context) error {
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
	return common.Ok(c, map[string]bool{
		"is_like": val,
	})
}

// BlogList 博客列表
func BlogList(c echo.Context) error {
	blogId := c.QueryParam("blogId")
	if blogId == "" {
		return common.Fail(c, global.BlogCode, global.BlogErr)
	}
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.BlogCode, global.UserIdentityErr)
	}
	//判断是否取过数据
	lens := 10
	list := make([]*models.Blog, 0, lens)
	blogText := new(models.Blog)
	val := global.Global.Redis.ZScore(global.Global.Ctx, global.BlogText+blogId+"user", id).Val()

	if val != 0 {
		result, err := global.Global.Redis.ZRevRange(global.Global.Ctx, global.BlogText+blogId, int64(val)+1, int64(val+11)).Result()
		if err != nil {
			return err
		}
		//更新索引
		_, err = global.Global.Redis.ZAdd(global.Global.Ctx, global.BlogText+blogId+"user", redis.Z{
			Score:  val + 11,
			Member: id,
		}).Result()
		if err != nil {
			return err
		}
		for i := 0; i < len(result); i++ {
			err := json.Unmarshal([]byte(result[i]), blogText)
			if err != nil {
				return err
			}
			if i == lens-1 {
				break
			}
			list = append(list, blogText)
		}
		return common.Ok(c, list)
	}
	//获取一天范围内的前十条说说
	t := time.Now()
	d := time.Date(t.Year(), t.Month(), t.Day()-1, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	res, err := global.Global.Redis.ZRevRangeByScore(global.Global.Ctx, global.BlogText+blogId, &redis.ZRangeBy{
		Min: strconv.FormatInt(d.Unix(), 10),
		Max: strconv.FormatInt(t.Unix(), 10),
	}).Result()
	if err != nil {
		return err
	}
	for i := 0; i < len(res); i++ {
		err := json.Unmarshal([]byte(res[i]), blogText)
		if err != nil {
			return err
		}
		if i == lens-1 {
			break
		}
		list = append(list, blogText)
	}
	return common.Ok(c, list)
}

// CollectBlog 收藏博客
func CollectBlog(c echo.Context) error {

	return common.Ok(c, nil)
}
