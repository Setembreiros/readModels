package comment_test

import (
	"bytes"
	mock_comment "readmodels/internal/comment/test/mock"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var loggerOutput bytes.Buffer
var repository *mock_comment.MockRepository

func SetUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository = mock_comment.NewMockRepository(ctrl)
	log.Logger = log.Output(&loggerOutput)
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
