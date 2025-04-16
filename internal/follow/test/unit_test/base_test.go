package unit_test_follow

import (
	"bytes"
	mock_database "readmodels/internal/db/test/mock"
	mock_follow "readmodels/internal/follow/test/mock"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var client *mock_database.MockDatabaseClient
var loggerOutput bytes.Buffer
var repository *mock_follow.MockRepository

func setUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	client = mock_database.NewMockDatabaseClient(ctrl)
	repository = mock_follow.NewMockRepository(ctrl)
	log.Logger = log.Output(&loggerOutput)
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
