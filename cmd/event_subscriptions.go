package main

import (
	"context"
	"readmodels/internal/events"
	"readmodels/internal/events/handlers"
)

type subscription struct {
	EventType string
	Handler   handler
}

var Subscriptions []subscription = []subscription{
	{
		EventType: "UserWasRegisteredEvent",
		Handler:   handlers.UserWasRegisteredEventHandler,
	},
}

type handler func(busChannel <-chan events.Event, ctx context.Context)
