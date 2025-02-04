package entities

type User struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Username         string `json:"username"`
	OperationStatus  string `json:"operation_status"`
}
