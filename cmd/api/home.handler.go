package main

import (
	"net/http"
	queryHelpers "vulh/soundcommunity/internal/utils"
)

func (app *application) GetHomeFeedsHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}

	podcasts, err := app.models.PodcastModel.GetPodcastsHomeFeeds(&queryHelpers.QueryConfig{
		FromTable:         "podcasts",
		WhereColumnName:   "podcast_name",
		SearchValue:       "*",
		OrderByColumnName: "created_at",
		Direction:         queryHelpers.QueryDirection.DESC,
		Skip:              0,
		Limit:             5,
	})

	if err != nil {
		res.status(404).json(envelop{"error": "cannot get home feeds"})
		return
	}

	episodes, _ := app.models.PodcastEpisodeModel.SearchEpisodes(
		&queryHelpers.QueryConfig{
			FromTable:         "podcast_episodes",
			WhereColumnName:   "podcast_name",
			SearchValue:       "*",
			OrderByColumnName: "created_at",
			Direction:         queryHelpers.QueryDirection.DESC,
			Skip:              0,
			Limit:             4,
		})

	res.status(http.StatusOK).json(envelop{"podcasts": podcasts, "episodes": episodes})
}
