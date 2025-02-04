package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/M1keTrike/practicaClienteServido/mainsServer/storage"
	"github.com/gin-gonic/gin"
)

func CheckChanges(c *gin.Context) {
	hasChanges := storage.HasChanges()
	c.JSON(http.StatusOK, gin.H{"changes": hasChanges})
}

func LongPollingReplication(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	for {
		time.Sleep(2 * time.Second)

		if user, found := storage.GetInsertedUser(); found {
			json.NewEncoder(c.Writer).Encode(user)
			c.Writer.Flush()
			return
		}
	}
}

func GetUpdatedUsers(c *gin.Context) {
	users := storage.GetUpdatedUsers()
	c.JSON(http.StatusOK, users)
}

func GetDeletedUsers(c *gin.Context) {
	ids := storage.GetDeletedUsers()
	c.JSON(http.StatusOK, ids)
}
