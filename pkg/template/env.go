package template

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/fredrikkvalvik/nots/pkg/template/eval"
	"github.com/fredrikkvalvik/nots/pkg/template/object"
)

// returns an env standard filters and values
func newEnv() *eval.Env {
	e := eval.NewEnv()

	e.RegisterFnValue("today_date_only", fnValueTodayDateOnly)
	e.RegisterFnValue("print_env", func() (eval.Object, error) {
		keys := slices.Collect(maps.Keys(e.Symbols))
		slices.Sort(keys)
		var str strings.Builder
		for _, key := range keys {
			sym := e.Symbols[key]
			fmt.Fprintf(&str, "%s\n", sym)
		}

		return &object.ObjectString{
			Val: str.String(),
		}, nil
	})

	e.RegisterFilter("uppercase", filterUppercase)
	e.RegisterFilter("lowercase", filterLowercase)
	e.RegisterFilter("to_title", filterToTitle)

	return e
}

// evaluates time and formats it to dateOnly (yyyy-mm-dd)
func fnValueTodayDateOnly() (object.Object, error) {
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
