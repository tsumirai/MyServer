package dto

import "time"

// GetMessageListByUID 根据uid获得消息列表
type GetMessageListByUID struct {
	UID      int64 `json:"uid"`       // 消息接收人的uid
	PageNum  int   `json:"page_num"`  // 消息页数
	PageSize int   `json:"page_size"` // 消息每页的数量
}

// Message 消息格式
type Message struct {
	ID            int64     `json:"id"`             //  消息的id
	SenderUID     int64     `json:"sender_uid"`     //  消息发送者的uid
	ReceiverUID   int64     `json:"receiver_uid"`   //  消息接收者的uid
	MessageType   int64     `json:"message_type"`   //  消息类型 1：被点赞内容  2：被点赞评论  3：被收藏  4：被下线  5：被关注  6：被封禁  7：解除封禁
	RelateID      int64     `json:"relate_id"`      //  相关id，如评论id，内容id，用户uid
	MessageStatus int64     `json:"message_status"` //  消息状态  0：未读，1：已读，2：删除
	ExtraInfo     string    `json:"extra_info"`     //  额外信息
	CreateTime    time.Time `json:"create_time"`    //  消息创建时间
	UpdateTime    time.Time `json:"update_time"`    //  消息更新时间
}
