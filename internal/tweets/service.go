package tweets

import (
	"context"
	"twitter-api/internal/tweets/dto"
	"twitter-api/internal/tweets/models"
)

type Service interface {
	Create(ctx context.Context, tweetDTO *dto.CreateTweet) (*dto.Tweet, error)
	Delete(ctx context.Context, id uint) error
}

type tweetsService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &tweetsService{
		repo: repo,
	}
}

func (s *tweetsService) Create(ctx context.Context, tweetDTO *dto.CreateTweet) (*dto.Tweet, error) {
	tweet := models.Tweet{
		UserID:  tweetDTO.UserID,
		Content: tweetDTO.Content,
	}
	newTweet, err := s.repo.Create(ctx, &tweet)
	if err != nil {
		return nil, err
	}

	return newTweet.ToDTO(), nil
}

func (s *tweetsService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
