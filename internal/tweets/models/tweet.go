package models

import (
	"time"
	"twitter-api/internal/tweets/dto"
)

type Tweet struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index;not null"`
	Content   string    `gorm:"size:280;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (t Tweet) ToDTO() *dto.Tweet {
	result := &dto.Tweet{
		ID:      t.ID,
		UserID:  t.UserID,
		Content: t.Content,
	}
	return result
}
