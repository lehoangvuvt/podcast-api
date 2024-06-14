package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	queryHelpers "vulh/soundcommunity/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int            `json:"id"`
	Username       string         `json:"username"`
	HashedPassword string         `json:"-"`
	Email          string         `json:"email"`
	CreatedAt      string         `json:"-"`
	UpdatedAt      sql.NullString `json:"-"`
}

type UserDetails struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CreateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserFavouriteEpisode struct {
	ID        int            `json:"id"`
	UserId    int            `json:"user_id"`
	EpisodeId int            `json:"episode_id"`
	CreatedAt string         `json:"-"`
	UpdatedAt sql.NullString `json:"-"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(createUserInput *CreateUserInput) error {
	stmt := `INSERT INTO users(username, hashed_password, email) VALUES($1, $2, $3)`
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(createUserInput.Password), 10)
	_, err := m.DB.Exec(stmt, createUserInput.Username, hashedPassword, createUserInput.Email)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) GetUserById(id int) (*UserDetails, error) {
	userDetails := &UserDetails{}
	queryBuilder := &queryHelpers.QueryBuilder{DB: m.DB}
	row := queryBuilder.
		Select("id", "username", "email").
		FromTable("users").
		WhereColumn("id").
		Equal(fmt.Sprintf("%v", id)).
		GetOne()
	err := row.Scan(&userDetails.ID, &userDetails.Username, &userDetails.Email)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	return userDetails, nil
}

func (m *UserModel) Login(loginInput *LoginInput) (*UserDetails, error) {
	user := &User{}
	queryBuilder := &queryHelpers.QueryBuilder{DB: m.DB}
	row := queryBuilder.
		Select("id", "username", "hashed_password", "email").
		FromTable("users").
		WhereColumn("username").
		Equal(loginInput.Username).
		GetOne()
	err := row.Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Email)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginInput.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	userDetails, err := m.GetUserById(user.ID)
	if err != nil {
		return nil, errors.New("cannot get user details")
	}
	return userDetails, nil
}

func (m *UserModel) CreateUserFavouriteEpisode(userId int, episodeId int) (*UserFavouriteEpisode, error) {
	userFavouriteSong := &UserFavouriteEpisode{}
	var lastInsertId int
	row := m.DB.QueryRow(`INSERT INTO user_favourite_episodes (user_id, episode_id) 
								VALUES ($1, $2) RETURNING id`, userId, episodeId)
	err := row.Scan(&lastInsertId)
	if err != nil {
		return nil, err
	}
	row = m.DB.QueryRow(`SELECT * FROM user_favourite_episodes WHERE id=$1`, lastInsertId)
	err = row.Scan(&userFavouriteSong.ID,
		&userFavouriteSong.UserId,
		&userFavouriteSong.EpisodeId,
		&userFavouriteSong.CreatedAt,
		&userFavouriteSong.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return userFavouriteSong, nil
}

func (m *UserModel) DeleteUserFavouriteEpisode(userId int, episodeId int) error {
	_, err := m.DB.Exec(`DELETE FROM user_favourite_episodes WHERE user_id=$1 AND episode_id=$2`, userId, episodeId)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) GetUserFavouriteEpisodes(userId int) ([]PodcastEpisode, error) {
	queryBuilder := &queryHelpers.QueryBuilder{DB: m.DB}
	rows, err := queryBuilder.
		Select("episode_id").
		FromTable("user_favourite_episodes").
		WhereColumn("user_id").
		Equal(fmt.Sprintf("%v", userId)).
		GetMany()
	if err == nil {
		var episodes []PodcastEpisode
		for rows.Next() {
			var episodeId int
			err = rows.Scan(&episodeId)
			log.Print(episodeId)
			if err == nil {
				var episode PodcastEpisode
				row := queryBuilder.
					Select("*").
					FromTable("podcast_episodes").
					WhereColumn("id").
					Equal(fmt.Sprintf("%v", episodeId)).
					GetOne()
				err = row.Scan(
					&episode.ID,
					&episode.UUID,
					&episode.PodcastId,
					&episode.EpisodeName,
					&episode.EpisodeNo,
					&episode.EpisodeDesc,
					&episode.SourceURL,
					&episode.CreatedAt,
					&episode.UpdatedAt,
				)
				if err == nil {
					episodes = append(episodes, episode)
				}
			}
		}
		return episodes, nil
	}
	return nil, err
}

func (m *UserModel) GetUserLikedPosts(userId int) ([]Post, error) {
	// rows, err := m.DB.Query("SELECT post_id FROM posts_likes WHERE user_id=$1", userId)
	// if err != nil {
	// 	return nil, err
	// }
	// for rows.Next() {
	// 	var postId int
	// 	err = rows.Scan(&postId)
	// 	if err == nil {
	// 		var post Post
	// 		rows, err = m.DB.Query("SELECT id, user_id, slug, title, short_content, thumbnail_url, created_at FROM posts WHERE id=$1", postId)
	// 		if err == nil {
	// 			for rows.Next() {

	// 			}
	// 		}
	// 	}
	// }
	return nil, nil
}
