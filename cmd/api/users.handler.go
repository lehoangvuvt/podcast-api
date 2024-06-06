package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"vulh/soundcommunity/internal/models"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	input := &models.CreateUserInput{}
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		errMsg, _ := json.Marshal(envelop{"error": "bad JSON format"})
		res.status(http.StatusBadRequest).json(errMsg)
	}
	err = app.models.UserModel.Insert(input)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			if strings.Contains(err.Error(), "username") {
				res.status(http.StatusBadRequest).json(envelop{"error": "username existed"})
			} else {
				res.status(http.StatusBadRequest).json(envelop{"error": "email existed"})
			}
		}
		return
	}
	res.status(http.StatusCreated).json(envelop{"message": "create user success"})
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	accessTokenExpireDuration := time.Duration(60)
	res := &Response{w: w}
	input := &models.LoginInput{}
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		errMsg, _ := json.Marshal(envelop{"error": "bad JSON format"})
		res.status(http.StatusBadRequest).json(errMsg)
	}
	user, err := app.models.UserModel.Login(input)
	if err != nil {
		res.status(http.StatusUnauthorized).json(envelop{"error": err.Error()})
		return
	}
	tokenStr, err := SignToken(envelop{"sub": user.ID}, accessTokenExpireDuration, "access_token")
	if err != nil {
		res.status(http.StatusUnauthorized).json(envelop{"error": "generate token error"})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenStr,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Path:     "/",
		MaxAge:   60 * 60,
	})
	res.status(http.StatusCreated).json(user)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	userId := r.Context().Value(ContextUserIdKey)
	user, err := app.models.UserModel.GetUserById(userId.(int))
	if err != nil {
		errMsg, _ := json.Marshal(envelop{"error": "user id not found"})
		res.status(http.StatusUnauthorized).json(errMsg)
		return
	}
	res.status(http.StatusAccepted).json(user)
}

func (app *application) invalidateHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	c := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   0,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	http.SetCookie(w, c)
	res.status(http.StatusAccepted).json(envelop{"message": "invalidate user success"})
}
