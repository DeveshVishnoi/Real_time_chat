package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"realtime_chat/models"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func unRegisterAndCloseConnection(c *Client) {
	c.hub.Unregister <- c
	c.webSocketConnection.Close()
}

func setSocketPayloadReadConfig(c *Client) {
	c.webSocketConnection.SetReadLimit(maxMessageSize)
	c.webSocketConnection.SetReadDeadline(time.Now().Add(pongWait))
	c.webSocketConnection.SetPongHandler(func(string) error { c.webSocketConnection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
}

func CreateNewSocketUser(hub *Hub, connection *websocket.Conn, emailID string) {
	client := &Client{
		hub:                 hub,
		webSocketConnection: connection,
		userId:              emailID,
		send:                make(chan SendData),
	}

	fmt.Println("Inside the new socket user")
	go client.ReadPump()
	go client.WritePump()

	client.hub.Register <- client

	fmt.Println(&client, client, *client)
}

func HandleUserRegisterEvent(hub *Hub, client *Client) {

	fmt.Println("Inside the register")
	hub.Clients[client] = true

	hub.HandleSocketData(client, SendData{
		MessageType: models.UserConnectWithChatRoom,
		Message:     client.userId,
	})

}

func HandleUserDisconnectEvent(hub *Hub, client *Client) {

	fmt.Println("Inside the Disconnect")
	_, ok := hub.Clients[client]
	if ok {
		delete(hub.Clients, client)
		close(client.send)

		hub.HandleSocketData(client, SendData{
			MessageType: models.UserDisconnectWithChatRoom,
			Message:     client.userId,
		})
	}

}

func (client *Client) ReadPump() {

	fmt.Println("Inside the readpump ")
	var socketEventPayload SendData

	// Unregistering the client and closing the connection
	defer unRegisterAndCloseConnection(client)

	// Setting up the Payload configuration
	setSocketPayloadReadConfig(client)

	for {
		// ReadMessage is a helper method for getting a reader using NextReader and reading from that reader to a buffer.
		_, payload, err := client.webSocketConnection.ReadMessage()

		decoder := json.NewDecoder(bytes.NewReader(payload))
		decoderErr := decoder.Decode(&socketEventPayload)

		if decoderErr != nil {
			log.Printf("error: %v", decoderErr)
			break
		}

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ===: %v", err)
			}
			break
		}

		//  Getting the proper Payload to send the client
		client.hub.HandleSocketData(client, socketEventPayload)
	}
}

func (client *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.webSocketConnection.Close()
	}()
	for {
		select {
		case payload, ok := <-client.send:

			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(payload)
			finalPayload := reqBodyBytes.Bytes()

			client.webSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.webSocketConnection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.webSocketConnection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(finalPayload)

			n := len(client.send)
			for i := 0; i < n; i++ {
				json.NewEncoder(reqBodyBytes).Encode(<-client.send)
				w.Write(reqBodyBytes.Bytes())
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.webSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.webSocketConnection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (hub *Hub) HandleSocketData(client *Client, data SendData) {

	fmt.Println("Data insid the ", data)
	type chatlistResponseStruct struct {
		Type     string      `json:"type"`
		Chatlist interface{} `json:"chatlist"`
	}

	fmt.Println("Inside the Handle socket data", data)
	switch data.MessageType {
	case models.UserConnectWithChatRoom:

		userid := (data.Message).(string)
		userDetails, err := hub.database.GetUserbyEmailId(userid)
		if err != nil {
			fmt.Println("unmable to find the user bu email id : ", userid)
		}
		if userDetails.Online == "N" {
			log.Println("A logged out user with userID " + userid + " tried to connect to Chat Server.")

		} else {
			newUserOnlinePayload := SendData{
				MessageType: "chatlist-response",
				Message: chatlistResponseStruct{
					Type: "new-user-joined",
					Chatlist: UserDetailsResponsePayloadStruct{
						Online:   userDetails.Online,
						EmailID:  userDetails.Email,
						Username: userDetails.UserName,
					},
				},
			}
			BroadcastSocketEventToAllClientExceptMe(client.hub, newUserOnlinePayload, userDetails.Email)

			users, err := hub.database.GetOnlineUsers(userid)
			if err != nil {
				fmt.Println("error getting all users", err)
			}
			allOnlineUsersPayload := SendData{
				MessageType: "chatlist-response",
				Message: chatlistResponseStruct{
					Type:     "my-chat-list",
					Chatlist: users,
				},
			}
			EmitToSpecificClient(client.hub, allOnlineUsersPayload, userDetails.Email)
		}

	case models.UserDisconnectWithChatRoom:

		fmt.Println("Hello insid ethe dosconnect chat room")
		if data.Message != nil {
			userid := data.Message.(string)
			userDetails, err := hub.database.GetUserbyEmailId(userid)
			if err != nil {
				fmt.Println("unmable to find the user bu email id : ", userid)
			}

			err = hub.database.UpdateUserStatus(userDetails.Email, "N")
			fmt.Println("Unable to update the user online status ", err)

			BroadcastSocketEventToAllClient(client.hub, SendData{
				MessageType: "chatlist-response",
				Message: chatlistResponseStruct{
					Type: "user-disconnected",
					Chatlist: UserDetailsResponsePayloadStruct{
						Online:   "N",
						EmailID:  userDetails.Email,
						Username: userDetails.UserName,
					},
				},
			})

		}

	case models.Message:
		fmt.Println("Hello bay ")
		fmt.Println("data", data)

		message := (data.Message.(map[string]interface{})["message"]).(string)
		fromUserID := (data.Message.(map[string]interface{})["fromUserID"]).(string)
		toUserID := (data.Message.(map[string]interface{})["toUserID"]).(string)

		if message != "" && fromUserID != "" && toUserID != "" {

			messagePacket := models.MessagePayloadStruct{
				FromUserID: fromUserID,
				Message:    message,
				ToUserID:   toUserID,
			}
			hub.database.StoreNewChatMessages(messagePacket)
			allOnlineUsersPayload := SendData{
				MessageType: "message-response",
				Message:     messagePacket,
			}
			EmitToSpecificClient(client.hub, allOnlineUsersPayload, toUserID)

		}
	}

}

// BroadcastSocketEventToAllClient will emit the socket events to all socket users
func BroadcastSocketEventToAllClient(hub *Hub, payload SendData) {
	for client := range hub.Clients {
		select {
		case client.send <- payload:
		default:
			close(client.send)
			delete(hub.Clients, client)
		}
	}
}

// BroadcastSocketEventToAllClientExceptMe will emit the socket events to all socket users,
// except the user who is emitting the event
func BroadcastSocketEventToAllClientExceptMe(hub *Hub, payload SendData, myUserID string) {
	for client := range hub.Clients {
		if client.userId != myUserID {
			select {
			case client.send <- payload:
			default:
				close(client.send)
				delete(hub.Clients, client)
			}
		}
	}
}

// EmitToSpecificClient will emit the socket event to specific socket user
func EmitToSpecificClient(hub *Hub, payload SendData, userID string) {

	fmt.Println("hub", hub)
	fmt.Println("paylod", payload)
	fmt.Println("UserId", userID)

	for client := range hub.Clients {
		if client.userId == userID {
			select {
			case client.send <- payload:
			default:
				close(client.send)
				delete(hub.Clients, client)
			}
		}
	}
}
