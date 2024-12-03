package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"realtime_chat/api/handlers"
	"realtime_chat/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	Engine *gin.Engine
}

func NewServer(apihandle *handlers.ApiHandler) *Server {

	hub := apihandle.NewHub()
	go hub.Run()

	router := gin.Default()

	corsMiddleware := utils.GetCorsConfig()
	router.Use(func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// Here we declare some routes.
	api := router.Group("/api")

	// for the root page.
	api.GET("/", apihandle.RenderHome)
	api.POST("/login", apihandle.UserLogin)
	api.POST("/reg", apihandle.Registration)

	// TODO: Please make sure this is query parameter.

	// http://localhost:8080/api/isAvailable?emailId=deveshvishnoi@suju.co
	api.GET("/isAvailable/:username", apihandle.IsUserAvailable)

	api.GET("/sessionStatus/:emailId", apihandle.UserSessionCheck)

	api.GET("/getConversation/:toUserId/:fromUserId", apihandle.GetMessagesHandler)

	api.GET("/ws/:emailId", func(c *gin.Context) {

		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		userID := c.Param("emailId")

		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		// Upgrade HTTP connection to WebSocket
		connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection for user %s: %v", userID, err)
			return
		}

		// Handle the new WebSocket user
		handlers.CreateNewSocketUser(hub, connection, userID)
	})

	return &Server{
		Engine: router,
	}
}

func (s *Server) Start() {

	log.Printf("Starting the HTTP Server: %s\n", os.Getenv("GIN_PORT"))

	if err := s.Engine.Run(os.Getenv("GIN_PORT")); err != nil && err != http.ErrServerClosed {
		fmt.Println("Error getting run the server : ", err)
		panic(err)
	}
}
