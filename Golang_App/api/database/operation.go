package database

import (
	"context"
	"errors"
	"fmt"
	"realtime_chat/models"
	"realtime_chat/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func (dbhelper *DBHelper) CreateUser(user models.User) error {

	// Check if the email already exists
	var existingUser models.User
	err := dbhelper.UserCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		fmt.Println("Email already exists in the database")
		return fmt.Errorf("email %s already exists", user.Email)
	}

	if err != mongo.ErrNoDocuments {
		fmt.Println("error checking if email exists - ", err)
		return err
	}

	// Hash the password
	hashPassword, err := utils.CreatePassword(user.Password)
	if err != nil {
		fmt.Println("error generating the hashed password - ", err)
		return err
	}

	// Insert the new user
	insertedId, err := dbhelper.UserCollection.InsertOne(context.TODO(), bson.M{
		"username": user.UserName,
		"email":    user.Email,
		"password": hashPassword,
		"online":   "N",
	})

	if err != nil {
		fmt.Println("error inserting the document - ", err)
		return err
	}

	fmt.Printf("Successfully inserted into the database with ID: %v\n", insertedId.InsertedID.(primitive.ObjectID).Hex())

	// Now change the status of the user.
	err = dbhelper.UpdateUserStatus(user.Email, "Y")
	if err != nil {
		fmt.Println("Failed to update the status of the user", err)
		return err
	}
	return nil
}

func (dbhelper *DBHelper) UpdateUserStatus(userId string, status string) error {

	// docId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	fmt.Println("error getting docId from the userId", err)
	// 	return err
	// }

	_, err := dbhelper.UserCollection.UpdateOne(context.TODO(), bson.M{"email": userId}, bson.M{"$set": bson.M{"online": status}})
	if err != nil {
		fmt.Println("Request failed to update the status of the user ", err)
		return err
	}
	return nil

}

func (dbhelper *DBHelper) GetUserbyEmailId(emailID string) (models.User, error) {

	var user models.User

	err := dbhelper.UserCollection.FindOne(context.TODO(), bson.M{"email": emailID}).Decode(&user)
	if err != nil {
		fmt.Println("Error getting user by user email id", err)
		return user, err
	}

	return user, nil
}

func (dbhelper *DBHelper) CheckUserExists(emailId string) (bool, error) {
	count, err := dbhelper.UserCollection.CountDocuments(context.TODO(), bson.M{"email": emailId})
	if err != nil {
		fmt.Println("Error checking user in the database", err)
		return false, err
	}
	return count > 0, nil
}

func (dbhelper *DBHelper) GetUserStatus(emailId string) (models.User, error) {

	var user models.User
	err := dbhelper.UserCollection.FindOne(context.TODO(), bson.M{"email": emailId}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New("user not found")
		}
		fmt.Println("Error fetching user status", err)
		return user, err
	}

	return user, nil
}

func (dbhelper *DBHelper) GetOnlineUsers(userId string) ([]models.User, error) {
	// Define the filter for online users
	filter := bson.M{"online": "Y", "email": bson.M{
		"$ne": userId,
	}}

	// Query the collection
	cursor, err := dbhelper.UserCollection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error finding online users:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and decode users
	var onlineUsers []models.User
	for cursor.Next(context.TODO()) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			fmt.Println("Error decoding user:", err)
			return nil, err
		}
		onlineUsers = append(onlineUsers, user)
	}

	if err := cursor.Err(); err != nil {
		fmt.Println("Cursor error:", err)
		return nil, err
	}

	return onlineUsers, nil
}

func (dbhelper *DBHelper) GetConversationBetweenTwoUsers(toUserId, fromUserId string) []models.ConversationStruct {
	var conversations []models.ConversationStruct

	queryCondition := bson.M{
		"$or": []bson.M{
			{
				"$and": []bson.M{
					{
						"toUserID": toUserId,
					},
					{
						"fromUserID": fromUserId,
					},
				},
			},
			{
				"$and": []bson.M{
					{
						"toUserID": fromUserId,
					},
					{
						"fromUserID": toUserId,
					},
				},
			},
		},
	}

	cursor, err := dbhelper.ConversationCollection.Find(context.TODO(), queryCondition)

	if err != nil {
		return conversations
	}

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var conversation models.ConversationStruct
		err := cursor.Decode(&conversation)

		if err == nil {
			conversations = append(conversations, models.ConversationStruct{
				ID:         conversation.ID,
				FromUserID: conversation.FromUserID,
				ToUserID:   conversation.ToUserID,
				Message:    conversation.Message,
			})
		}
	}
	return conversations
}

func (dbhelper *DBHelper) StoreNewChatMessages(messagePayload models.MessagePayloadStruct) bool {
	_, err := dbhelper.ConversationCollection.InsertOne(context.TODO(), bson.M{
		"fromUserID": messagePayload.FromUserID,
		"message":    messagePayload.Message,
		"toUserID":   messagePayload.ToUserID,
	})

	if err != nil {
		fmt.Println("Error inserting the message struct ", err)
		return false
	}
	return true
}
