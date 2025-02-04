package main

import (
	"github.com/M1keTrike/practicaClienteServido/mainsServer/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/users", handlers.CreateUser)
	r.PUT("/users/:id", handlers.UpdateUserByID)
	r.DELETE("/users/:id", handlers.DeleteUserByID)

	r.GET("/changes", handlers.CheckChanges)
	r.GET("/longpolling", handlers.LongPollingReplication)
	r.GET("/updated-users", handlers.GetUpdatedUsers)
	r.GET("/deleted-users", handlers.GetDeletedUsers)

	r.Run(":8081") 
}
