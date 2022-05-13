package dto

import "time"

type ContentRes struct {
	ID               int64     `json:"id"`                 // 内容id
	Title            string    `json:"title"`              // 标题
	Content          string    `json:"content"`            // 内容
	ImageUrls        []string  `json:"image_urls" `        // 图片链接
	AuthorUID        int64     `json:"author_uid"`         // 作者的uid
	VideoUrl         string    `json:"video_url" `         // 视频url
	ContentType      int64     `json:"content_type"`       // 内容类型：0 图文  1  视频
	ContentSpaceType int64     `json:"content_space_type"` // 内容所属空间类型 0：普通空间  1：隐私空间
	LocCityID        int64     `json:"loc_city_id" `       // 定位城市
	ContentStatus    int64     `json:"content_status" `    // 上线状态 0：审核中  1：审核未通过  2：上线  3：用户删除  4：系统下线
	AuditReason      int64     `json:"audit_reason" `      // 审核未通过原因  1：暴力色情
	CreateTime       time.Time `json:"create_time"`        // 创建时间
	UpdateTime       time.Time `json:"update_time"`        // 更新时间
}

// CreateContentReq 创建内容请求
type CreateContentReq struct {
	Title            string   `json:"title"`              // 标题
	Content          string   `json:"content"`            // 内容
	ImageUrls        []string `json:"image_urls"`         // 图片url
	AuthorUID        int64    `json:"author_uid"`         // 作者uid
	VideoUrl         string   `json:"video_url"`          // 视频url
	ContentType      int64    `json:"content_type"`       // 内容类型 0 图文  1  视频
	ContentSpaceType int64    `json:"content_space_type"` // 内容所属空间类型 0：普通空间  1：隐私空间
	LocCityID        int64    `json:"loc_city_id"`        // 定位城市id
}

// ContentReq 获得内容请求
type ContentReq struct {
	ID              int64 `json:"id"`                // 内容id
	AuthorUID       int64 `json:"author_uid"`        // 作者的UID
	PageNum         int   `json:"page_num"`          // 页码
	PageSize        int   `json:"page_size"`         // 每页数量
	ContentListType int   `json:"content_list_type"` // 内容列表的类型 1：喜欢的内容  2：收藏的内容
}

// ContentPermissionReq 修改内容权限请求
type ContentPermissionReq struct {
	ID         int64 `json:"id"`         // 内容id
	AuthorUID  int64 `json:"author_uid"` // 用户的uid
	Permission int   `json:"permission"` // 内容的权限
}

// ContentSpaceReq 修改内容空间请求
type ContentSpaceReq struct {
	ID        int64 `json:"id"`         // 内容id
	AuthorUID int64 `json:"author_uid"` // 用户的uid
	Space     int   `json:"space"`      // 内容的空间
}
