package requests

import "time"

type EditPostRequest struct {
	Title       string    `json:"title"`
	Assign      string    `json:"assign"`
	PublishDate time.Time `json:"publishing"`
	Deadline    time.Time `json:"deadline"`
	Content     string    `json:"content"`
}
