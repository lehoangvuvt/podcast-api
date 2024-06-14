package models

import (
	"database/sql"
)

type PostLike struct {
	ID        int            `json:"id"`
	PostID    int            `json:"post_id"`
	UserID    int            `json:"user_id"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt sql.NullString `json:"-"`
}

type PostLikeModel struct {
	DB *sql.DB
}

func (m *PostLikeModel) Delete(userId int, postId int) error {
	_, err := m.DB.Exec("DELETE FROM posts_likes WHERE post_id=$1 AND user_id=$2", postId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostLikeModel) Insert(userId int, postId int) error {
	_, err := m.DB.Exec("INSERT INTO posts_likes (post_id, user_id) VALUES ($1, $2)", postId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostLikeModel) GetPostLikesByPostId(postId int) ([]PostLike, error) {
	rows, err := m.DB.Query("SELECT * FROM posts_likes WHERE post_id=$1", postId)
	if err != nil {
		return make([]PostLike, 0), err
	}
	var postLikes []PostLike
	for rows.Next() {
		var postLike PostLike
		err = rows.Scan(&postLike.ID, &postLike.PostID, &postLike.UserID, &postLike.CreatedAt, &postLike.UpdatedAt)
		if err == nil {
			postLikes = append(postLikes, postLike)
		}
	}
	return postLikes, nil
}
