package reaction_test

import (
	"encoding/json"
	"readmodels/internal/model"
	reaction_handler "readmodels/internal/reaction/handler"
	mock_reaction_handler "readmodels/internal/reaction/handler/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

var userLikedPostEventService *mock_reaction_handler.MockUserLikedPostEventService
var userLikedPostEventHandler *reaction_handler.UserLikedPostEventHandler

func setUpUserLikedPostEventHandler(t *testing.T) {
	SetUp(t)
	userLikedPostEventService = mock_reaction_handler.NewMockUserLikedPostEventService(ctrl)
	userLikedPostEventHandler = reaction_handler.NewUserLikedPostEventHandler(userLikedPostEventService)
}

func TestHandleUserLikedPostEvent(t *testing.T) {
	setUpUserLikedPostEventHandler(t)
	data := &reaction_handler.UserLikedPostEvent{
		Username: "user123",
		PostId:   "post123",
	}
	event, _ := json.Marshal(data)
	expectedPostLike := &model.PostLike{
		User: &model.UserMetadata{
			Username: "user123",
		},
		PostId: "post123",
	}
	userLikedPostEventService.EXPECT().CreatePostLike(expectedPostLike)

	userLikedPostEventHandler.Handle(event)
}

func TestInvalidDataInUserLikedPostEventHandler(t *testing.T) {
	setUpUserLikedPostEventHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	userLikedPostEventHandler.Handle(event)

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
