package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/M1keTrike/practicaClienteServido/entities"
	"github.com/gin-gonic/gin"
)

var (
	replicatedUsers []entities.User
	mutex           sync.Mutex
)

func checkForChanges() {
	for {
		resp, err := http.Get("http://localhost:8081/changes")
		if err != nil {
			fmt.Println("Error en short polling:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		var result map[string]bool
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if result["changes"] {
			
			go receiveNewRecords()

			
			getUpdatedRecords()

		
			getDeletedRecords()
		}

		time.Sleep(5 * time.Second)
	}
}


func receiveNewRecords() {
	for {
		resp, err := http.Get("http://localhost:8081/longpolling")
		if err != nil {
			fmt.Println("Error en long polling:", err)
			time.Sleep(2 * time.Second) 
			continue
		}

		decoder := json.NewDecoder(resp.Body)
		var user entities.User
		if err := decoder.Decode(&user); err == nil {
			mutex.Lock()
			replicatedUsers = append(replicatedUsers, user)
			mutex.Unlock()
			fmt.Println("Nuevo usuario replicado:", user)
		}

		resp.Body.Close()
	}
}


func getUpdatedRecords() {
	resp, err := http.Get("http://localhost:8081/updated-users")
	if err != nil {
		fmt.Println("Error obteniendo actualizaciones:", err)
		return
	}
	defer resp.Body.Close()

	var updatedUsers []entities.User
	json.NewDecoder(resp.Body).Decode(&updatedUsers)

	mutex.Lock()
	for i, user := range replicatedUsers {
		for _, updated := range updatedUsers {
			if user.Id == updated.Id {
				replicatedUsers[i] = updated
			}
		}
	}
	mutex.Unlock()
}


func getDeletedRecords() {
	resp, err := http.Get("http://localhost:8081/deleted-users")
	if err != nil {
		fmt.Println("Error obteniendo eliminaciones:", err)
		return
	}
	defer resp.Body.Close()

	var deletedIDs []int
	json.NewDecoder(resp.Body).Decode(&deletedIDs)

	mutex.Lock()
	for _, id := range deletedIDs {
		for i, user := range replicatedUsers {
			if user.Id == id {
				replicatedUsers = append(replicatedUsers[:i], replicatedUsers[i+1:]...)
				break
			}
		}
	}
	mutex.Unlock()
}

func main() {
	go checkForChanges()

	r := gin.Default()
	r.GET("/replicated-users", func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()
		c.JSON(http.StatusOK, replicatedUsers)
	})

	r.Run(":8082")
}
