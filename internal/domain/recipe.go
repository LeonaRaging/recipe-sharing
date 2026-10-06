package domain

import "time"

type Recipe struct {
	ID          int64
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int64  `json:"userid"`
	PublishedAt time.Time
}
