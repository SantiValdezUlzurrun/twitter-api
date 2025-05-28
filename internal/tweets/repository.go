package tweets

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"twitter-api/internal/clients/sql"
	"twitter-api/internal/tweets/models"
)

type Repository interface {
	Create(ctx context.Context, tweet *models.Tweet) (*models.Tweet, error)
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	sqlClient sql.Client
	redis     *redis.Client
}

func NewRepository(sqlClient sql.Client, redis *redis.Client) Repository {
	return &repository{sqlClient: sqlClient, redis: redis}
}

func (r *repository) Create(ctx context.Context, tweet *models.Tweet) (*models.Tweet, error) {

	if err := r.sqlClient.CreateTweet(ctx, tweet); err != nil {
		return nil, err
	}
	if err := r.cacheTweet(ctx, tweet); err != nil {
		return nil, err
	}

	return tweet, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	tweet, err := r.sqlClient.GetTweet(ctx, id)

	err = r.sqlClient.DeleteTweet(ctx, tweet)
	if err != nil {
		return err
	}

	err = r.redisDeleteTweet(ctx, tweet)
	if err != nil {
		return fmt.Errorf("[Repository] Error deleting tweet: %w", err)
	}

	return nil
}

func (r *repository) redisDeleteTweet(ctx context.Context, tweet *models.Tweet) error {
	tweetKey := fmt.Sprintf("tweet:%s", tweet.ID)
	timelineKey := fmt.Sprintf("timeline:%s", tweet.UserID)

	pipe := r.redis.Pipeline()
	pipe.Del(ctx, tweetKey)
	pipe.ZRem(ctx, timelineKey, tweet.ID)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *repository) cacheTweet(ctx context.Context, tweet *models.Tweet) error {
	tweetKey := fmt.Sprintf("tweets:%s", tweet.ID)

	tweetData, err := newTweet(tweet)
	if err != nil {
		return fmt.Errorf("[Repository] error serializing tweet: %w", err)
	}

	pipe := r.redis.Pipeline()
	pipe.Set(ctx, tweetKey, tweetData, 0)

	if err := pipe.LPush(ctx, "tweet_queue", tweet.ID).Err(); err != nil {
		return fmt.Errorf("[Repository] error pushing tweet: %w", err)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("[Repository] Redis error: %w", err)
	}

	return nil
}

func newTweet(tw *models.Tweet) ([]byte, error) {
	jsonData, err := json.Marshal(struct {
		UserID  string `json:"userId"`
		Content string `json:"content"`
	}{
		UserID:  fmt.Sprintf("%d", tw.UserID),
		Content: tw.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("[Repository] error serializing tweet: %w", err)
	}
	return jsonData, nil
}
