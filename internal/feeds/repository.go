package feeds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"twitter-api/internal/tweets/models"
)

type Repository interface {
	GetTweetIDFromQueue(ctx context.Context) (uint, error)
	GetTweet(ctx context.Context, id uint) (*models.Tweet, error)
	GetFollowers(ctx context.Context, userID uint) ([]uint, error)
	PushTweetToFollowersFeed(ctx context.Context, tweetID uint, followers []uint) error
}

type repository struct {
	redis *redis.Client
}

func NewRepository(redis *redis.Client) Repository {
	return &repository{redis: redis}
}

func (r *repository) GetTweetIDFromQueue(ctx context.Context) (uint, error) {
	result, err := r.getTweetIDFromQueue(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}
	if result == "" {
		return 0, nil
	}
	tweetID, err := stouint(result)
	if err != nil {
		return 0, fmt.Errorf("[Repository] Error unexpected data type: %T", result)
	}
	return tweetID, nil
}

func (r *repository) GetTweet(ctx context.Context, id uint) (*models.Tweet, error) {
	tweetKey := fmt.Sprintf("tweets:%d", id)

	tweetData, err := r.redis.Get(ctx, tweetKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("[Repository] Error tweet ID %d does not exist", id)
		}
		return nil, fmt.Errorf("[Repository] Error obtaining tweet: %w", err)
	}

	var tweetJson struct {
		UserID  string `json:"userId"`
		Content string `json:"content"`
	}
	err = json.Unmarshal([]byte(tweetData), &tweetJson)
	if err != nil {
		return nil, fmt.Errorf("[Repository] Error deserializing tweet JSON: %w", err)
	}

	userID, _ := stouint(tweetJson.UserID)
	tweet := models.Tweet{
		ID:      id,
		UserID:  userID,
		Content: tweetJson.Content,
	}

	return &tweet, nil
}

func (r *repository) GetFollowers(ctx context.Context, userID uint) ([]uint, error) {
	key := fmt.Sprintf("followers:%d", userID)

	followers, err := r.redis.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("[Repository] Error obtaining followers: %w", err)
	}
	followersIDs := make([]uint, len(followers))
	for i, follower := range followers {
		followersIDs[i], err = stouint(follower)
		if err != nil {
			return nil, fmt.Errorf("[Repository] Error casting follower ids: %w", err)
		}
	}
	return followersIDs, nil
}

func (r *repository) PushTweetToFollowersFeed(ctx context.Context, tweetID uint, followers []uint) error {
	pipe := r.redis.Pipeline()
	for _, followerID := range followers {
		feedKey := fmt.Sprintf("feed:%d", followerID)
		pipe.LPush(ctx, feedKey, tweetID)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("[Repository] Error updating feeds: %w", err)
	}
	return nil
}

func (r *repository) getTweetIDFromQueue(ctx context.Context) (string, error) {
	tweetID, err := r.redis.RPop(ctx, "tweet_queue").Result()
	if err != nil {
		return "", err
	}
	exists, err := r.redis.SIsMember(ctx, "processed_tweets", tweetID).Result()
	if err != nil {
		return "", err
	}
	if !exists {
		r.redis.SAdd(ctx, "processed_tweets", tweetID).Result()
		return tweetID, nil
	}
	return "", nil
}

func stouint(s string) (uint, error) {
	u64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(u64), nil
}
