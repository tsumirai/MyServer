package consts

const (
	CommentTable    = "comment"
	CommentTableNum = 10
)

const CommentLengthLimit = 200 // 评论字数限制

const DefaultCommentPageSize = 10 // 默认每页拉取10条评论

const (
	CommentStatusAuditing     = 0 // 审核中
	CommentStatusAuditFailed  = 1 // 审核未通过
	CommentStatusOnline       = 2 // 上线
	CommentStatusUserDelete   = 3 // 用户删除
	CommentStatusSystemDelete = 4 // 系统下线
)
