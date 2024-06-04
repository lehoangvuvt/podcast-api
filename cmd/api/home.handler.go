package main

import (
	"net/http"
	queryHelpers "vulh/soundcommunity/internal/utils"
)

func (app *application) GetHomeFeedsHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	podcasts, _ := app.models.PodcastModel.SearchPodcasts(&queryHelpers.QueryConfig{
		FromTable:         "podcasts",
		WhereColumnName:   "podcast_name",
		SearchValue:       "*",
		OrderByColumnName: "created_at",
		Direction:         queryHelpers.QueryDirection.DESC,
		Skip:              0,
		Limit:             5,
	})
	episodes, _ := app.models.PodcastEpisodeModel.SearchEpisodes(&queryHelpers.QueryConfig{
		FromTable:         "podcast_episodes",
		WhereColumnName:   "episode_name",
		SearchValue:       "*",
		OrderByColumnName: "created_at",
		Direction:         queryHelpers.QueryDirection.DESC,
		Skip:              0,
		Limit:             5,
	})
	res.status(http.StatusOK).json(envelop{"podcasts": podcasts, "episodes": episodes})
}
