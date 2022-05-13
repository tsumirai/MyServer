package consts

const (
	ContentImageNumLimit = 9 // 图片上限
)

const (
	ContentTable    = "content" // 表名
	ContentTableNum = 10        // 分表数量
)

const (
	ContentTitleLengthLimit = 20  // 标题字数限制
	ContentLengthLimit      = 300 // 内容字数限制
)

const DefaultContentSize = 10 // 每页默认的内容数

const (
	ContentStatusAuditing     = 0 // 审核中
	ContentStatusAuditFailed  = 1 // 审核未通过
	ContentStatusOnline       = 2 // 上线
	ContentStatusUserDelete   = 3 // 用户删除
	ContentStatusSystemDelete = 4 // 系统下线
)

const (
	ContentListTypeLike    = 1 // 点赞的内容
	ContentListTypeCollect = 2 // 收藏的内容
)

const (
	ContentPermissionNone   = 0 // 无权限限制
	ContentPermissionFriend = 1 // 好友可见
	ContentPermissionSelf   = 2 // 自己可见
)

const (
	ContentCommonSpace = 0 // 普通空间
	ContentSecretSpace = 1 // 隐私空间
)
