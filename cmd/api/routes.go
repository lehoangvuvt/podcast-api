package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://podcast-client.vercel.app"},
		AllowCredentials: true,
	})

	router.HandlerFunc(http.MethodPost, "/users", app.createUserHandler)
	router.HandlerFunc(http.MethodPost, "/login", app.loginHandler)
	router.HandlerFunc(http.MethodGet, "/users/authenticate", app.AuthGuard(app.authenticate))
	router.HandlerFunc(http.MethodGet, "/users/invalidate", app.AuthGuard(app.invalidateHandler))
	router.HandlerFunc(http.MethodPost, "/users/favourites", app.AuthGuard(app.modifyUserFavouriteItemHandler))

	router.HandlerFunc(http.MethodPost, "/genres", app.createGenreHandler)
	router.HandlerFunc(http.MethodGet, "/genres", app.getAllGenres)
	router.HandlerFunc(http.MethodGet, "/genres/:uuid", app.getGenreDetailsHandler)

	router.HandlerFunc(http.MethodGet, "/podcasts", app.getAllPodcastsHandler)
	router.HandlerFunc(http.MethodGet, "/podcasts/:uuid", app.getPodcastDetailsHandler)
	router.HandlerFunc(http.MethodPost, "/podcasts", app.createPodcastHandler)
	router.HandlerFunc(http.MethodPut, "/podcasts/genres", app.createGenrePodcastHandler)
	router.HandlerFunc(http.MethodPost, "/podcasts/episode", app.createPodcastEpisodeHandler)

	router.HandlerFunc(http.MethodGet, "/search/:q", app.SearchHandler)

	router.HandlerFunc(http.MethodGet, "/episodes/:uuid", app.getEpisodeDetailsHandler)

	router.HandlerFunc(http.MethodGet, "/relative/episodes", app.getRelativeEpisodesHandler)

	router.HandlerFunc(http.MethodGet, "/home/feeds", app.GetHomeFeedsHandler)

	handler := crs.Handler(router)
	return handler
}
