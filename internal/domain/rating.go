package domain

import "time"

type Rating struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	RecipeID  int64     `json:"recipe_id" gorm:"not null;index"`
	UserID    int64     `json:"user_id" gorm:"not null"`
	Score     int       `json:"score" gorm:"not null"`
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
