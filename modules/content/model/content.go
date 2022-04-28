package model

import "time"

type Content struct {
	ID               int64     `gorm:"column:id" db:"id" json:"id" form:"id"`
	Title            string    `gorm:"column:title" db:"title" json:"title" form:"title"`                                                     //  标题
	Content          string    `gorm:"column:content" db:"content" json:"content" form:"content"`                                             //  内容
	ImageUrls        string    `gorm:"column:image_urls" db:"image_urls" json:"image_urls" form:"image_urls"`                                 //  图片链接
	AuthorUID        string    `gorm:"column:author_uid" db:"author_uid" json:"author_uid" form:"author_uid"`                                 //  作者的uid
	VideoUrl         string    `gorm:"column:video_url" db:"video_url" json:"video_url" form:"video_url"`                                     //  视频url
	ContentType      int64     `gorm:"column:content_type" db:"content_type" json:"content_type" form:"content_type"`                         //  内容类型：0 图文  1  视频
	ContentSpaceType int64     `gorm:"column:content_space_type" db:"content_space_type" json:"content_space_type" form:"content_space_type"` //  内容所属空间类型 0：普通空间  1：隐私空间
	LocCityID        int64     `gorm:"column:loc_city_id" db:"loc_city_id" json:"loc_city_id" form:"loc_city_id"`                             // 定位城市
	ContentStatus    int64     `gorm:"column:content_status" db:"content_status" json:"content_status" form:"content_status"`                 // 上线状态 0：审核中  1：审核未通过  2：上线  3：用户删除  4：系统下线
	AuditReason      int64     `gorm:"column:audit_reason" db:"audit_reason" json:"audit_reason" form:"audit_reason"`                         // 审核未通过原因  1：暴力色情
	CreateTime       time.Time `gorm:"column:create_time" db:"create_time" json:"create_time" form:"create_time"`                             //  创建时间
	UpdateTime       time.Time `gorm:"column:update_time" db:"update_time" json:"update_time" form:"update_time"`                             //  更新时间
}
