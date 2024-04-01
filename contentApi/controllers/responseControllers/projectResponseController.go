package responseControllers

import (
	"contentApi/dto/responses"
	"contentApi/models"
)

func GetProjectResponse(project models.Project) responses.ProjectResponse {
	projectResponse := responses.ProjectResponse{
		ID:   project.ID,
		Name: project.Name,
	}

	projectResponse.Creator = GetUserResponse(project.Creator)
	projectResponse.Team = GetProjectTeamResponse(project.Team)
	projectResponse.Posts = []responses.PostResponse{}

	for _, post := range project.Posts {
		projectResponse.Posts = append(projectResponse.Posts, GetPostResponse(post))
	}

	return projectResponse
}
