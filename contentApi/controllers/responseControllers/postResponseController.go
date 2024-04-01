package responseControllers

import (
	"contentApi/dto/responses"
	"contentApi/models"
)

func GetPostResponse(post models.Post) responses.PostResponse {
	postResponse := responses.PostResponse{
		ID:          post.ID,
		Title:       post.Title,
		PublishDate: post.PublishDate,
		Deadline:    post.Deadline,
		Content:     post.Content,
	}

	postResponse.Assign = GetUserResponse(post.Assign)

	for _, tag := range post.Tags {
		postResponse.Tags = append(postResponse.Tags, GetTagResponse(tag))
	}

	return postResponse
}

func GetSinglePostResponse(post models.Post) responses.SinglePostResponse {
	postResponse := responses.SinglePostResponse{
		ProjectName: post.Project.Name,
		Info:        GetPostResponse(post),
	}

	return postResponse
}
