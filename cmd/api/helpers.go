package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
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

func RemoveVNPunctuation(input string) string {
	var Regexp_A = `à|á|ạ|ã|ả|ă|ắ|ằ|ẳ|ẵ|ặ|â|ấ|ầ|ẩ|ẫ|ậ`
	var Regexp_E = `è|ẻ|ẽ|é|ẹ|ê|ề|ể|ễ|ế|ệ`
	var Regexp_I = `ì|ỉ|ĩ|í|ị`
	var Regexp_U = `ù|ủ|ũ|ú|ụ|ư|ừ|ử|ữ|ứ|ự`
	var Regexp_Y = `ỳ|ỷ|ỹ|ý|ỵ`
	var Regexp_O = `ò|ỏ|õ|ó|ọ|ô|ồ|ổ|ỗ|ố|ộ|ơ|ờ|ở|ỡ|ớ|ợ`
	var Regexp_D = `Đ|đ`
	reg_a := regexp.MustCompile(Regexp_A)
	reg_e := regexp.MustCompile(Regexp_E)
	reg_i := regexp.MustCompile(Regexp_I)
	reg_o := regexp.MustCompile(Regexp_O)
	reg_u := regexp.MustCompile(Regexp_U)
	reg_y := regexp.MustCompile(Regexp_Y)
	reg_d := regexp.MustCompile(Regexp_D)
	input = reg_a.ReplaceAllLiteralString(input, "a")
	input = reg_e.ReplaceAllLiteralString(input, "e")
	input = reg_i.ReplaceAllLiteralString(input, "i")
	input = reg_o.ReplaceAllLiteralString(input, "o")
	input = reg_u.ReplaceAllLiteralString(input, "u")
	input = reg_y.ReplaceAllLiteralString(input, "y")
	input = reg_d.ReplaceAllLiteralString(input, "d")

	// regexp remove charaters in ()
	var RegexpPara = `\(.*\)`
	reg_para := regexp.MustCompile(RegexpPara)
	input = reg_para.ReplaceAllLiteralString(input, "")

	return strings.ToLower(input)
}
