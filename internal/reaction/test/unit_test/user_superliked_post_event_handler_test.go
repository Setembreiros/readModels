package reaction_test

import (
	"encoding/json"
	"readmodels/internal/model"
	reaction_handler "readmodels/internal/reaction/handler"
	mock_reaction_handler "readmodels/internal/reaction/handler/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

var userSuperlikedPostEventService *mock_reaction_handler.MockUserSuperlikedPostEventService
var userSuperlikedPostEventHandler *reaction_handler.UserSuperlikedPostEventHandler

func setUpUserSuperlikedPostEventHandler(t *testing.T) {
	SetUp(t)
	userSuperlikedPostEventService = mock_reaction_handler.NewMockUserSuperlikedPostEventService(ctrl)
	userSuperlikedPostEventHandler = reaction_handler.NewUserSuperlikedPostEventHandler(userSuperlikedPostEventService)
}

func TestHandleUserSuperlikedPostEvent(t *testing.T) {
	setUpUserSuperlikedPostEventHandler(t)
	data := &reaction_handler.UserSuperlikedPostEvent{
		Username: "user123",
		PostId:   "post123",
	}
	event, _ := json.Marshal(data)
	expectedPostSuperlike := &model.PostSuperlike{
		User: &model.UserMetadata{
			Username: "user123",
		},
		PostId: "post123",
	}
	userSuperlikedPostEventService.EXPECT().CreatePostSuperlike(expectedPostSuperlike)

	userSuperlikedPostEventHandler.Handle(event)
}

func TestInvalidDataInUserSuperlikedPostEventHandler(t *testing.T) {
	setUpUserSuperlikedPostEventHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	userSuperlikedPostEventHandler.Handle(event)

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
