package models

import "database/sql"

type Models struct {
	UserModel           *UserModel
	GenreModel          *GenreModel
	PodcastModel        *PodcastModel
	PodcastEpisodeModel *PodcastEpisodeModel
	GenerePodcastModel  *GenerePodcastModel
}

func NewModels(DB *sql.DB) *Models {
	return &Models{
		UserModel:           &UserModel{DB: DB},
		GenreModel:          &GenreModel{DB: DB},
		PodcastModel:        &PodcastModel{DB: DB},
		PodcastEpisodeModel: &PodcastEpisodeModel{DB: DB},
		GenerePodcastModel:  &GenerePodcastModel{DB: DB},
	}
}
