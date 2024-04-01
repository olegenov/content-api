package responses

import (
	"time"
)

type PostResponse struct {
	ID          uint          `json:"id"`
	Title       string        `json:"title"`
	Assign      UserResponse  `json:"assign"`
	PublishDate time.Time     `json:"publishing"`
	Deadline    time.Time     `json:"deadline"`
	Tags        []TagResponse `json:"tags"`
	Content     string        `json:"content"`
}

type SinglePostResponse struct {
	ProjectName string       `json:"project_name"`
	Info        PostResponse `json:"info"`
}
