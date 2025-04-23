package reaction_test

import (
	"encoding/json"
	"readmodels/internal/model"
	reaction_handler "readmodels/internal/reaction/handler"
	mock_reaction_handler "readmodels/internal/reaction/handler/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

var userUnsuperlikedPostEventService *mock_reaction_handler.MockUserUnsuperlikedPostEventService
var userUnsuperlikedPostEventHandler *reaction_handler.UserUnsuperlikedPostEventHandler

func setUpUserUnsuperlikedPostEventHandler(t *testing.T) {
	SetUp(t)
	userUnsuperlikedPostEventService = mock_reaction_handler.NewMockUserUnsuperlikedPostEventService(ctrl)
	userUnsuperlikedPostEventHandler = reaction_handler.NewUserUnsuperlikedPostEventHandler(userUnsuperlikedPostEventService)
}

func TestHandleUserUnsuperlikedPostEvent(t *testing.T) {
	setUpUserUnsuperlikedPostEventHandler(t)
	data := &reaction_handler.UserUnsuperlikedPostEvent{
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
	userUnsuperlikedPostEventService.EXPECT().DeletePostSuperlike(expectedPostSuperlike)

	userUnsuperlikedPostEventHandler.Handle(event)
}

func TestInvalidDataInUserUnsuperlikedPostEventHandler(t *testing.T) {
	setUpUserUnsuperlikedPostEventHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	userUnsuperlikedPostEventHandler.Handle(event)

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
