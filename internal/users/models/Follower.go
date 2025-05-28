package models

type Follower struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint `gorm:"index;not null"`
	FollowerID uint `gorm:"index;not null"`
}
