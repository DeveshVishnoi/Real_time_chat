package database

import (
	"realtime_chat/models"

	"go.mongodb.org/mongo-driver/mongo"
)

// For collection
type DBHelper struct {
	//Mongo Client provider
	MongoClient *mongo.Client

	//Collections
	UserCollection         *mongo.Collection
	ConversationCollection *mongo.Collection
}

// For fucntions
type DBHelperProvider interface {
	CreateUser(models.User) error
	UpdateUserStatus(string, string) error
	GetUserbyEmailId(string) (models.User, error)
	CheckUserExists(string) (bool, error)
	GetUserStatus(string) (models.User, error)
	GetOnlineUsers(string) ([]models.User, error)
	GetConversationBetweenTwoUsers(string, string) []models.ConversationStruct
	StoreNewChatMessages(models.MessagePayloadStruct) bool
}

// Initialize database.
func NewDbProvider(client *mongo.Client) DBHelperProvider {

	return &DBHelper{
		MongoClient:            client,
		UserCollection:         client.Database(models.DataBase_Name).Collection(models.UserCollection_Name),
		ConversationCollection: client.Database(models.DataBase_Name).Collection(models.Conversation_Name),
	}
}
