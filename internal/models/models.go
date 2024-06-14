package models

import "database/sql"

type Models struct {
	UserModel           *UserModel
	GenreModel          *GenreModel
	PodcastModel        *PodcastModel
	PodcastEpisodeModel *PodcastEpisodeModel
	GenerePodcastModel  *GenerePodcastModel
	PostModel           *PostModel
	TopicModel          *TopicModel
	PostLikeModel       *PostLikeModel
}

func NewModels(DB *sql.DB) *Models {
	return &Models{
		UserModel:           &UserModel{DB: DB},
		GenreModel:          &GenreModel{DB: DB},
		PodcastModel:        &PodcastModel{DB: DB},
		PodcastEpisodeModel: &PodcastEpisodeModel{DB: DB},
		GenerePodcastModel:  &GenerePodcastModel{DB: DB},
		PostModel:           &PostModel{DB: DB},
		TopicModel:          &TopicModel{DB: DB},
		PostLikeModel:       &PostLikeModel{DB: DB},
	}
}
