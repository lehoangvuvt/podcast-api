package main

import (
	"encoding/json"
	"net/http"
)

type UserFavouriteItemType string

var UserFavouriteItem = struct {
	Episode UserFavouriteItemType
	Podcast UserFavouriteItemType
}{
	Episode: "Episode",
	Podcast: "Podcast",
}

type UserFavouriteOperatorType string

var UserFavouriteOperator = struct {
	Add    UserFavouriteOperatorType
	Remove UserFavouriteOperatorType
}{
	Add:    "Add",
	Remove: "Remove",
}

type ModifyUserFavouriteItemInput struct {
	Type     UserFavouriteItemType     `json:"type"`
	ItemId   int                       `json:"item_id"`
	Operator UserFavouriteOperatorType `json:"operator"`
}

func (app *application) modifyUserFavouriteItemHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	userId := r.Context().Value(ContextUserIdKey)
	input := &ModifyUserFavouriteItemInput{}
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		errMsg, _ := json.Marshal(envelop{"error": "bad JSON format"})
		res.status(http.StatusBadRequest).json(errMsg)
	}
	switch input.Type {
	case UserFavouriteItem.Episode:
		if input.Operator == UserFavouriteOperator.Add {
			userFavouriteEpisode, err := app.models.UserModel.CreateUserFavouriteEpisode(userId.(int), input.ItemId)
			if err != nil {
				res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
				return
			}
			res.status(http.StatusCreated).json(envelop{
				"type":     UserFavouriteItem.Episode,
				"operator": UserFavouriteOperator.Add,
				"item_id":  userFavouriteEpisode.EpisodeId,
			})
			return
		} else {
			err = app.models.UserModel.DeleteUserFavouriteEpisode(userId.(int), input.ItemId)
			if err != nil {
				res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
				return
			}
			res.status(http.StatusOK).json(envelop{
				"type":     UserFavouriteItem.Episode,
				"operator": UserFavouriteOperator.Remove,
				"item_id":  input.ItemId,
			})
			return
		}
	}
}