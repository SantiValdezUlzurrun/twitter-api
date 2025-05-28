package models

import (
	"time"
	"twitter-api/internal/users/dto"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"not null;unique"`
	Followers int       `gorm:"default:0"`
	Following int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (u User) ToDTO() *dto.User {
	result := &dto.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Followers: u.Followers,
		Following: u.Following,
	}
	return result
}
