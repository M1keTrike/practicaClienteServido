package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/M1keTrike/practicaClienteServido/entities"
	"github.com/M1keTrike/practicaClienteServido/replicationServer/storage"
	"github.com/gin-gonic/gin"
)

func CheckForChanges() {
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
			go ReceiveNewRecords()
			GetUpdatedRecords()
			GetDeletedRecords()
		}

		time.Sleep(5 * time.Second)
	}
}

func ReceiveNewRecords() {
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
			storage.AddUser(user)
			fmt.Println("Nuevo usuario replicado:", user)
		}

		resp.Body.Close()
	}
}

func GetUpdatedRecords() {
	resp, err := http.Get("http://localhost:8081/updated-users")
	if err != nil {
		fmt.Println("Error obteniendo actualizaciones:", err)
		return
	}
	defer resp.Body.Close()

	var updatedUsers []entities.User
	json.NewDecoder(resp.Body).Decode(&updatedUsers)

	for _, updatedUser := range updatedUsers {
		storage.UpdateUser(updatedUser)
	}
}

func GetDeletedRecords() {
	resp, err := http.Get("http://localhost:8081/deleted-users")
	if err != nil {
		fmt.Println("Error obteniendo eliminaciones:", err)
		return
	}
	defer resp.Body.Close()

	var deletedIDs []int
	json.NewDecoder(resp.Body).Decode(&deletedIDs)

	for _, id := range deletedIDs {
		storage.DeleteUser(id)
	}
}

func GetReplicatedUsers(c *gin.Context) {
	users := storage.GetUsers()
	c.JSON(http.StatusOK, users)
}
