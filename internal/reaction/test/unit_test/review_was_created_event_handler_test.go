package reaction_test

import (
	"encoding/json"
	"readmodels/internal/model"
	reaction_handler "readmodels/internal/reaction/handler"
	mock_reaction_handler "readmodels/internal/reaction/handler/test/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var reviewWasCreatedEventHandler *reaction_handler.ReviewWasCreatedEventHandler
var reviewWasCreatedEventService *mock_reaction_handler.MockReviewWasCreatedEventService

func setUpHandler(t *testing.T) {
	SetUp(t)
	reviewWasCreatedEventService = mock_reaction_handler.NewMockReviewWasCreatedEventService(ctrl)
	reviewWasCreatedEventHandler = reaction_handler.NewReviewWasCreatedEventHandler(reviewWasCreatedEventService)
}

func TestHandleReviewWasCreatedEvent(t *testing.T) {
	setUpHandler(t)
	timeNow := time.Now().UTC().Format(model.TimeLayout)
	data := &reaction_handler.ReviewWasCreatedEvent{
		ReviewId:  uint64(123456),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		Rating:    3,
		CreatedAt: timeNow,
	}
	event, _ := json.Marshal(data)
	expectedTime, _ := time.Parse(model.TimeLayout, timeNow)
	expectedReview := &model.Review{
		ReviewId:  uint64(123456),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		Rating:    3,
		CreatedAt: expectedTime,
	}
	reviewWasCreatedEventService.EXPECT().CreateReview(expectedReview)

	reviewWasCreatedEventHandler.Handle(event)
}

func TestInvalidDataInReviewWasCreatedEventHandler(t *testing.T) {
	setUpHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	reviewWasCreatedEventHandler.Handle(event)

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
