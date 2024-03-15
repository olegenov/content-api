package responseControllers

import (
	"contentApi/dto/responses"
	"contentApi/models"
)

func GetPostResponse(post models.Post) responses.PostResponse {
	postResponse := responses.PostResponse{
		Title:       post.Title,
		PublishDate: post.PublishDate,
		Deadline:    post.Deadline,
		Content:     post.Content,
	}

	for _, tag := range post.Tags {
		postResponse.Tags = append(postResponse.Tags, GetTagResponse(tag))
	}

	return postResponse
}
