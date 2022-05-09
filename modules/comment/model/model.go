package model

import "time"

type Comment struct {
	ID              int64     `gorm:"column:id" db:"id" json:"id" form:"id"`                                                             //  评论id
	ContentID       int64     `gorm:"column:content_id" db:"content_id" json:"content_id" form:"content_id"`                             //  对应的内容id
	ParentCommentID int64     `gorm:"column:parent_comment_id" db:"parent_comment_id" json:"parent_comment_id" form:"parent_comment_id"` //  父级评论的id
	AuthorUID       int64     `gorm:"column:author_uid" db:"author_uid" json:"author_uid" form:"author_uid"`                             //  评论作者的uid
	Comment         string    `gorm:"column:comment" db:"comment" json:"comment" form:"comment"`                                         //  评论的内容
	CommentStatus   int64     `gorm:"column:comment_status" db:"comment_status" json:"comment_status" form:"comment_status"`             //  评论状态  0：审核中  1：审核未通过  2：上线  3：用户删除  4：系统下线
	AuditReason     int64     `gorm:"column:audit_reason" db:"audit_reason" json:"audit_reason" form:"audit_reason"`                     //  审核未通过原因  1：色情暴力
	CreateTime      time.Time `gorm:"column:create_time" db:"create_time" json:"create_time" form:"create_time"`                         //  评论创建时间
	UpdateTime      time.Time `gorm:"column:update_time" db:"update_time" json:"update_time" form:"update_time"`                         //  评论更新时间
}
