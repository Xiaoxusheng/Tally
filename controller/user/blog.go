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

/*
博客模块
*/

type Blog struct {
	BlogText string `form:"blogText"  validate:"required"`
	Url      string `form:"url" validate:"required"`
}

type UpdateStatus struct {
	BlogId string `json:"blogId,omitempty" form:"blogId" query:"blogId" param:"blogId"`
	Status int    `json:"status,omitempty" form:"status" query:"status" param:"status"`
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
		IsHide:       false,
		Likes:        0,
		IP:           c.RealIP(),
		ViolateRule:  false,
	}
	//写进数据库
	err = dao.InsertBlog(blogs)
	if err != nil {
		global.Global.Log.Warn(err)
		return common.Fail(c, global.BlogCode, global.BlogErr)
	}
	global.Global.Pool.Submit(func() {
		text, err := json.Marshal(blogs)
		if err != nil {
			return
		}
		//写进sortSet
		_, err = global.Global.Redis.ZAdd(global.Global.Ctx, global.BlogText, redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: text,
		}).Result()
		if err != nil {
			global.Global.Log.Warn(err)
		}
		result, err := global.Global.Redis.ZRange(global.Global.Ctx, global.BlogText, 0, time.Now().Unix()).Result()
		if err != nil {
			return
		}
		fmt.Println(result)
		//把博客id存入set
		_, err = global.Global.Redis.SAdd(global.Global.Ctx, global.BlogText+":IdList", blogId).Result()
		if err != nil {
			global.Global.Log.Warn(err)
		}
		//点赞
		_, err = global.Global.Redis.Set(global.Global.Ctx, global.BlogLikesKey+blogId, 0, 0).Result()
		if err != nil {
			global.Global.Log.Warn(err)
		}
	})
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
	//判断博客是否存在
	value := global.Global.Redis.SIsMember(global.Global.Ctx, global.BlogText+":IdList", blogId).Val()
	if !value {
		return common.Fail(c, global.LikesCode, global.BlogNotFound)
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
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.BlogCode, global.UserIdentityErr)
	}
	//判断是否取过数据
	list := make([]models.Blog, 0, global.Count)
	blogText := new(models.Blog)
	val := global.Global.Redis.ZScore(global.Global.Ctx, global.BlogText+"user", id).Val()
	var t1 float64 = 0
	var count int64 = 1
	fmt.Println("val", val)
	//val
	if val != global.Fail {
		t := time.Now()
		d := time.Date(t.Year(), t.Month(), t.Day()-1, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
		//获取
		l := len(global.Global.Redis.ZRangeByScore(global.Global.Ctx, global.BlogText, &redis.ZRangeBy{
			Min:    strconv.FormatInt(d.Unix(), 10),
			Max:    strconv.FormatInt(t.Unix(), 10),
			Offset: 0,
			Count:  global.Count,
		}).Val())
		if l <= global.Collect {
			fmt.Println("l", l)
			val = 0
		}
		result, err := global.Global.Redis.ZRevRangeByScoreWithScores(global.Global.Ctx, global.BlogText, &redis.ZRangeBy{
			Min:    strconv.FormatInt(d.Unix(), 10),
			Max:    strconv.FormatInt(t.Unix(), 10),
			Offset: int64(0),
			Count:  global.Count,
		}).Result()

		fmt.Println("res", result)
		if err != nil {
			return err
		}
		global.Global.Log.Warn("len", len(result))
		for i := 0; i < len(result); i++ {
			err = json.Unmarshal([]byte(result[i].Member.(string)), blogText)
			if err != nil {
				return err
			}
			if t1 == result[i].Score {
				count++
			} else {
				t1 = result[i].Score
				count = 1
			}
			list = append(list, *blogText)
		}
		fmt.Println("list", list)
		//存入有几条相同的数据
		_, err = global.Global.Redis.ZAdd(global.Global.Ctx, global.BlogText+"user", redis.Z{
			Score:  float64(count),
			Member: id,
		}).Result()
		if err != nil {
			return err
		}
		return common.Ok(c, list)
	}
	//获取一天范围内的前十条说说
	t := time.Now()
	d := time.Date(t.Year(), t.Month(), t.Day()-1, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	res, err := global.Global.Redis.ZRevRangeByScoreWithScores(global.Global.Ctx, global.BlogText, &redis.ZRangeBy{
		Min:    strconv.FormatInt(d.Unix(), 10),
		Max:    strconv.FormatInt(t.Unix(), 10),
		Offset: 0,
		Count:  global.Count,
	}).Result()
	fmt.Println(res)
	//第一次取

	if err != nil {
		return err
	}

	for i := 0; i < len(res); i++ {
		err := json.Unmarshal([]byte(res[i].Member.(string)), blogText)
		if err != nil {
			return err
		}
		if t1 == res[i].Score {
			count++
		} else {
			t1 = res[i].Score
			count = 1
		}
		list = append(list, *blogText)
	}
	if len(list) == 0 {
		count = 0
	}
	global.Global.Redis.ZAdd(global.Global.Ctx, global.BlogText+"user", redis.Z{
		Score:  float64(count),
		Member: id,
	})
	return common.Ok(c, list)
}

// CollectBlog 收藏博客
func CollectBlog(c echo.Context) error {
	blogId := c.QueryParam("blogId")
	if blogId == "" {
		return common.Fail(c, global.LikesCode, global.QueryErr)
	}
	id := utils.GetIdentity(c, "identity")
	if id == "" {
		return common.Fail(c, global.LikesCode, global.UserIdentityErr)
	}
	//判断博客id是否存在
	val := global.Global.Redis.SIsMember(global.Global.Ctx, global.BlogText+":IdList", blogId).Val()
	if !val {
		return common.Fail(c, global.LikesCode, global.BlogNotFound)
	}
	//判断是否收藏
	value := global.Global.Redis.SIsMember(global.Global.Ctx, global.BlogCollects+id, blogId).Val()
	if value {
		//在玩家的收藏列表中删除收藏
		res, err := global.Global.Redis.SRem(global.Global.Ctx, global.BlogCollects+id, blogId).Result()
		global.Global.Log.Info(res, err)
		//添加进删除的
		res, err = global.Global.Redis.SAdd(global.Global.Ctx, global.BlogCollectRem+id, blogId).Result()
		global.Global.Log.Info(res, err)
		if err != nil {
			global.Global.Log.Warn(err)
			return common.Fail(c, global.Collect, global.BlogCollect)
		}
		return common.Fail(c, global.Collect, global.BlogCollect)
	}
	//收藏博客的set,未收藏
	_, err := global.Global.Redis.SAdd(global.Global.Ctx, global.BlogCollects+id, blogId).Result()
	if err != nil {
		global.Global.Log.Warn(err)
		return common.Fail(c, global.Collect, global.BlogCollect)
	}
	res, err := global.Global.Redis.SRem(global.Global.Ctx, global.BlogCollectRem+id, blogId).Result()
	global.Global.Log.Info(res, err)
	if err != nil {
		global.Global.Log.Warn(err)
		return common.Fail(c, global.Collect, global.BlogCollect)
	}
	if uid := dao.GetIdByBlog(blogId); uid != "" {
		//	存进数据库
		err = dao.InsertBlogCollect(&models.Collect{
			Identity:     utils.GetUidV4(),
			UserIdentity: id,
			CollectId:    uid,
			BlogId:       blogId,
		})
		if err != nil {
			global.Global.Log.Warn(err)
			return common.Ok(c, "")
		}
	}
	return common.Ok(c, nil)
}

// GetBlogHistoryList 浏览历史
func GetBlogHistoryList(c echo.Context) error {

	return common.Ok(c, nil)
}

// UpdateBlogStatus 修改博客的状态
func UpdateBlogStatus(c echo.Context) error {
	updateStatus := new(UpdateStatus)
	err := c.Bind(updateStatus)
	if err != nil {
		return common.Fail(c, global.BlogCode, global.QueryErr)
	}
	if updateStatus.Status > 2 || updateStatus.Status < 0 {
		return common.Fail(c, global.BlogCode, global.QueryErr)
	}
	//判断是否是这个人的博客

	//判断博客是否存在
	value := global.Global.Redis.SIsMember(global.Global.Ctx, global.BlogText+":IdList", updateStatus.BlogId).Val()
	if !value {
		return common.Fail(c, global.BlogCode, global.BlogNotFound)
	}
	//修改数据库
	err = dao.UpdateStatus(updateStatus.BlogId, updateStatus.Status)
	if err != nil {
		global.Global.Log.Warn(err)
		return common.Fail(c, global.BlogCode, "修改失败")
	}
	return common.Ok(c, nil)
}

//博客详情

// DeleteBlog 删除博客
func DeleteBlog(c echo.Context) error {
	//博客id
	blogId := c.QueryParam("blog_id")
	if blogId == "" {
		return common.Fail(c, global.BlogCode, global.QueryErr)
	}
	//判断博客是否存在
	value := global.Global.Redis.SIsMember(global.Global.Ctx, global.BlogText+":IdList", blogId).Val()
	if !value {
		return common.Fail(c, global.BlogCode, global.BlogNotFound)
	}
	//删除
	if global.Global.Redis.SRem(global.Global.Ctx, global.BlogText+":IdList", blogId).Val() == global.Fail {
		return common.Fail(c, global.BlogCode, global.DeleteBlogFail)
	}
	return common.Ok(c, nil)
}

//
