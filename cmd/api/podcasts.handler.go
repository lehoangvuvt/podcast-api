package main

import (
	"encoding/json"
	"net/http"
	"vulh/soundcommunity/internal/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) createPodcastHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	input := &models.CreatePodcastInput{}
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		errMsg, _ := json.Marshal(envelop{"error": "bad JSON format"})
		res.status(http.StatusBadRequest).json(errMsg)
	}
	err = app.models.PodcastModel.Insert(input)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusCreated).json(envelop{"message": "create podcast success"})
}

func (app *application) getAllPodcastsHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	podcasts, err := app.models.PodcastModel.GetAllPodcasts()
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusOK).json(envelop{"podcasts": podcasts})
}

func (app *application) getPodcastDetailsHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	params := httprouter.ParamsFromContext(r.Context())
	uuid := params.ByName("uuid")
	podcastDetails, err := app.models.PodcastModel.GetPodcastDetails(uuid)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusOK).json(envelop{"podcast_details": podcastDetails})
}

func (app *application) SearchPodcastsByNameHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	params := httprouter.ParamsFromContext(r.Context())
	q := params.ByName("q")
	podcasts, err := app.models.PodcastModel.SearchPodcastsByName(q)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusOK).json(envelop{"podcasts": podcasts})
}
