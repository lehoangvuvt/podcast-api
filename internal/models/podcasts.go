package models

import (
	"database/sql"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Podcast struct {
	ID           int            `json:"id"`
	UUID         string         `json:"uuid"`
	OwnerId      int            `json:"owner_id"`
	PodcastName  string         `json:"podcast_name"`
	PodcastDesc  string         `json:"podcast_desc"`
	ThumbnailURL string         `json:"thumbnail_url"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    sql.NullString `json:"-"`
}

type PodcastDetails struct {
	ID           int              `json:"id"`
	UUID         string           `json:"uuid"`
	OwnerId      int              `json:"owner_id"`
	PodcastName  string           `json:"podcast_name"`
	PodcastDesc  string           `json:"podcast_desc"`
	ThumbnailURL string           `json:"thumbnail_url"`
	Episodes     []PodcastEpisode `json:"episodes"`
	CreatedAt    string           `json:"created_at"`
	UpdatedAt    sql.NullString   `json:"-"`
}

type CreatePodcastInput struct {
	OwnerId      int    `json:"owner_id"`
	PodcastName  string `json:"podcast_name"`
	PodcastDesc  string `json:"podcast_desc"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type PodcastModel struct {
	DB *sql.DB
}

func (m *PodcastModel) Insert(input *CreatePodcastInput) error {
	uuid, _ := gonanoid.New()
	stmt := "INSERT INTO podcasts (uuid, owner_id, podcast_name, podcast_desc, thumbnail_url) VALUES($1, $2, $3, $4, $5)"
	_, err := m.DB.Exec(stmt, uuid, input.OwnerId, input.PodcastName, input.PodcastDesc, input.ThumbnailURL)
	if err != nil {
		return err
	}
	return nil
}

func (m *PodcastModel) GetAllPodcasts() ([]Podcast, error) {
	rows, err := m.DB.Query("SELECT * FROM podcasts")
	if err != nil {
		return nil, err
	}
	var podcasts []Podcast
	for rows.Next() {
		var podcast Podcast
		err = rows.Scan(&podcast.ID,
			&podcast.UUID,
			&podcast.OwnerId,
			&podcast.PodcastName,
			&podcast.PodcastDesc,
			&podcast.ThumbnailURL,
			&podcast.CreatedAt,
			&podcast.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		podcasts = append(podcasts, podcast)
	}
	return podcasts, nil
}

func (m *PodcastModel) GetPodcastDetails(uuid string) (*PodcastDetails, error) {
	podcast := &PodcastDetails{}
	row := m.DB.QueryRow("SELECT * FROM podcasts WHERE uuid=$1", uuid)
	err := row.Scan(&podcast.ID,
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
	var podcastEpisodes []PodcastEpisode
	rows, err := m.DB.Query("SELECT * FROM podcast_episodes WHERE podcast_id=$1", podcast.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var episode PodcastEpisode
		err = rows.Scan(&episode.ID,
			&episode.UUID,
			&episode.PodcastId,
			&episode.EpisodeName,
			&episode.EpisodeNo,
			&episode.EpisodeDesc,
			&episode.SourceURL,
			&episode.CreatedAt,
			&episode.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		podcastEpisodes = append(podcastEpisodes, episode)
	}
	podcast.Episodes = podcastEpisodes
	return podcast, nil
}

func (m *PodcastModel) SearchPodcastsByName(name string) ([]Podcast, error) {
	var podcasts []Podcast
	rows, err := m.DB.Query(`SELECT * FROM podcasts WHERE podcast_name ILIKE $1`, "%"+name+"%")
	if err != nil {
		return nil, err
	}
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
	return podcasts, nil
}
