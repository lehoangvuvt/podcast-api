package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) getAllTopicsHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	topics, err := app.models.TopicModel.GetAllTopics()
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusOK).json(envelop{"topics": topics})
}

func (app *application) searchTopicsByNameHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	params := httprouter.ParamsFromContext(r.Context())
	q := params.ByName("q")
	topics, err := app.models.TopicModel.SearchTopicsByName(q)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusOK).json(envelop{"topics": topics})
}

func (app *application) getRelativeTopicsHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	params := httprouter.ParamsFromContext(r.Context())
	slug := params.ByName("slug")
	topics, err := app.models.TopicModel.GetRelativeTopics(slug)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusOK).json(envelop{"topics": topics})
}

func (app *application) getRecommendedTopics(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	topics, err := app.models.TopicModel.GetRecommendedTopics()
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusOK).json(envelop{"topics": topics})
}
