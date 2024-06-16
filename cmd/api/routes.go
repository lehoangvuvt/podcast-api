package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://podcast-client.vercel.app", "https://podcast.healing-journey.asia", "https://stories.healing-journey.asia"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})

	router.HandlerFunc(http.MethodPost, "/users", app.createUserHandler)
	router.HandlerFunc(http.MethodPost, "/login", app.loginHandler)
	router.HandlerFunc(http.MethodGet, "/users/authenticate", app.AuthGuard(app.authenticate))
	router.HandlerFunc(http.MethodGet, "/users/invalidate", app.AuthGuard(app.invalidateHandler))
	router.HandlerFunc(http.MethodPost, "/users/favourites", app.AuthGuard(app.modifyUserFavouriteItemHandler))
	router.HandlerFunc(http.MethodGet, "/users/favourites", app.AuthGuard(app.getUserFavouriteItemsHandler))

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

	router.HandlerFunc(http.MethodPost, "/posts", app.AuthGuard(app.createPostHandler))
	router.HandlerFunc(http.MethodGet, "/posts/search", app.getPostsHandler)
	router.HandlerFunc(http.MethodGet, "/posts/post/:slug", app.getPostBySlugHandler)
	router.HandlerFunc(http.MethodDelete, "/posts/likes/:id", app.AuthGuard(app.unlikePostHandler))
	router.HandlerFunc(http.MethodPut, "/posts/likes/:id", app.AuthGuard(app.likePostHandler))
	router.HandlerFunc(http.MethodGet, "/posts/likes/:id", app.getPostLikesByPostIdHandler)
	router.HandlerFunc(http.MethodPost, "/posts/comments/:id", app.AuthGuard(app.createCommentToPostHandler))
	router.HandlerFunc(http.MethodGet, "/posts/comments/:id", app.getCommentsByPostIdHandler)

	router.HandlerFunc(http.MethodGet, "/topics", app.getAllTopicsHandler)
	router.HandlerFunc(http.MethodGet, "/topics/search/:q", app.searchTopicsByNameHandler)
	router.HandlerFunc(http.MethodGet, "/topics/posts/:slug", app.getPostsByTopicHandler)
	router.HandlerFunc(http.MethodGet, "/topics/relative/:slug", app.getRelativeTopicsHandler)
	router.HandlerFunc(http.MethodGet, "/topics/recommended", app.getRecommendedTopics)

	router.HandlerFunc(http.MethodPost, "/upload", app.AuthGuard(app.uploadFileHandler))

	handler := crs.Handler(router)
	return handler
}
