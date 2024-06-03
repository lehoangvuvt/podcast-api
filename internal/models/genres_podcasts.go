package models

import (
	"database/sql"
)

type GenerePodcast struct {
	ID        int            `json:"id"`
	GenreId   int            `json:"genre_id"`
	PodcastId int            `json:"podcast_id"`
	CreatedAt string         `json:"-"`
	UpdatedAt sql.NullString `json:"-"`
}

type CreateGenerePodcastInput struct {
	GenreIds  []int `json:"genres_ids"`
	PodcastId int   `json:"podcast_id"`
}

type GenerePodcastModel struct {
	DB *sql.DB
}

func (m *GenerePodcastModel) Insert(input *CreateGenerePodcastInput) error {
	tx, err := m.DB.Begin()

	if err != nil {
		return err
	}

	for _, genreId := range input.GenreIds {
		_, err := tx.Exec("INSERT INTO genres_podcasts(genre_id, podcast_id) VALUES($1, $2)", genreId, input.PodcastId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
