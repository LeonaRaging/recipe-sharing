package domain

type Favorite struct {
	UserID   int64 `gorm:"primaryKey;autoIncrement:false"`
	RecipeID int64 `gorm:"primaryKey;autoIncrement:false"`
}
