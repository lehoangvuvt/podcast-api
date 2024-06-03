package models

import (
	"database/sql"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Genre struct {
	ID        int            `json:"id"`
	UUID      string         `json:"uuid"`
	GenreName string         `json:"genre_name"`
	GenreDesc string         `json:"genre_desc"`
	BgImage   string         `json:"bg_image"`
	CreatedAt string         `json:"-"`
	UpdatedAt sql.NullString `json:"-"`
}

type GenreDetails struct {
	ID        int            `json:"id"`
	UUID      string         `json:"uuid"`
	GenreName string         `json:"genre_name"`
	GenreDesc string         `json:"genre_desc"`
	BgImage   string         `json:"bg_image"`
	Podcasts  []Podcast      `json:"podcasts"`
	CreatedAt string         `json:"-"`
	UpdatedAt sql.NullString `json:"-"`
}

type CreateGenreInput struct {
	GenreName string `json:"genre_name"`
	GenreDesc string `json:"genre_desc"`
	BgImage   string `json:"bg_image"`
}

type GenreModel struct {
	DB *sql.DB
}

func (m *GenreModel) Insert(createGenreInput *CreateGenreInput) error {
	genreUUID, _ := gonanoid.New()
	stmt := "INSERT INTO genres(uuid, genre_name, genre_desc, bg_image) VALUES($1, $2, $3, $4)"
	_, err := m.DB.Exec(stmt, genreUUID, createGenreInput.GenreName, createGenreInput.GenreDesc, createGenreInput.BgImage)
	if err != nil {
		return err
	}
	return nil
}

func (m *GenreModel) GetAllGenres() ([]Genre, error) {
	rows, err := m.DB.Query("SELECT * FROM genres")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var genres []Genre
	for rows.Next() {
		var genre Genre
		err = rows.Scan(
			&genre.ID,
			&genre.UUID,
			&genre.GenreName,
			&genre.GenreDesc,
			&genre.BgImage,
			&genre.CreatedAt,
			&genre.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (m *GenreModel) GetGenreDetails(uuid string) (*GenreDetails, error) {
	genre := &GenreDetails{}
	row := m.DB.QueryRow("SELECT * FROM genres WHERE uuid=$1", uuid)
	err := row.Scan(
		&genre.ID,
		&genre.UUID,
		&genre.GenreName,
		&genre.GenreDesc,
		&genre.BgImage,
		&genre.CreatedAt,
		&genre.UpdatedAt)
	if err != nil {
		return nil, err
	}
	stmt := `SELECT p.id, 
				p.uuid, 
				p.owner_id, 
				p.podcast_name, 
				p.podcast_desc, 
				p.thumbnail_url, 
				p.created_at, 
				p.updated_at FROM podcasts AS p 
				INNER JOIN genres_podcasts 
				ON genres_podcasts.podcast_id = p.id 
				WHERE genres_podcasts.genre_id=$1`
	rows, err := m.DB.Query(stmt, genre.ID)
	if err != nil {
		return nil, err
	}
	var podcasts []Podcast
	for rows.Next() {
		var podcast Podcast
		err = rows.Scan(
			&podcast.ID,
			&podcast.UUID,
			&podcast.OwnerId,
			&podcast.PodcastName,
			&podcast.PodcastDesc,
			&podcast.ThumbnailURL,
			&podcast.CreatedAt,
			&podcast.UpdatedAt)
		if err != nil {
			return nil, err
		}
		podcasts = append(podcasts, podcast)
	}
	if podcasts == nil {
		podcasts = []Podcast{}
	}
	genre.Podcasts = podcasts
	return genre, nil
}
