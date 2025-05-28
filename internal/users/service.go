package users

import (
	"context"
	"fmt"
	"twitter-api/internal/users/dto"
	"twitter-api/internal/users/models"
)

type Service interface {
	Create(ctx context.Context, user *dto.CreateUser) (*dto.User, error)
	Follow(ctx context.Context, id, followerID uint) error
	Unfollow(ctx context.Context, id, followerID uint) error
}

type userService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Create(ctx context.Context, userDTO *dto.CreateUser) (*dto.User, error) {
	user := models.User{
		Name:  userDTO.Name,
		Email: userDTO.Email,
	}
	newUser, err := s.repo.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	return newUser.ToDTO(), nil
}

func (s *userService) Follow(ctx context.Context, id, followerID uint) error {
	if id == followerID {
		return fmt.Errorf("[Service] User cant follow itself")
	}
	return s.repo.Follow(ctx, id, followerID)
}

func (s *userService) Unfollow(ctx context.Context, id, followerID uint) error {
	if id == followerID {
		return fmt.Errorf("[Service] User cant follow itself")
	}

	return s.repo.Unfollow(ctx, id, followerID)
}
