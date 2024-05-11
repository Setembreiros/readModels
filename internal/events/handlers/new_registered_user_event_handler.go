package handlers

import (
	"encoding/json"
	"fmt"
)

type UserWasRegisteredEvent struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	Region   string `json:"region"`
	FullName string `json:"full_name"`
}

func service(event UserWasRegisteredEvent) {
	fmt.Println("New user registered:")
	fmt.Println("ID:", event.UserId)
}

func UserWasRegisteredEventHandler(event []byte) {
	var userWasRegisteredEvent UserWasRegisteredEvent

	err := Decode(event, &userWasRegisteredEvent)
	if err != nil {
		fmt.Printf("Invalid event data, err: %s\n", err)
		return
	}

	service(userWasRegisteredEvent)
}

func Decode(datab []byte, data any) error {
	return json.Unmarshal(datab, &data)
}
