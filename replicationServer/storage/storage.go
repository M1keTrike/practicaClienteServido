package storage

import (
	"sync"

	"github.com/M1keTrike/practicaClienteServido/entities"
)

var (
	replicatedUsers []entities.User
	mutex           sync.Mutex
)

func AddUser(user entities.User) {
	mutex.Lock()
	defer mutex.Unlock()
	replicatedUsers = append(replicatedUsers, user)
}

func UpdateUser(updatedUser entities.User) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, user := range replicatedUsers {
		if user.Id == updatedUser.Id {
			replicatedUsers[i] = updatedUser
			break
		}
	}
}

func DeleteUser(id int) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, user := range replicatedUsers {
		if user.Id == id {
			replicatedUsers = append(replicatedUsers[:i], replicatedUsers[i+1:]...)
			break
		}
	}
}

func GetUsers() []entities.User {
	mutex.Lock()
	defer mutex.Unlock()
	return replicatedUsers
}
