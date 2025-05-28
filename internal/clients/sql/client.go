package sql

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	tweetsModels "twitter-api/internal/tweets/models"
	usersModels "twitter-api/internal/users/models"
)

type Client interface {
	CreateUser(ctx context.Context, user *usersModels.User) error
	Follow(ctx context.Context, userID, followerID uint) error
	Unfollow(ctx context.Context, userID, followerID uint) error
	CreateTweet(ctx context.Context, tweet *tweetsModels.Tweet) error
	GetTweet(ctx context.Context, id uint) (*tweetsModels.Tweet, error)
	DeleteTweet(ctx context.Context, tweet *tweetsModels.Tweet) error
}

type sql struct {
	db *gorm.DB
}

func NewSqlClient(db *gorm.DB) Client {
	return &sql{db: db}
}

func (c *sql) CreateUser(ctx context.Context, user *usersModels.User) error {
	if err := c.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(context.DeadlineExceeded, ctx.Err()) {
			return fmt.Errorf("[SQL Client] Deadline exceeded")
		}
		return fmt.Errorf("[SQL Client] Error creating user: %w", err)
	}
	return nil
}

func (c *sql) Follow(ctx context.Context, userID, followerID uint) error {
	err := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user, follower usersModels.User

		if err := tx.First(&user, "id = ?", userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("[SQL Client] User not found: %w", err)
			}
			return fmt.Errorf("[SQL Client] Error: %w", err)
		}

		if err := tx.First(&follower, "id = ?", followerID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("[SQL Client] User not found: %w", err)
			}
			return fmt.Errorf("[SQL Client] Error: %w", err)
		}

		var count int64
		if err := tx.Model(&usersModels.Follower{}).
			Where("user_id = ? AND follower_id = ?", userID, followerID).
			Count(&count).Error; err != nil {
			return fmt.Errorf("[SQL Client] Error checking if user already follows: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("[SQL Client] User already follows")
		}

		followerRecord := &usersModels.Follower{
			UserID:     userID,
			FollowerID: followerID,
		}
		if err := tx.Create(followerRecord).Error; err != nil {
			return fmt.Errorf("[SQL Client] Error following user: %w", err)
		}

		if err := tx.Model(&user).UpdateColumn("followers", gorm.Expr("followers + ?", 1)).Error; err != nil {
			return fmt.Errorf("[SQL Client] Error updating followers: %w", err)
		}

		if err := tx.Model(&follower).UpdateColumn("following", gorm.Expr("following + ?", 1)).Error; err != nil {
			return fmt.Errorf("[SQL Client] Error updating following:  %w", err)
		}

		return nil
	})
	return err
}

func (c *sql) Unfollow(ctx context.Context, userID, followerID uint) error {
	err := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var followerRecord usersModels.Follower
		if err := tx.Where("user_id = ? AND follower_id = ?", userID, followerID).First(&followerRecord).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("[SQL Client] User does not follow: %w", err)
			}
			return fmt.Errorf("[SQL Client] Error: %w", err)
		}

		if err := tx.Delete(&followerRecord).Error; err != nil {
			return fmt.Errorf("[SQL Client] Error deleting follow: %w", err)
		}

		if err := tx.Model(&usersModels.User{}).Where("id = ?", userID).UpdateColumn("followers", gorm.Expr("followers - ?", 1)).Error; err != nil {
			return fmt.Errorf("[SQL Client] Error updating followers: %w", err)
		}

		if err := tx.Model(&usersModels.User{}).Where("id = ?", followerID).UpdateColumn("following", gorm.Expr("following - ?", 1)).Error; err != nil {
			return fmt.Errorf("[SQL Client] Error updating following:  %w", err)
		}

		return nil
	})
	return err
}

func (c *sql) CreateTweet(ctx context.Context, tweet *tweetsModels.Tweet) error {
	err := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(tweet).Error; err != nil {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return fmt.Errorf("[SQL Client] Deadline exceeded")
			}
			return fmt.Errorf("[SQL Client] Error creating tweet: %w", err)
		}
		return nil
	})
	return err
}

func (c *sql) GetTweet(ctx context.Context, id uint) (*tweetsModels.Tweet, error) {
	tweet := &tweetsModels.Tweet{}
	if err := c.db.WithContext(ctx).First(tweet, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("[SQL Client] Error tweet not found")
		}
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("[SQL Client] Deadline exceeded")
		}
		return nil, fmt.Errorf("[SQL Client] Error obtaining tweet: %w", err)
	}
	return tweet, nil
}

func (c *sql) DeleteTweet(ctx context.Context, tweet *tweetsModels.Tweet) error {
	err := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(tweet).Error; err != nil {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return fmt.Errorf("[SQL Client] Deadline exceeded")
			}
			return fmt.Errorf("[SQL Client] Error deleting tweet:: %w", err)
		}
		return nil
	})
	return err
}
