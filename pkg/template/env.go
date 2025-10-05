package template

import (
	"strings"
	"time"

	"github.com/fredrikkvalvik/nots/pkg/template/eval"
	"github.com/fredrikkvalvik/nots/pkg/template/object"
)

// returns an env standard filters and values
func newEnv() *eval.Env {
	e := eval.NewEnv()

	e.RegisterFnValue("today_date_only", fnValueTodayDateOnly)

	e.RegisterFilter("uppercase", filterUppercase)
	e.RegisterFilter("lowercase", filterLowercase)
	e.RegisterFilter("to_title", filterToTitle)

	return e
}

// evaluates time and formats it to dateOnly (yyyy-mm-dd)
func fnValueTodayDateOnly() (eval.Object, error) {
	return &object.ObjectString{
		Val: time.Now().Format(time.DateOnly),
	}, nil
}

// return a string with all letters uppercased
func filterUppercase(obj object.Object) (object.Object, error) {
	return &object.ObjectString{
		Val: strings.ToUpper(obj.ToString()),
	}, nil
}

// return a string with all letters lowercased
func filterLowercase(obj object.Object) (object.Object, error) {
	return &object.ObjectString{
		Val: strings.ToLower(obj.ToString()),
	}, nil
}

// return a string with "title case"
func filterToTitle(obj object.Object) (object.Object, error) {
	str := strings.ToTitle(strings.ToLower(obj.ToString()))
	return &object.ObjectString{
		Val: str,
	}, nil
}
