package integration_test_reaction

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	database "readmodels/internal/db"
	"readmodels/internal/model"
	"readmodels/internal/reaction"
	reaction_handler "readmodels/internal/reaction/handler"
	integration_test_arrange "readmodels/test/integration_test_common/arrange"
	integration_test_assert "readmodels/test/integration_test_common/assert"
	"readmodels/test/test_common"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

var db *database.Database
var controller *reaction.ReactionController
var userLikedPostEventHandler *reaction_handler.UserLikedPostEventHandler
var userSuperlikedPostEventHandler *reaction_handler.UserSuperlikedPostEventHandler
var userUnlikedPostEventHandler *reaction_handler.UserUnlikedPostEventHandler
var userUnsuperlikedPostEventHandler *reaction_handler.UserUnsuperlikedPostEventHandler
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUp(t *testing.T) {
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)

	// Real infrastructure and services
	db = integration_test_arrange.CreateTestDatabase(t, ginContext)
	repository := reaction.NewReactionRepository(db)
	service := reaction.NewReactionService(repository)
	controller = reaction.NewReactionController(service)
	userLikedPostEventHandler = reaction_handler.NewUserLikedPostEventHandler(service)
	userSuperlikedPostEventHandler = reaction_handler.NewUserSuperlikedPostEventHandler(service)
	userUnlikedPostEventHandler = reaction_handler.NewUserUnlikedPostEventHandler(service)
	userUnsuperlikedPostEventHandler = reaction_handler.NewUserUnsuperlikedPostEventHandler(service)
}

func tearDown() {
	db.Client.Truncate()
}

func TestCreatePostLike_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	existingPost := &database.PostMetadata{
		PostId:   "post123",
		Username: "username1",
		Type:     "TEXT",
		Likes:    0,
	}
	existingUserProfile := &model.UserProfile{
		Username: "user123",
		Name:     "fullname123",
	}
	integration_test_arrange.AddPostToDatabase(t, db, existingPost)
	integration_test_arrange.AddUserProfileToDatabase(t, db, existingUserProfile)
	data := &reaction_handler.UserLikedPostEvent{
		Username: "user123",
		PostId:   existingPost.PostId,
	}
	event, _ := test_common.SerializeData(data)
	expectedLike := &model.PostLike{
		User: &model.UserMetadata{
			Username: data.Username,
			Name:     existingUserProfile.Name,
		},
		PostId: data.PostId,
	}

	userLikedPostEventHandler.Handle(event)

	integration_test_assert.AssertPostLikeExists(t, db, expectedLike)
	integration_test_assert.AssertPostLikesIncreased(t, db, existingPost.PostId)
}

func TestCreatePostSuperlike_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	existingPost := &database.PostMetadata{
		PostId:     "post123",
		Username:   "username1",
		Type:       "TEXT",
		Superlikes: 0,
	}
	existingUserProfile := &model.UserProfile{
		Username: "user123",
		Name:     "fullname123",
	}
	integration_test_arrange.AddPostToDatabase(t, db, existingPost)
	integration_test_arrange.AddUserProfileToDatabase(t, db, existingUserProfile)
	data := &reaction_handler.UserSuperlikedPostEvent{
		Username: "user123",
		PostId:   existingPost.PostId,
	}
	event, _ := test_common.SerializeData(data)
	expectedSuperlike := &model.PostSuperlike{
		User: &model.UserMetadata{
			Username: data.Username,
			Name:     existingUserProfile.Name,
		},
		PostId: data.PostId,
	}

	userSuperlikedPostEventHandler.Handle(event)

	integration_test_assert.AssertPostSuperlikeExists(t, db, expectedSuperlike)
	integration_test_assert.AssertPostSuperlikesIncreased(t, db, existingPost.PostId)
}

func TestGetPostLikesMetadata_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	populatePostLikesDb(t)
	postId := "post1"
	lastUsername := "username2"
	limit := 4
	ginContext.Request, _ = http.NewRequest("GET", "/postLikes", nil)
	ginContext.Params = []gin.Param{{Key: "postId", Value: postId}}
	u := url.Values{}
	u.Add("lastUsername", lastUsername)
	u.Add("limit", strconv.Itoa(limit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"postLikes":[	
			{
				"username":  "username3",
				"name": 	 "fullname3"
			},		
			{
				"username":  "username4",
				"name": 	 "fullname4"
			},		
			{
				"username":  "username5",
				"name": 	 "fullname5"
			},		
			{
				"username":  "username6",
				"name": 	 "fullname6"
			}
			],
			"lastUsername":"username6"
		}
	}`

	controller.GetPostLikesMetadata(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}

func TestGetPostSuperlikesMetadata_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	populatePostSuperlikesDb(t)
	postId := "post1"
	lastUsername := "username2"
	limit := 4
	ginContext.Request, _ = http.NewRequest("GET", "/postSuperlikes", nil)
	ginContext.Params = []gin.Param{{Key: "postId", Value: postId}}
	u := url.Values{}
	u.Add("lastUsername", lastUsername)
	u.Add("limit", strconv.Itoa(limit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"postSuperlikes":[	
			{
				"username":  "username3",
				"name": 	 "fullname3"
			},		
			{
				"username":  "username4",
				"name": 	 "fullname4"
			},		
			{
				"username":  "username5",
				"name": 	 "fullname5"
			},		
			{
				"username":  "username6",
				"name": 	 "fullname6"
			}
			],
			"lastUsername":"username6"
		}
	}`

	controller.GetPostSuperlikesMetadata(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}
