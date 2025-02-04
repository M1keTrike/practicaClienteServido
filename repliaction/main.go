package main

import (
	"github.com/M1keTrike/practicaClienteServido/entities"
	"github.com/gin-gonic/gin"
)


type Strorage struct {
	Users []entities.User
}




func main() {
	r := gin.Default()


	r.Run("8081")
}