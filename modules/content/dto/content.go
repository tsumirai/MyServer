package dto

// CreateContentReq 创建内容请求
type CreateContentReq struct {
	Title            string   `json:"title"`              // 标题
	Content          string   `json:"content"`            // 内容
	ImageUrls        []string `json:"image_urls"`         // 图片url
	AuthorUID        string   `json:"author_uid"`         // 作者uid
	VideoUrl         string   `json:"video_url"`          // 视频url
	ContentType      int64    `json:"content_type"`       // 内容类型 0 图文  1  视频
	ContentSpaceType int64    `json:"content_space_type"` // 内容所属空间类型 0：普通空间  1：隐私空间
	LocCityID        int64    `json:"loc_city_id"`        // 定位城市id
}
