package cache

import "fmt"

func GetUserInfoRedisKey(uid int64) string {
	return fmt.Sprintf("user_info_uid:%v", uid)
}
