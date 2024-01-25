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

type Comment struct {
	BlogID   string `json:"blog_id," query:"blog_id" form:"blog_id"  param:"blog_id"`
	Text     string `json:"text" query:"text" form:"text" param:"text"`
	ParentID string `json:"parent_id,omitempty" query:"parent_id" form:"parent_id" param:"parent_id"`
}

type DelComment struct {
	BlogID    string `json:"blog_id," query:"blog_id" form:"blog_id"  param:"blog_id"`
	CommentId string `json:"comment_id" query:"comment_id" form:"comment_id"  param:"comment_id"`
}
