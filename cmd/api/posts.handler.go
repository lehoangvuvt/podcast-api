package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"vulh/soundcommunity/internal/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}

	userId := r.Context().Value(ContextUserIdKey)
	input := &models.CreatePostInput{}
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		errMsg, _ := json.Marshal(envelop{"error": "bad JSON format"})
		res.status(http.StatusBadRequest).json(errMsg)
	}
	err = app.models.PostModel.Insert(input, userId.(int))
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusCreated).json(envelop{"message": "create post success"})
}

func (app *application) getPostBySlugHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	params := httprouter.ParamsFromContext(r.Context())
	slug := params.ByName("slug")
	post, err := app.models.PostModel.GetPostBySlug(slug)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	var relativePosts []models.PostWithUserInfo
	relativePostMap := make(map[int]models.PostWithUserInfo)
	for _, topic := range post.Topics {
		rPosts, err := app.models.PostModel.GetPostsByTopic(topic.Slug)
		if err == nil {
			for _, rPost := range rPosts {
				if rPost.ID == post.ID {
					continue
				}
				relativePostMap[rPost.ID] = rPost
			}
		}
	}

	for _, post := range relativePostMap {
		relativePosts = append(relativePosts, post)
	}

	res.status(http.StatusOK).json(envelop{"post": post, "relative_posts": relativePosts})
}

func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	params := httprouter.ParamsFromContext(r.Context())
	q := params.ByName("q")
	posts, err := app.models.PostModel.GetPosts(q)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusOK).json(envelop{"posts": posts})
}

func (app *application) getPostsByTopicHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	params := httprouter.ParamsFromContext(r.Context())
	topicSlug := params.ByName("slug")
	posts, err := app.models.PostModel.GetPostsByTopic(topicSlug)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusOK).json(envelop{"posts": posts})
}

func (app *application) getPostLikesByPostIdHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	params := httprouter.ParamsFromContext(r.Context())
	postId, err := strconv.Atoi(params.ByName("id"))
	if err != nil || postId < 0 {
		res.status(http.StatusBadRequest).json(envelop{"error": "invalid post id"})
		return
	}
	postLikes, err := app.models.PostLikeModel.GetPostLikesByPostId(postId)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusOK).json(envelop{"post_likes": postLikes})
}

func (app *application) getCommentsByPostIdHandler(w http.ResponseWriter, r *http.Request) {
	res := &Response{w: w}
	params := httprouter.ParamsFromContext(r.Context())
	postId, err := strconv.Atoi(params.ByName("id"))
	if err != nil || postId < 0 {
		res.status(http.StatusBadRequest).json(envelop{"error": "invalid post id"})
		return
	}
	comments, err := app.models.PostCommentModel.GetCommentsByPostId(postId)
	if err != nil {
		res.status(http.StatusBadRequest).json(envelop{"error": err.Error()})
		return
	}
	res.status(http.StatusOK).json(envelop{"comments": comments})
}
