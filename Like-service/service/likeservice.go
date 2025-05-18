package service

import (
	"github.com/shrey209/Like-Service/model"
	"github.com/shrey209/Like-Service/repository"
)

type LikeService struct {
	Repo *repository.LikeRepository
}

func NewLikeService(repo *repository.LikeRepository) *LikeService {
	return &LikeService{Repo: repo}
}

func (s *LikeService) LikePost(postID string) error {
	return s.Repo.IncrementLike(postID)
}

func (s *LikeService) GetPostLikes(postID string) (int64, error) {
	return s.Repo.GetLikeCount(postID)
}

func (s *LikeService) BatchLikePosts(likes []model.PostLike) error {
	return s.Repo.BatchIncrementLikes(likes)
}
func (s *LikeService) InitPostLike(postID string) error {
	return s.Repo.InitPostLike(postID)
}
