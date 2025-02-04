package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func MakeGetRequest(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error en la solicitud HTTP:", err)
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
