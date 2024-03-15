package responses

import "time"

type PostResponse struct {
	Title       string `json:"title"`
	PublishDate time.Time
	Deadline    time.Time
	Tags        []TagResponse `json:"tags"`
	Content     string        `json:"content"`
}
