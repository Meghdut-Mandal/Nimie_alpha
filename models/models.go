package models

import (
	"Nimie/config"
	"Nimie/utils"
	"gorm.io/gorm"
	"time"
)

type Conversation struct {
	ConversationId int64 `json:"conversation_id" gorm:"primaryKey;autoIncrement;notNull"`
	UserIdA        int64 `json:"user_id_a"`
	UserIdB        int64 `json:"user_id_b"`
	CreatedAt      int64 `json:"created_at" gorm:"autoCreateTime"`
	StatusId       int64 `json:"status_id"`
}

type User struct {
	UserId     int64  `json:"user_id" gorm:"primaryKey;autoIncrement;notNull"`
	CreateTime int64  `json:"create_time" gorm:"autoCreateTime"`
	PublicKey  string `json:"public_key"`
}

type ChatMessage struct {
	MessageId      int64  `json:"message_id" gorm:"primaryKey;autoIncrement;notNull"`
	ConversationId int64  `json:"conversation_id"`
	CreateTime     int64  `json:"create_time" gorm:"autoCreateTime"`
	UserId         int64  `json:"user_id"`
	Message        string `json:"message"`
}

type Status struct {
	StatusId   int64  `json:"status_id" gorm:"primaryKey;autoIncrement;notNull"`
	UserId     int64  `json:"user_id"`
	CreateTime int64  `json:"create_time" gorm:"autoCreateTime"`
	HeaderText string `json:"header_text"`
	LinkId     string `json:"link_id"`
}

var db *gorm.DB

func init() {
	config.Connect()
	db = config.GetSqlDB()
	db.AutoMigrate(&Conversation{}, &User{}, &ChatMessage{}, &Status{})
}

// AddConversation Add conversion to database
func AddConversation(conversation *Conversation) {
	db.Create(conversation)
}

// GetConversation Get conversation from database
func GetConversation(conversationId int64) *Conversation {
	conversation := &Conversation{}
	db.Where("conversation_id = ?", conversationId).First(conversation)
	return conversation
}

// AddNewUser Add user to db
func AddNewUser(publicKey string) *User {
	user := &User{
		PublicKey: publicKey,
	}
	db.Create(user)
	return user
}

// AddStatus Create a new status
func AddStatus(text *string, userId int64) *Status {
	status := &Status{
		UserId:     userId,
		HeaderText: *text,
		LinkId:     utils.Base62Encode(time.Now().UnixNano() - utils.UnixTime2021),
	}
	db.Create(status)
	return status
}

// AddMessage Add messages to db
func AddMessage(message *ChatMessage) {
	db.Create(message)
}

// RemoveStatus remove status from the database
func RemoveStatus(statusId int64, userId int64) string {
	status := &Status{}
	db.Where("status_id = ?", statusId).First(status)

	if status.CreateTime == 0 {
		return "Status not found"
	} else if status.UserId != userId {
		return "You are not the owner of this status"
	}

	if status.UserId == userId {
		db.Delete(status)
		return "Status deleted"
	}
	return "You are not allowed to delete this status"
}
