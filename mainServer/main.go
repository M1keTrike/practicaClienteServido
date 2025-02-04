package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/M1keTrike/practicaClienteServido/entities"
	"github.com/gin-gonic/gin"
)

var (
	users         []entities.User
	mutex         sync.Mutex
	nextID        = 1
	insertedUsers []entities.User
	updatedUsers  []entities.User
	deletedUsers  []int
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
	newUser.OperationStatus = "INSERT"
	users = append(users, newUser)
	insertedUsers = append(insertedUsers, newUser) 


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
			updatedUser.OperationStatus = "UPDATE"
			users[i] = updatedUser
			updatedUsers = append(updatedUsers, updatedUser) 
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
			deletedUsers = append(deletedUsers, id) 
			c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
}

func checkChanges(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	hasChanges := len(insertedUsers) > 0 || len(updatedUsers) > 0 || len(deletedUsers) > 0

	c.JSON(http.StatusOK, gin.H{"changes": hasChanges})
}

func longPollingReplication(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	for {
		time.Sleep(2 * time.Second)

		mutex.Lock()
		if len(insertedUsers) > 0 {
			user := insertedUsers[0]
			insertedUsers = insertedUsers[1:] 
			mutex.Unlock()

		
			json.NewEncoder(c.Writer).Encode(user)
			c.Writer.Flush()
			return 
		}
		mutex.Unlock()
	}
}


func getUpdatedUsers(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	c.JSON(http.StatusOK, updatedUsers)
	updatedUsers = []entities.User{}
}

func getDeletedUsers(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	c.JSON(http.StatusOK, deletedUsers)
	deletedUsers = []int{}
}

func main() {
	r := gin.Default()

	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUserByID)
	r.DELETE("/users/:id", deleteUserByID)

	r.GET("/changes", checkChanges)
	r.GET("/longpolling", longPollingReplication)
	r.GET("/updated-users", getUpdatedUsers)
	r.GET("/deleted-users", getDeletedUsers)

	r.Run(":8081")
}
