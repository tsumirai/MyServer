package cache

import "fmt"

// GetUserInfoRedisKey 根据UID获得用户数据
func GetUserInfoRedisKey(uid int64) string {
	return fmt.Sprintf("user_info:uid:%v", uid)
}

// GetContentDataByIDsRedisKey 根据作者uid和内容id获得内容数据，hash表结构，subKey为contentID
func GetContentDataByIDsRedisKey(authorUID int64) string {
	return fmt.Sprintf("content_data_by_id:uid:%v", authorUID)
}

// GetContentIDsByAuthorUID 根据作者的UID获得内容ID
func GetContentIDsByAuthorUID(authorUID int64, pageNum, pageSize int) string {
	return fmt.Sprintf("content_list_by_authorUID:authorUID:%v:pageNum:%v:pageSize:%v", authorUID, pageNum, pageSize)
}

// GetCommentDataByIDsRedisKey 根据内容id和评论id获得评论数据，hash表结构，subKey为commentID
func GetCommentDataByIDsRedisKey(contentID int64) string {
	return fmt.Sprintf("comment_data_by_id:id:%v", contentID)
}

// GetCommentIDByContentID 根据内容id获得评论id
func GetCommentIDByContentID(contentID int64, pageNum, pageSize int) string {
	return fmt.Sprintf("comments_by_contentID:contentID:%v:pageNum:%v:pageSize:%v", contentID, pageNum, pageSize)
}

// GetMessageIDsByUID 根据接收者的uid获得消息列表id
func GetMessageIDsByUID(uid int64, pageNum, pageSize int) string {
	return fmt.Sprintf("message_by_uid:uid:%v:pageNum:%v:pageSize:%v", uid, pageNum, pageSize)
}

// GetMessageDataByIDsRedisKey 根据消息id获得消息的数据
func GetMessageDataByIDsRedisKey(receiverUID int64) string {
	return fmt.Sprintf("message_data_by_uid:uid:%v", receiverUID)
}
