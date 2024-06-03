package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type ContextKey string

const ContextUserIdKey ContextKey = "userId"

type envelop = map[string]any

type Response struct {
	w          http.ResponseWriter
	statusCode int
}

func (res *Response) status(statusCode int) *Response {
	res.statusCode = statusCode
	return res
}

func (res *Response) json(data any) {
	res.w.Header().Set("Content-Type", "application/json")
	jsn, err := json.Marshal(data)
	if err != nil {
		res.w.WriteHeader(http.StatusBadRequest)
		errorMsg, _ := json.Marshal(envelop{"error": "bad JSON format"})
		res.w.Write(errorMsg)
		return
	}
	res.w.WriteHeader(res.statusCode)
	res.w.Write(jsn)
}

func SignToken(data map[string]any, expiredInMinutes time.Duration, secretKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	for k, v := range data {
		claims[k] = v
	}
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(expiredInMinutes * time.Minute).Unix()
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func ParseToken(tokenStr string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	invalidClaim := token.Claims.Valid()
	if invalidClaim != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}
