package template

import (
	"strings"
	"time"

	"github.com/fredrikkvalvik/nots/pkg/template/eval"
)

// returns an env standard filters and values
func newEnv() *eval.Env {
	e := eval.NewEnv()

	e.RegisterStringValue("today", time.Now().Format(time.DateOnly))

	e.RegisterFilter("uppercase", symbolUppercase)
	e.RegisterFilter("lowercase", symbolLowercase)
	e.RegisterFilter("to_title", symbolToTitle)

	return e
}

// return a string with all letters uppercased
func symbolUppercase(obj eval.Object) (eval.Object, error) {
	return &eval.ObjectString{
		Val: strings.ToUpper(obj.ToString()),
	}, nil
}

// return a string with all letters lowercased
func symbolLowercase(obj eval.Object) (eval.Object, error) {
	return &eval.ObjectString{
		Val: strings.ToLower(obj.ToString()),
	}, nil
}

// return a string with "title case"
func symbolToTitle(obj eval.Object) (eval.Object, error) {
	str := strings.ToTitle(strings.ToLower(obj.ToString()))
	return &eval.ObjectString{
		Val: str,
	}, nil
}
