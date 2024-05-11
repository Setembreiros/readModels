package main

import (
	"readmodels/internal/events"
	"readmodels/internal/events/handlers"
)

var Subscriptions []events.EventSubscription = []events.EventSubscription{
	{
		EventType: "UserWasRegisteredEvent",
		Handler:   handlers.UserWasRegisteredEventHandler,
	},
}
