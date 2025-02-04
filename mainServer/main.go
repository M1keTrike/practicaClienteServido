package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/M1keTrike/practicaClienteServido/entities"
	"github.com/gin-gonic/gin"
)

var (
	users  []entities.User
	mutex  sync.Mutex
	nextID = 1
)

func createUser(c *gin.Context) {
	var newUser entities.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	mutex.Lock()
	newUser.Id = nextID
	nextID++
	users = append(users, newUser)
	mutex.Unlock()

	c.JSON(http.StatusCreated, newUser)
}



func updateUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var updatedUser entities.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, user := range users {
		if user.Id == id {
			updatedUser.Id = id
			users[i] = updatedUser
			c.JSON(http.StatusOK, updatedUser)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
}

func deleteUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
}

func main() {
	r := gin.Default()

	
	r.POST("/users", createUser)      

	r.PUT("/users/:id", updateUserByID) 
	r.DELETE("/users/:id", deleteUserByID)

	r.Run(":8081") 
}
