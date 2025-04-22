package kafka

func getTopics() []string {
	return []string{
		"UserWasRegisteredEvent",
		"UserProfileUpdatedEvent",
		"PostWasCreatedEvent",
		"PostsWereDeletedEvent",
		"UserAFollowedUserBEvent",
		"UserAUnfollowedUserBEvent",
		"CommentWasCreatedEvent",
		"CommentWasUpdatedEvent",
		"CommentWasDeletedEvent",
		"UserLikedPostEvent",
		"UserUnlikedPostEvent",
		"UserSuperlikedPostEvent",
		"UserUnsuperlikedPostEvent",
	}
}
