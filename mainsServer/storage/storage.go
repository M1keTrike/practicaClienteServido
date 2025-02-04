package storage

import (
	"sync"

	"github.com/M1keTrike/practicaClienteServido/entities"
)

var (
	users         []entities.User
	mutex         sync.Mutex
	nextID        = 1
	insertedUsers []entities.User
	updatedUsers  []entities.User
	deletedUsers  []int
)

func AddUser(user entities.User) {
	mutex.Lock()
	defer mutex.Unlock()

	user.Id = nextID
	nextID++
	user.OperationStatus = "INSERT"

	users = append(users, user)
	insertedUsers = append(insertedUsers, user)
}

func UpdateUser(id int, updatedUser entities.User) bool {
	mutex.Lock()
	defer mutex.Unlock()

	for i, user := range users {
		if user.Id == id {
			updatedUser.Id = id
			updatedUser.OperationStatus = "UPDATE"
			users[i] = updatedUser
			updatedUsers = append(updatedUsers, updatedUser)
			return true
		}
	}
	return false
}

func DeleteUser(id int) bool {
	mutex.Lock()
	defer mutex.Unlock()

	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)
			deletedUsers = append(deletedUsers, id)
			return true
		}
	}
	return false
}

func HasChanges() bool {
	mutex.Lock()
	defer mutex.Unlock()

	return len(insertedUsers) > 0 || len(updatedUsers) > 0 || len(deletedUsers) > 0
}

func GetInsertedUser() (entities.User, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	if len(insertedUsers) > 0 {
		user := insertedUsers[0]
		insertedUsers = insertedUsers[1:]
		return user, true
	}
	return entities.User{}, false
}

func GetUpdatedUsers() []entities.User {
	mutex.Lock()
	defer mutex.Unlock()

	usersCopy := updatedUsers
	updatedUsers = []entities.User{}
	return usersCopy
}

func GetDeletedUsers() []int {
	mutex.Lock()
	defer mutex.Unlock()

	idsCopy := deletedUsers
	deletedUsers = []int{}
	return idsCopy
}
