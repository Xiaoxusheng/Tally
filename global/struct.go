package global

// UrlList 图片返回结构体
type UrlList struct {
	Url   string `json:"url,omitempty"`
	Index int    `json:"index,omitempty"`
}

const (
	Count = 15 //博客列表数量
)

type User struct {
	Phone    string `json:"phone,omitempty" form:"phone" param:"phone" query:"phone"`
	GithubId string `json:"githubId,omitempty" form:"githubId" param:"githubId" query:"githubId"`
}

//*处理异步数据的管道*/
