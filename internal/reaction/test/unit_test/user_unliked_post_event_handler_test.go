package reaction_test

import (
	"encoding/json"
	"readmodels/internal/model"
	reaction_handler "readmodels/internal/reaction/handler"
	mock_reaction_handler "readmodels/internal/reaction/handler/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

var userUnlikedPostEventService *mock_reaction_handler.MockUserUnlikedPostEventService
var userUnlikedPostEventHandler *reaction_handler.UserUnlikedPostEventHandler

func setUpUserUnlikedPostEventHandler(t *testing.T) {
	SetUp(t)
	userUnlikedPostEventService = mock_reaction_handler.NewMockUserUnlikedPostEventService(ctrl)
	userUnlikedPostEventHandler = reaction_handler.NewUserUnlikedPostEventHandler(userUnlikedPostEventService)
}

func TestHandleUserUnlikedPostEvent(t *testing.T) {
	setUpUserUnlikedPostEventHandler(t)
	data := &reaction_handler.UserUnlikedPostEvent{
		Username: "user123",
		PostId:   "post123",
	}
	event, _ := json.Marshal(data)
	expectedLikePost := &model.LikePost{
		Username: "user123",
		PostId:   "post123",
	}
	userUnlikedPostEventService.EXPECT().DeleteLikePost(expectedLikePost)

	userUnlikedPostEventHandler.Handle(event)
}

func TestInvalidDataInUserUnlikedPostEventHandler(t *testing.T) {
	setUpUserUnlikedPostEventHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	userUnlikedPostEventHandler.Handle(event)

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
