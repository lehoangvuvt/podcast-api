package main

import (
	"net/http"
	"vulh/soundcommunity/internal/models"

	"github.com/julienschmidt/httprouter"
)

type SearchResult struct {
	Podcasts []models.Podcast               `json:"podcasts"`
	Episodes []models.PodcastEpisodeDetails `json:"episodes"`
	Genres   []models.Genre                 `json:"genres"`
}

func (app *application) SearchHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	result := &SearchResult{}
	params := httprouter.ParamsFromContext(r.Context())
	q := params.ByName("q")
	podcasts, err := app.models.PodcastModel.SearchPodcastsByName(q)
	if err == nil {
		result.Podcasts = podcasts
	}
	episodes, err := app.models.PodcastEpisodeModel.SearchEpisodesByName(q)
	if err == nil {
		result.Episodes = episodes
	}
	genres, err := app.models.GenreModel.SearchGenresByName(q)
	if err == nil {
		result.Genres = genres
	}
	res.status(http.StatusOK).json(envelop{"result": result})
}
