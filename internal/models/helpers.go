package models

import (
	"regexp"
	"strings"
)

func removePunctuation(input string) string {
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
	input = strings.ReplaceAll(input, "?", "")
	input = strings.ReplaceAll(input, ":", "")
	input = strings.ReplaceAll(input, "'", "")
	input = strings.ReplaceAll(input, `"`, "")
	input = strings.ReplaceAll(input, `”`, "")
	return strings.ToLower(input)
}
