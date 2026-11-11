package model

import (
	"gorm.io/gorm"
)
type ChatHistory struct {
	gorm.Model
	UserID uint `json:"user_id"`//属于哪个用户
	Role string `json:"role"` //角色：user或assistant
	Content string`json:"content"`//聊天内容
}
// 自动迁移建表
func InitChatHistoryTable() {
	DB.AutoMigrate(&ChatHistory{})
}
//保存一条聊天记录
func SaveChat(userID uint,role string,content string){
	DB.Create(&ChatHistory{
		UserID: userID,
		Role: role,
		Content: content,
	})
}
// 获取最近 N 条聊天记录（注意：数据库查出来是倒序的，要反转一下）
func GetRecentChats(userID uint,limit int)[]ChatHistory{
	var chats []ChatHistory
	//按ID倒序拿最近的10条
	DB.Where("user_id=?",userID).Order("id desc").Limit(limit).Find(&chats)
	//反转切片
	for i,j := 0,len(chats)-1;i<j;i,j=i+1,j-1{
		chats[i],chats[j] = chats[j],chats[i]
	}
	return chats
}