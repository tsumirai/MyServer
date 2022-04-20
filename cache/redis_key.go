package cache

import "fmt"

func GetUserInfoRedisKey(uid string) string {
	return fmt.Sprintf("user_info_uid:%v", uid)
}
