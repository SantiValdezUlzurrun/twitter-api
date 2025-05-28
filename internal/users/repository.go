package users

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"twitter-api/internal/clients/sql"
	"twitter-api/internal/users/models"
)

type Repository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Follow(ctx context.Context, id, followerID uint) error
	Unfollow(ctx context.Context, id, followerID uint) error
}

type repository struct {
	sqlClient sql.Client
	redis     *redis.Client
}

func NewRepository(sqlClient sql.Client, redis *redis.Client) Repository {
	return &repository{sqlClient: sqlClient, redis: redis}
}

func (r *repository) Create(ctx context.Context, userModel *models.User) (*models.User, error) {

	err := r.sqlClient.CreateUser(ctx, userModel)
	if err != nil {
		return nil, fmt.Errorf("[Repository] error creating user: %w", err)
	}

	userData, err := json.Marshal(struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{
		ID:   fmt.Sprintf("%d", userModel.ID),
		Name: userModel.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("[Repository] error serializing user: %w", err)
	}

	key := fmt.Sprintf("users:%d", userModel.ID)
	err = r.redis.Set(ctx, key, userData, 0).Err()
	if err != nil {
		return nil, err
	}

	return userModel, nil
}

func (r *repository) Follow(ctx context.Context, userID, followerID uint) error {

	if err := r.sqlClient.Follow(ctx, userID, followerID); err != nil {
		return err
	}

	pipe := r.redis.TxPipeline()
	pipe.SAdd(ctx, fmt.Sprintf("following:%d", followerID), userID)
	pipe.SAdd(ctx, fmt.Sprintf("followers:%d", userID), followerID)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("[Repository] redis error updating followers: %w", err)
	}

	return nil
}

func (r *repository) Unfollow(ctx context.Context, userID, followerID uint) error {

	if err := r.sqlClient.Unfollow(ctx, userID, followerID); err != nil {
		return err
	}

	pipe := r.redis.TxPipeline()
	pipe.SRem(ctx, fmt.Sprintf("following:%d", followerID), userID)
	pipe.SRem(ctx, fmt.Sprintf("followers:%d", userID), followerID)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("[Repository] redis error updating followers: %w", err)
	}
	return nil
}
