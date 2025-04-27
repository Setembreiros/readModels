package reaction_test

import (
	"bytes"
	"net/http/httptest"
	mock_database "readmodels/internal/db/test/mock"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var ctrl *gomock.Controller
var loggerOutput bytes.Buffer
var client *mock_database.MockDatabaseClient
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func SetUp(t *testing.T) {
	ctrl = gomock.NewController(t)
	client = mock_database.NewMockDatabaseClient(ctrl)
	log.Logger = log.Output(&loggerOutput)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
