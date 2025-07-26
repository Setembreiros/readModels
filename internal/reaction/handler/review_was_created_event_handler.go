package reaction_handler

import (
	common_data "readmodels/internal/common/data"
	"readmodels/internal/model"
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=review_was_created_event_handler.go -destination=test/mock/review_was_created_event_handler.go

type ReviewWasCreatedEvent struct {
	ReviewId  uint64 `json:"reviewId"`
	Username  string `json:"username"`
	PostId    string `json:"postId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Rating    int    `json:"rating"`
	CreatedAt string `json:"createdAt"`
}

type ReviewWasCreatedEventService interface {
	CreateReview(data *model.Review)
}

type ReviewWasCreatedEventHandler struct {
	service ReviewWasCreatedEventService
}

func NewReviewWasCreatedEventHandler(service ReviewWasCreatedEventService) *ReviewWasCreatedEventHandler {
	return &ReviewWasCreatedEventHandler{
		service: service,
	}
}

func (handler *ReviewWasCreatedEventHandler) Handle(event []byte) {
	var reviewWasCreatedEvent ReviewWasCreatedEvent
	log.Info().Msg("Handling ReviewWasCreatedEvent")

	err := common_data.DeserializeData(event, &reviewWasCreatedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	data, err := mapData(reviewWasCreatedEvent)
	if err != nil {
		return
	}

	handler.service.CreateReview(data)
}

func mapData(event ReviewWasCreatedEvent) (*model.Review, error) {
	parsedCreatedAt, err := time.Parse(model.TimeLayout, event.CreatedAt)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error parsing time CreatedAt")
		return nil, err
	}

	return &model.Review{
		ReviewId:  event.ReviewId,
		Username:  event.Username,
		PostId:    event.PostId,
		Title:     event.Title,
		Content:   event.Content,
		Rating:    event.Rating,
		CreatedAt: parsedCreatedAt,
	}, nil
}
