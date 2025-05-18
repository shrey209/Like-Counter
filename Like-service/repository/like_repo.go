package repository

import (
	"github.com/gocql/gocql"
	"github.com/shrey209/Like-Service/model"
)

type LikeRepository struct {
	Session *gocql.Session
}

func NewLikeRepository(session *gocql.Session) *LikeRepository {
	return &LikeRepository{Session: session}
}

func (r *LikeRepository) IncrementLike(postID string) error {
	query := `UPDATE post_likes SET like_count = like_count + 1 WHERE post_id = ?`
	return r.Session.Query(query, postID).Exec()
}

func (r *LikeRepository) GetLikeCount(postID string) (int64, error) {
	var count int64
	query := `SELECT like_count FROM post_likes WHERE post_id = ?`
	if err := r.Session.Query(query, postID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *LikeRepository) BatchIncrementLikes(likes []model.PostLike) error {
	batch := r.Session.NewBatch(gocql.UnloggedBatch)

	for _, like := range likes {
		batch.Query(`UPDATE post_likes SET like_count = like_count + ? WHERE post_id = ?`, like.LikeCount, like.PostID)
	}

	return r.Session.ExecuteBatch(batch)
}

func (r *LikeRepository) CreatePostLikeEntry(postID string) error {
	query := `INSERT INTO post_likes (post_id, like_count) VALUES (?, 0)`
	return r.Session.Query(query, postID).Exec()
}

func (r *LikeRepository) InitPostLike(postID string) error {
	query := `UPDATE post_likes SET like_count = like_count + 0 WHERE post_id = ?`
	return r.Session.Query(query, postID).Exec()
}
