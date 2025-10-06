package template

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/fredrikkvalvik/nots/pkg/template/eval"
	"github.com/fredrikkvalvik/nots/pkg/template/object"
)

// returns an env standard filters and values
func newEnv() *eval.Env {
	e := eval.NewEnv()

	e.RegisterFnValue("today_date_only", "returns todays date in yyyy-mm-dd format", fnValueTodayDateOnly)
	e.RegisterFnValue("print_env", "print all available symbols", func() (eval.Object, error) {
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

	e.RegisterFilter("uppercase", "return the input string in all uppercase", filterUppercase)
	e.RegisterFilter("lowercase", "return the input string in all lowercase", filterLowercase)
	e.RegisterFilter("to_title", "return the input string with each word capitalized", filterToTitle)
	e.RegisterFilter("sh", "executes the input string as a shell command in the current environment", filterExec)

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

// filter for executing shell commands.
// no safeguards. no guardrails. careful.
func filterExec(obj object.Object) (object.Object, error) {
	strObj, ok := obj.(*object.ObjectString)
	if !ok {
		return nil, fmt.Errorf("cannot execute non-string command")
	}
	str := strObj.Val

	if len(str) == 0 {
		return nil, fmt.Errorf("cannot execute empty string")
	}

	cmd := exec.Command(os.ExpandEnv("$SHELL"), "-c", str)
	res, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return &object.ObjectString{Val: string(res)}, nil

}
