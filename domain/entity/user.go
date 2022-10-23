package entity

import "time"

// UserInfo 用户玩家
type UserInfo struct {
	ID         string    `json:"id"`         // 玩家ID
	Username   string    `json:"username"`   // 用户名
	Nickname   string    `json:"nickname"`   // 玩家昵称
	Avatar     string    `json:"avatar"`     // 头像
	CreateTime time.Time `json:"createTime"` // 玩家注册时间
}
