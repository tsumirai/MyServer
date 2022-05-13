package router

import (
	commentController "MyServer/modules/comment/controller"
	contentController "MyServer/modules/content/controller"
	interactController "MyServer/modules/interact/controller"
	messageController "MyServer/modules/message/controller"
	testController "MyServer/modules/test_tool/controller"
	userController "MyServer/modules/user/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter(R *gin.Engine) {
	R.GET("/ping", new(testController.TestToolController).Ping)

	/*------------------- user ---------------------*/
	R.POST("/user/userLogin", userController.NewUserController().UserLogin)               // 用户登录，新用户自动注册
	R.POST("/user/getUserInfoByUID", userController.NewUserController().GetUserInfoByUID) // 根据UID获取用户信息
	R.POST("/user/updateUserInfo", userController.NewUserController().UpdateUserInfo)     // 更新用户信息

	/*-------------------- content -----------------*/
	R.POST("/content/feed", contentController.NewContentController().Feed)                                 // 获得内容的feed流
	R.POST("/content/createContent", contentController.NewContentController().CreateContent)               // 创建内容
	R.POST("/content/deleteContentByID", contentController.NewContentController().DeleteContentByID)       // 用户根据ID删除内容
	R.POST("/content/getContentByID", contentController.NewContentController().GetContentByID)             // 根据内容ID获得内容
	R.POST("/content/getContentList", contentController.NewContentController().GetContentList)             // 获得内容列表（获得收藏的内容，获得点赞的内容）
	R.POST("/content/shareContent", contentController.NewContentController().ShareContent)                 // 分享内容
	R.POST("/content/setContentPermission", contentController.NewContentController().SetContentPermission) // 设置内容的权限 0：无权限  1：好友可见  2：自己可见
	R.POST("/content/setContentSpace", contentController.NewContentController().SetContentSpace)           // 设置内容的空间 0：普通空间  1：隐私空间

	/*--------------------- comment ----------------*/
	R.POST("/comment/createComment", commentController.NewCommentController().CreateComment)                           // 创建评论
	R.POST("/comment/getCommentsByContentID", commentController.NewCommentController().GetCommentsByContentID)         // 获得内容下的评论
	R.POST("/comment/getCommentCountByContentID", commentController.NewCommentController().GetCommentCountByContentID) // 获得内容下的评论数量

	/*-------------------- message -----------------*/
	R.POST("/message/getMessageListByUID", messageController.NewMessageController().GetMessageListByUID) // 获得消息列表
	R.POST("/message/updateMessageStatus", messageController.NewMessageController().UpdateMessageStatus) // 修改消息状态

	/*-------------------- interact ---------------*/
	R.POST("/interact/like", interactController.NewInteractController().Like)             // 点赞、取消点赞，收藏、取消收藏
	R.POST("/interact/report", interactController.NewInteractController().Report)         // 举报用户、举报内容、举报评论
	R.POST("/interact/followUser", interactController.NewInteractController().FollowUser) // 关注、取消关注用户
}
