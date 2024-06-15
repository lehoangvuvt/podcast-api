package models

import (
	"database/sql"
	"fmt"
)

type PostComment struct {
	ID               int            `json:"id"`
	PostID           int            `json:"post_id"`
	UserID           int            `json:"user_id"`
	ReplyToCommentId sql.NullInt32  `json:"reply_to_comment_id"`
	Content          string         `json:"content"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        sql.NullString `json:"-"`
}

type PostCommentWithUser struct {
	ID               int           `json:"id"`
	PostID           int           `json:"post_id"`
	UserID           int           `json:"user_id"`
	ReplyToCommentId sql.NullInt32 `json:"reply_to_comment_id"`
	UserDetails      UserDetails   `json:"user_details"`
	Content          string        `json:"content"`
	CreatedAt        string        `json:"created_at"`
}

type CreateCommentInput struct {
	ReplyToCommentId int    `json:"reply_to_comment_id"`
	Content          string `json:"content"`
}

type PostCommentModel struct {
	DB *sql.DB
}

func (m *PostCommentModel) Insert(userId int, postId int, createCommentInput *CreateCommentInput) error {
	var stmt string
	if createCommentInput.ReplyToCommentId != -1 {
		stmt = fmt.Sprintf(`INSERT INTO 
							posts_comments (post_id, user_id, reply_to_comment_id, content) 
							VALUES (%v, %v, %v, '%v')`, postId, userId, createCommentInput.ReplyToCommentId, createCommentInput.Content)
	} else {
		stmt = fmt.Sprintf(`INSERT INTO 
							posts_comments (post_id, user_id, content) 
							VALUES (%v, %v, '%v')`, postId, userId, createCommentInput.Content)
	}
	_, err := m.DB.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostCommentModel) GetCommentsByPostId(postId int) ([]PostCommentWithUser, error) {
	var postComments []PostCommentWithUser
	rows, err := m.DB.Query(`SELECT id, post_id, user_id, reply_to_comment_id, content, created_at 
							FROM posts_comments WHERE post_id = $1 ORDER BY created_at DESC`, postId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var postComment PostCommentWithUser
		err = rows.Scan(&postComment.ID, &postComment.PostID, &postComment.UserID, &postComment.ReplyToCommentId,
			&postComment.Content, &postComment.CreatedAt)
		if err == nil {
			var userDetails UserDetails
			row := m.DB.QueryRow(`SELECT id, username, email from users WHERE id=$1`, postComment.UserID)
			err = row.Scan(&userDetails.ID, &userDetails.Username, &userDetails.Email)
			if err == nil {
				postComment.UserDetails = userDetails
			}
		}
		postComments = append(postComments, postComment)
	}
	return postComments, nil
}
