package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"vulh/soundcommunity/internal/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) likePostHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	params := httprouter.ParamsFromContext(r.Context())
	postId, err := strconv.Atoi(params.ByName("id"))
	if err != nil || postId < 0 {
		res.status(http.StatusBadRequest).json(envelop{"error": "invalid post id"})
		return
	}
	userId := r.Context().Value(ContextUserIdKey)
	err = app.models.PostLikeModel.Insert(userId.(int), postId)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err})
		return
	}
	res.status(http.StatusCreated).json(envelop{"message": "like post success"})
}

func (app *application) unlikePostHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	params := httprouter.ParamsFromContext(r.Context())
	postId, err := strconv.Atoi(params.ByName("id"))
	if err != nil || postId < 0 {
		res.status(http.StatusBadRequest).json(envelop{"error": "invalid post id"})
		return
	}
	userId := r.Context().Value(ContextUserIdKey)
	err = app.models.PostLikeModel.Delete(userId.(int), postId)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err})
		return
	}
	res.status(http.StatusOK).json(envelop{"message": "unlike post success"})
}

func (app *application) createCommentToPostHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	params := httprouter.ParamsFromContext(r.Context())
	postId, err := strconv.Atoi(params.ByName("id"))
	if err != nil || postId < 0 {
		res.status(http.StatusBadRequest).json(envelop{"error": "invalid post id"})
		return
	}
	userId := r.Context().Value(ContextUserIdKey)
	input := &models.CreateCommentInput{}
	err = json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		errMsg, _ := json.Marshal(envelop{"error": "bad JSON format"})
		res.status(http.StatusBadRequest).json(errMsg)
	}
	err = app.models.PostCommentModel.Insert(userId.(int), postId, input)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err})
		return
	}
	res.status(http.StatusOK).json(envelop{"message": "create comment success"})
}
