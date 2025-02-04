package main

import (
	"github.com/M1keTrike/practicaClienteServido/replicationServer/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	
	go handlers.CheckForChanges()

	r := gin.Default()


	r.GET("/replicated-users", handlers.GetReplicatedUsers)

	r.Run(":8082")
}