func TestDeletePostLike_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	existingPost := &database.PostMetadata{
		PostId:   "post123",
		Username: "username1",
		Type:     "TEXT",
		Likes:    1,
	}
	integration_test_arrange.AddPostToDatabase(t, db, existingPost)
	data := &reaction_handler.UserUnlikedPostEvent{
		Username: "user123",
		PostId:   existingPost.PostId,
	}
	event, _ := test_common.SerializeData(data)
	expectedLike := &model.PostLike{
		User: &model.UserMetadata{
			Username: data.Username,
		},
		PostId: data.PostId,
	}

	userUnlikedPostEventHandler.Handle(event)

	integration_test_assert.AssertPostLikeDoesNotExists(t, db, expectedLike)
	integration_test_assert.AssertPostLikesDecreased(t, db, existingPost.PostId)
}

func TestDeletePostSuperlike_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	existingPost := &database.PostMetadata{
		PostId:     "post123",
		Username:   "username1",
		Type:       "TEXT",
		Superlikes: 1,
	}
	integration_test_arrange.AddPostToDatabase(t, db, existingPost)
	data := &reaction_handler.UserUnsuperlikedPostEvent{
		Username: "user123",
		PostId:   existingPost.PostId,
	}
	event, _ := test_common.SerializeData(data)
	expectedSuperlike := &model.PostSuperlike{
		User: &model.UserMetadata{
			Username: data.Username,
		},
		PostId: data.PostId,
	}

	userUnsuperlikedPostEventHandler.Handle(event)

	integration_test_assert.AssertPostSuperlikeDoesNotExists(t, db, expectedSuperlike)
	integration_test_assert.AssertPostSuperlikesDecreased(t, db, existingPost.PostId)
}

func populatePostLikesDb(t *testing.T) {
	existingPostLikes := []*database.PostLikeMetadata{
		{
			PostId:   "post1",
			Username: "username1",
			Name:     "fullname1",
		},
		{
			PostId:   "post2",
			Username: "username1",
			Name:     "fullname1",
		},
		{
			PostId:   "post3",
			Username: "username1",
			Name:     "fullname1",
		},
		{
			PostId:   "post1",
			Username: "username2",
			Name:     "fullname2",
		},
		{
			PostId:   "post1",
			Username: "username3",
			Name:     "fullname3",
		},
		{
			PostId:   "post2",
			Username: "username1",
			Name:     "fullname1",
		},
		{
			PostId:   "post2",
			Username: "username1",
			Name:     "fullname1",
		},
		{
			PostId:   "post1",
			Username: "username4",
			Name:     "fullname4",
		},
		{
			PostId:   "post1",
			Username: "username5",
			Name:     "fullname5",
		},
		{
			PostId:   "post3",
			Username: "username1",
			Name:     "fullname1",
		},
		{
			PostId:   "post1",
			Username: "username6",
			Name:     "fullname6",
		},
		{
			PostId:   "post1",
			Username: "username7",
			Name:     "fullname7",
		},
	}

	for _, existingPostLike := range existingPostLikes {
		integration_test_arrange.AddPostLikeToDatabase(t, db, existingPostLike)
	}
}

func populatePostSuperlikesDb(t *testing.T) {
	existingPostSuperlikes := []*database.PostSuperlikeMetadata{
		{
			PostId:   "post1",
			Username: "username1",
			Name:     "fullname1",
		},
		{
			PostId:   "post2",
			Username: "username1",
			Name:     "fullname1",
		},
		{
			PostId:   "post3",
			Username: "username1",
			Name:     "fullname1",
		},
		{
			PostId:   "post1",
			Username: "username2",
			Name:     "fullname2",
		},
		{
			PostId:   "post1",
			Username: "username3",
			Name:     "fullname3",
		},
		{
			PostId:   "post2",
			Username: "username1",
			Name:     "fullname1",
		},
		{
			PostId:   "post2",
			Username: "username1",
			Name:     "fullname1",
		},
		{
			PostId:   "post1",
			Username: "username4",
			Name:     "fullname4",
		},
		{
			PostId:   "post1",
			Username: "username5",
			Name:     "fullname5",
		},
		{
			PostId:   "post3",
			Username: "username1",
			Name:     "fullname1",
		},
		{
			PostId:   "post1",
			Username: "username6",
			Name:     "fullname6",
		},
		{
			PostId:   "post1",
			Username: "username7",
			Name:     "fullname7",
		},
	}

	for _, existingPostSuperlike := range existingPostSuperlikes {
		integration_test_arrange.AddPostSuperlikeToDatabase(t, db, existingPostSuperlike)
	}
}
