package dto

import "time"

// CreateCommentReq 创建评论请求
type CreateCommentReq struct {
	ContentID       int64  //  对应的内容id
	ParentCommentID int64  //  父级评论的id
	AuthorUID       int64  //  评论作者的uid
	Comment         string //  评论的内容
}

// GetCommentByContentIDReq 根据内容id获得评论的请求
type GetCommentByContentIDReq struct {
	ContentID int64 // 内容的id
	PageNum   int   // 页码
	PageSize  int   // 每页的数量
}

// GetCommentContByContentIDReq 获得内容下的评论数量
type GetCommentContByContentIDReq struct {
	ContentID int64 // 内容的id
}

// CommentRes 评论的返回结构体
type CommentRes struct {
	ID              int64     `json:"id"`                //  评论id
	ContentID       int64     `json:"content_id"`        //  对应的内容id
	ParentCommentID int64     `json:"parent_comment_id"` //  父级评论的id
	AuthorUID       int64     `json:"author_uid"`        //  评论作者的uid
	Comment         string    `json:"comment"`           //  评论的内容
	CommentStatus   int64     `json:"comment_status"`    //  评论状态  0：审核中  1：审核未通过  2：上线  3：用户删除  4：系统下线
	AuditReason     int64     `json:"audit_reason"`      //  审核未通过原因  1：色情暴力
	CreateTime      time.Time `json:"create_time"`       //  评论创建时间
	UpdateTime      time.Time `json:"update_time"`       //  评论更新时间
}
