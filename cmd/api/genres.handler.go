package main

import (
	"encoding/json"
	"net/http"
	"vulh/soundcommunity/internal/models"
)

func (app *application) createGenreHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	input := &models.CreateGenreInput{}
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		errMsg, _ := json.Marshal(envelop{"error": "bad JSON format"})
		res.status(http.StatusBadRequest).json(errMsg)
	}
	err = app.models.GenreModel.Insert(input)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusCreated).json(envelop{"message": "create genre success"})
}

func (app *application) getAllGenres(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	genres, err := app.models.GenreModel.GetAllGenres()
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusOK).json(envelop{"genres": genres})
}
