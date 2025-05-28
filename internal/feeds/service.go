package feeds

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Service interface {
	CreateFeeds()
}

type feedsService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &feedsService{
		repo: repo,
	}
}

func (s *feedsService) CreateFeeds() {
	ctx := context.Background()
	numWorkers := 5

	tweetIDs := make(chan uint)

	for i := 0; i < numWorkers; i++ {
		go s.worker(tweetIDs)
	}

	for {
		tweetID, err := s.repo.GetTweetIDFromQueue(ctx)
		if err != nil || tweetID == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		tweetIDs <- tweetID
	}
}

func (s *feedsService) worker(tweetIDs <-chan uint) {
	for tweetID := range tweetIDs {
		if err := s.sendTweetToFollowers(tweetID); err != nil {
			log.Printf("[FeedsService] Error processing tweet ID=%d: %v", tweetID, err)
		}
	}
}

func (s *feedsService) sendTweetToFollowers(tweetID uint) error {
	ctx := context.Background()
	tweet, err := s.repo.GetTweet(ctx, tweetID)
	if err != nil {
		return fmt.Errorf("[FeedsService] Error obtaining tweet: %w", err)
	}

	followers, err := s.repo.GetFollowers(ctx, tweet.UserID)
	if err != nil {
		return fmt.Errorf("[FeedsService] Error obtaining followers: %w", err)
	}

	err = s.repo.PushTweetToFollowersFeed(ctx, tweetID, followers)
	if err != nil {
		return err
	}
	return nil
}
