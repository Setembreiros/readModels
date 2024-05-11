package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"readmodels/internal/events"
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

func UserWasRegisteredEventHandler(busChannel <-chan events.Event, ctx context.Context) {
	for {
		select {
		case event := <-busChannel:
			var userWasRegisteredEvent UserWasRegisteredEvent
			err := Decode(event.Data, &userWasRegisteredEvent)
			if err != nil {
				fmt.Printf("Invalid event data, err: %s\n", err)
				continue
			}
			service(userWasRegisteredEvent)
		case <-ctx.Done():
			return
		}
	}
}

func Decode(datab []byte, data *UserWasRegisteredEvent) error {
	return json.Unmarshal(datab, &data)
}
