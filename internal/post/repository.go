package post

import database "readmodels/internal/db"

type PostRepository database.Database

func (r PostRepository) AddNewPostMetadata(data *PostMetadata) error {
	return r.Client.InsertData("PostMetadata", data)
}

type PostMetadataKey struct {
	PostId string
}

func (r PostRepository) GetPostMetadatasByUser(username string, lastPostId, lastPostCreatedAt string, limit int) ([]*PostMetadata, string, string, error) {
	data, lastPostId, lastPostCreatedAt, err := r.Client.GetPostsByIndexUser(username, lastPostId, lastPostCreatedAt, limit)
	if err != nil {
		return []*PostMetadata{}, "", "", err
	}

	var posts []*PostMetadata
	for _, post := range data {
		posts = append(posts, mapToDomain(post))
	}

	return posts, lastPostId, lastPostCreatedAt, nil
}

func (r PostRepository) RemovePostMetadata(postIds []string) error {
	postKeys := make([]any, len(postIds))
	for i, v := range postIds {
		postKeys[i] = &PostMetadataKey{
			PostId: v,
		}
	}

	err := r.Client.RemoveMultipleData("PostMetadata", postKeys)

	return err
}

func mapToDomain(data *database.PostMetadata) *PostMetadata {
	return &PostMetadata{
		PostId:      data.PostId,
		Username:    data.Username,
		Type:        data.Type,
		Title:       data.Title,
		Description: data.Description,
		CreatedAt:   data.CreatedAt,
		LastUpdated: data.LastUpdated,
	}
}
