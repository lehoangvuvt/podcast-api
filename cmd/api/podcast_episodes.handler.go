package main

import (
	"encoding/json"
	"net/http"
	"vulh/soundcommunity/internal/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) createPodcastEpisodeHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	input := &models.CreatePodcastEpisodeInput{}
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		errMsg, _ := json.Marshal(envelop{"error": "bad JSON format"})
		res.status(http.StatusBadRequest).json(errMsg)
	}
	err = app.models.PodcastEpisodeModel.Insert(input)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusCreated).json(envelop{"message": "create podcast episode success"})
}

func (app *application) getEpisodeDetailsHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	params := httprouter.ParamsFromContext(r.Context())
	uuid := params.ByName("uuid")
	episodeDetails, err := app.models.PodcastEpisodeModel.GetEpisodeDetails(uuid)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusCreated).json(envelop{"episode_details": episodeDetails})
}
