package main

import (
	"encoding/json"
	"net/http"
	"vulh/soundcommunity/internal/models"
)

func (app *application) createGenrePodcastHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	input := &models.CreateGenerePodcastInput{}
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		errMsg, _ := json.Marshal(envelop{"error": "bad JSON format"})
		res.status(http.StatusBadRequest).json(errMsg)
	}
	err = app.models.GenerePodcastModel.Insert(input)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusCreated).json(envelop{"message": "create genre <-> podcast success"})
}
