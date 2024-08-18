package post

import database "readmodels/internal/db"

type PostRepository database.Database

func (r PostRepository) AddNewPostMetadata(data *PostMetadata) error {
	return r.Client.InsertData("PostMetadata", data)
}
