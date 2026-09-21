package service

import (
	"context"
	"errors"

	"minisocial/internal/domain"
	"minisocial/internal/repository"

	"github.com/google/uuid"
)

type SocialService struct {
	repo *repository.PostgresRepository
}

func NewSocialService(repo *repository.PostgresRepository) *SocialService {
	return &SocialService{repo: repo}
}

func (s *SocialService) Register(ctx context.Context, username string) (domain.User, error) {
	if len(username) < 3 {
		return domain.User{}, errors.New("username must be at least 3 characters")
	}
	return s.repo.CreateUser(ctx, username)
}

func (s *SocialService) CreatePost(ctx context.Context, authorID uuid.UUID, content string) (domain.Post, error) {
	if len(content) == 0 {
		return domain.Post{}, errors.New("post content cannot be empty")
	}
	return s.repo.CreatePost(ctx, authorID, content)
}

func (s *SocialService) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}
	return s.repo.FollowUser(ctx, followerID, followingID)
}

func (s *SocialService) Like(ctx context.Context, userID, postID uuid.UUID) error {
	return s.repo.LikePost(ctx, userID, postID)
}

func (s *SocialService) GetFeed(ctx context.Context, userID uuid.UUID) ([]domain.Post, error) {
	return s.repo.GetFeed(ctx, userID)
}
