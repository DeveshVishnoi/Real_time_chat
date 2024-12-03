package main

import (
	"fmt"
	db "realtime_chat/api/database"
	"realtime_chat/api/handlers"
	server "realtime_chat/api/routes"
	"realtime_chat/pkg/config"
)

func main() {

	config.LoadEnvData()

	client := db.ConnectDatabase()
	fmt.Println("Client -- ", client)

	dbHelper := db.NewDbProvider(client.MongoClient)
	handler := handlers.NewHandlerProvider(dbHelper)
	server := server.NewServer(handler)

	server.Start()

}
