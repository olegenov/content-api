package responseControllers

import (
	"contentApi/dto/responses"
	"contentApi/models"
)

func GetTagResponse(tag models.Tag) responses.TagResponse {
	tagResponse := responses.TagResponse{
		Name:  tag.Name,
		Color: tag.Color,
	}

	return tagResponse
}
