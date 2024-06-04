package models

import (
	"database/sql"
	queryHelpers "vulh/soundcommunity/internal/utils"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type PodcastEpisode struct {
	ID          int            `json:"id"`
	UUID        string         `json:"uuid"`
	PodcastId   int            `json:"podcast_id"`
	EpisodeName string         `json:"episode_name"`
	EpisodeNo   int            `json:"episode_no"`
	EpisodeDesc string         `json:"episode_desc"`
	SourceURL   string         `json:"source_url"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   sql.NullString `json:"-"`
}

type PodcastEpisodeDetails struct {
	ID          int            `json:"id"`
	UUID        string         `json:"uuid"`
	PodcastId   int            `json:"podcast_id"`
	EpisodeName string         `json:"episode_name"`
	EpisodeNo   int            `json:"episode_no"`
	EpisodeDesc string         `json:"episode_desc"`
	SourceURL   string         `json:"source_url"`
	Podcast     Podcast        `json:"podcast"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   sql.NullString `json:"-"`
}

type CreatePodcastEpisodeInput struct {
	PodcastId   int    `json:"podcast_id"`
	EpisodeName string `json:"episode_name"`
	EpisodeNo   int    `json:"episode_no"`
	EpisodeDesc string `json:"episode_desc"`
	SourceURL   string `json:"source_url"`
}

type PodcastEpisodeModel struct {
	DB *sql.DB
}

func (m *PodcastEpisodeModel) Insert(input *CreatePodcastEpisodeInput) error {
	uuid, _ := gonanoid.New()
	stmt := `INSERT INTO podcast_episodes(uuid, podcast_id, episode_name, episode_no, episode_desc, source_url) 
			VALUES($1, $2, $3, $4, $5, $6)`
	_, err := m.DB.Exec(stmt, uuid, input.PodcastId, input.EpisodeName, input.EpisodeNo, input.EpisodeDesc, input.SourceURL)
	if err != nil {
		return err
	}
	return nil
}

func (m *PodcastEpisodeModel) GetEpisodeDetails(uuid string) (*PodcastEpisodeDetails, error) {
	episode := &PodcastEpisodeDetails{}
	row := m.DB.QueryRow("SELECT * FROM podcast_episodes WHERE uuid=$1", uuid)
	err := row.Scan(
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
	if err != nil {
		return nil, err
	}
	podcast := Podcast{}
	row = m.DB.QueryRow("SELECT * FROM podcasts WHERE id=$1", episode.PodcastId)
	err = row.Scan(
		&podcast.ID,
		&podcast.UUID,
		&podcast.OwnerId,
		&podcast.PodcastName,
		&podcast.PodcastDesc,
		&podcast.ThumbnailURL,
		&podcast.CreatedAt,
		&podcast.UpdatedAt,
	)
	if err == nil {
		episode.Podcast = podcast
	}
	return episode, nil
}

func (m *PodcastEpisodeModel) SearchEpisodesByName(name string) ([]PodcastEpisodeDetails, error) {
	var episodes []PodcastEpisodeDetails
	rows, err := m.DB.Query(`SELECT * FROM podcast_episodes WHERE episode_name ILIKE $1`, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var episode PodcastEpisodeDetails
		err = rows.Scan(
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
		if err != nil {
			return nil, err
		}
		episodes = append(episodes, episode)
	}
	for index, item := range episodes {
		var podcast Podcast
		row := m.DB.QueryRow("SELECT * FROM podcasts WHERE id=$1", item.PodcastId)
		err = row.Scan(
			&podcast.ID,
			&podcast.UUID,
			&podcast.OwnerId,
			&podcast.PodcastName,
			&podcast.PodcastDesc,
			&podcast.ThumbnailURL,
			&podcast.CreatedAt,
			&podcast.UpdatedAt,
		)
		if err == nil {
			episodes[index].Podcast = podcast
		}
	}
	return episodes, nil
}

func (m *PodcastEpisodeModel) SearchEpisodes(queryConfig *queryHelpers.QueryConfig) ([]PodcastEpisodeDetails, error) {
	var episodes []PodcastEpisodeDetails
	queryBuilder := &queryHelpers.QueryBuilder{DB: m.DB}
	queryStr := queryBuilder.
		FromTable(queryConfig.FromTable).
		WhereColumn(queryConfig.WhereColumnName).
		Search(queryConfig.Operator, queryConfig.SearchValue).
		OrderBy(queryConfig.OrderByColumnName, queryConfig.Direction).
		Skip(queryConfig.Skip).
		Limit(queryConfig.Limit).
		GetQuery()
	rows, err := m.DB.Query(queryStr)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var episode PodcastEpisodeDetails
		err = rows.Scan(
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
		if err != nil {
			return nil, err
		}
		episodes = append(episodes, episode)
	}
	for index, item := range episodes {
		var podcast Podcast
		row := m.DB.QueryRow("SELECT * FROM podcasts WHERE id=$1", item.PodcastId)
		err = row.Scan(
			&podcast.ID,
			&podcast.UUID,
			&podcast.OwnerId,
			&podcast.PodcastName,
			&podcast.PodcastDesc,
			&podcast.ThumbnailURL,
			&podcast.CreatedAt,
			&podcast.UpdatedAt,
		)
		if err == nil {
			episodes[index].Podcast = podcast
		}
	}
	return episodes, nil
}
