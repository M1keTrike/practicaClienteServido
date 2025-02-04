package handlers

import (
	"net/http"
	"strconv"

	"github.com/M1keTrike/practicaClienteServido/entities"
	"github.com/M1keTrike/practicaClienteServido/mainsServer/storage"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var newUser entities.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	storage.AddUser(newUser)

	c.JSON(http.StatusCreated, newUser)
}

func UpdateUserByID(c *gin.Context) {
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

	success := storage.UpdateUser(id, updatedUser)
	if success {
		c.JSON(http.StatusOK, updatedUser)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
	}
}

func DeleteUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	success := storage.DeleteUser(id)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
	}
}
