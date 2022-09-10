package entity

import "time"

// User 用户玩家
type User struct {
	ID         string    // 玩家ID
	Username   string    // 用户名
	Nickname   string    // 玩家昵称
	CreateTime time.Time // 玩家注册时间
}
