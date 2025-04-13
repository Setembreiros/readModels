package integration_test_followers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	database "readmodels/internal/db"
	"readmodels/internal/follow"
	"readmodels/internal/userprofile"
	integration_test_arrange "readmodels/test/integration_test_common/arrange"
	integration_test_assert "readmodels/test/integration_test_common/assert"
	"testing"

	"github.com/gin-gonic/gin"
)

var db *database.Database
var controller *follow.FollowController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUp(t *testing.T) {
	// Mocks
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)

	// Real infrastructure and services
	db = integration_test_arrange.CreateTestDatabase(t, ginContext)
	repository := follow.FollowRepository(*db)
	service := follow.NewFollowService(repository)
	controller = follow.NewFollowController(service)
}

func tearDown() {
	db.Client.Clean()
}

func TestGetFollowersMetadata_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	ginContext.Request, _ = http.NewRequest("GET", "/followers", nil)
	followerId1 := "USERA"
	followerId2 := "USERB"
	followerId3 := "USERC"
	u := url.Values{}
	u.Add("followerId", followerId1)
	u.Add("followerId", followerId2)
	u.Add("followerId", followerId3)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedData := []follow.FollowerMetadata{
		{
			Username: followerId1,
			Name:     "fullname1",
		},
		{
			Username: followerId2,
			Name:     "fullname2",
		},
		{
			Username: followerId3,
			Name:     "fullname3",
		},
	}
	populateDb(t, expectedData)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {"followers":[
		{
			"username":      "` + followerId1 + `",
			"fullname":   "` + expectedData[0].Name + `"
		},
		{
			"username":      "` + followerId2 + `",
			"fullname":   "` + expectedData[1].Name + `"
		},
		{
			"username":      "` + followerId3 + `",
			"fullname":   "` + expectedData[2].Name + `"
		}
		]
		}
	}`

	controller.GetFollowersMetadata(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}

func TestGetFolloweesMetadata_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	ginContext.Request, _ = http.NewRequest("GET", "/followees", nil)
	followeeId1 := "USERA"
	followeeId2 := "USERB"
	followeeId3 := "USERC"
	u := url.Values{}
	u.Add("followeeId", followeeId1)
	u.Add("followeeId", followeeId2)
	u.Add("followeeId", followeeId3)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedData := []follow.FollowerMetadata{
		{
			Username: followeeId1,
			Name:     "fullname1",
		},
		{
			Username: followeeId2,
			Name:     "fullname2",
		},
		{
			Username: followeeId3,
			Name:     "fullname3",
		},
	}
	populateDb(t, expectedData)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {"followees":[
		{
			"username":      "` + followeeId1 + `",
			"fullname":   "` + expectedData[0].Name + `"
		},
		{
			"username":      "` + followeeId2 + `",
			"fullname":   "` + expectedData[1].Name + `"
		},
		{
			"username":      "` + followeeId3 + `",
			"fullname":   "` + expectedData[2].Name + `"
		}
		]
		}
	}`

	controller.GetFolloweesMetadata(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}

func populateDb(t *testing.T, data []follow.FollowerMetadata) {
	for _, follower := range data {
		integration_test_arrange.AddUserProfileToDatabase(t, db, &userprofile.UserProfile{
			Username: follower.Username,
			Name:     follower.Name,
		})

	}
}
