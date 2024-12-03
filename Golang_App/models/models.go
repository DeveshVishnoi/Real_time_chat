package models

type User struct {
	UserName string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
	Online   string `json:"online" bson:"online"`
}

// ConversationStruct is a universal struct for mapping the conversations
type ConversationStruct struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	Message    string `json:"message"`
	ToUserID   string `json:"toUserID"`
	FromUserID string `json:"fromUserID"`
}

type MessagePayloadStruct struct {
	FromUserID string `json:"fromUserID"`
	ToUserID   string `json:"toUserID"`
	Message    string `json:"message"`
}
