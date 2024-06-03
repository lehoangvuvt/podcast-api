package main

import (
	"context"
	"log"
	"net/http"
)

func (app *application) AuthGuard(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthGuard")
		res := &Response{w: w}
		cookieValue, err := r.Cookie("access_token")

		if err != nil {
			res.status(http.StatusUnauthorized).json(envelop{"error": "unauthorized user"})
			return
		}
		tokenStr := cookieValue.Value
		claims, err := ParseToken(tokenStr, "access_token")
		if err != nil {
			res.status(http.StatusUnauthorized).json(envelop{"error": "unauthorized user"})
			return
		}
		userId := -1
		for k, v := range claims {
			if k == "sub" {
				userId = int(v.(float64))
			}
		}
		if userId == -1 {
			res.status(http.StatusUnauthorized).json(envelop{"error": "unauthorized user"})
			return
		}
		ctx := context.WithValue(r.Context(), ContextUserIdKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
