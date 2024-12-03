package handlers

import (
	"realtime_chat/api/database"

	"github.com/gorilla/websocket"
)

type SendData struct {
	MessageType string      `json:"message_type" bson:"message_type"`
	Message     interface{} `json:"message" bson:"message"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub                 *Hub
	webSocketConnection *websocket.Conn
	userId              string
	send                chan SendData
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {

	// Registerd client
	Clients map[*Client]bool

	// Register request from the client
	Register chan *Client

	// Unregister request from the client
	Unregister chan *Client

	database database.DBHelperProvider
}

type UserDetailsResponsePayloadStruct struct {
	Username string `json:"username"`
	EmailID  string `json:"email"`
	Online   string `json:"online"`
}
