package model

import "time"

type Message struct {
	ID            int64     `gorm:"column:id" db:"id" json:"id" form:"id"`                                                 //  消息的id
	SenderUID     int64     `gorm:"column:sender_uid" db:"sender_uid" json:"sender_uid" form:"sender_uid"`                 //  消息发送者的uid
	ReceiverUID   int64     `gorm:"column:receiver_uid" db:"receiver_uid" json:"receiver_uid" form:"receiver_uid"`         //  消息接收者的uid
	MessageType   int64     `gorm:"column:message_type" db:"message_type" json:"message_type" form:"message_type"`         //  消息类型 1：被点赞内容  2：被点赞评论  3：被收藏  4：被下线  5：被关注  6：被封禁  7：解除封禁
	RelateID      int64     `gorm:"column:relate_id" db:"relate_id" json:"relate_id" form:"relate_id"`                     //  相关id，如评论id，内容id，用户uid
	MessageStatus int64     `gorm:"column:message_status" db:"message_status" json:"message_status" form:"message_status"` //  消息状态  0：未读，1：已读，2：删除
	ExtraInfo     string    `gorm:"column:extra_info" db:"extra_info" json:"extra_info" form:"extra_info"`                 //  额外信息
	CreateTime    time.Time `gorm:"column:create_time" db:"create_time" json:"create_time" form:"create_time"`             //  消息创建时间
	UpdateTime    time.Time `gorm:"column:update_time" db:"update_time" json:"update_time" form:"update_time"`             //  消息更新时间
}
