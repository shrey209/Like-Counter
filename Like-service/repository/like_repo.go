package repository

import (
	"github.com/gocql/gocql"
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
