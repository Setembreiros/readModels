package post

import database "readmodels/internal/db"

type PostRepository database.Database

func (r PostRepository) AddNewPostMetadata(data *PostMetadata) error {
	return r.Client.InsertData("PostMetadata", data)
}

func (r PostRepository) GetPostMetadatasByUser(username string) ([]*PostMetadata, error) {
	data, err := r.Client.GetPostsByIndexUser(username)
	if err != nil {
		return []*PostMetadata{}, err
	}

	var posts []*PostMetadata
	for _, post := range data {
		posts = append(posts, mapToDomain(post))
	}

	return posts, nil
}

func (r PostRepository) RemovePostMetadata(postIds []string) error {
	anySlice := make([]any, len(postIds))
	for i, v := range postIds {
		anySlice[i] = v
	}

	err := r.Client.RemoveMultipleData("Posts", anySlice)

	return err
}

func mapToDomain(data *database.PostMetadata) *PostMetadata {
	return &PostMetadata{
		PostId:      data.PostId,
		Username:    data.Username,
		Type:        data.Type,
		FileType:    data.FileType,
		Title:       data.Title,
		Description: data.Description,
		CreatedAt:   data.CreatedAt,
		LastUpdated: data.LastUpdated,
	}
}
