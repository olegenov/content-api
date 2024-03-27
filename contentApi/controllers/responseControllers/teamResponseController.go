package responseControllers

import (
	"contentApi/dto/responses"
	"contentApi/models"
)

func GetTeamResponse(team models.Team) responses.TeamResponse {
	teamResponse := responses.TeamResponse{
		ID:   team.ID,
		Name: team.Name,
	}

	teamResponse.Creator = GetUserResponse(team.Creator)

	for _, user := range team.Users {
		teamResponse.Users = append(teamResponse.Users, GetUserResponse(user))
	}

	return teamResponse
}

func GetProjectTeamResponse(team models.Team) responses.ProjectTeamResponse {
	teamResponse := responses.ProjectTeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		CreatorID: team.CreatorID,
	}

	return teamResponse
}
