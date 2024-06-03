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
