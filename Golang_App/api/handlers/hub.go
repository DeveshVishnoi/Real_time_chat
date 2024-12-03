package handlers

import "fmt"

func (apihandle *ApiHandler) NewHub() *Hub {

	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		database:   apihandle.DBHelperProvider,
	}
}

// Run will execute Go Routines to check incoming Socket events
func (hub *Hub) Run() {

	fmt.Println("Hello")
	for {
		select {
		case client := <-hub.Register:
			HandleUserRegisterEvent(hub, client)

		case client := <-hub.Unregister:
			HandleUserDisconnectEvent(hub, client)
		}
	}
}
