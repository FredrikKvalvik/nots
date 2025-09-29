package datetemplate

import (
	"fmt"
	"strings"
	"time"
)

func Generate(template string) string {
	t := time.Now()

	result := strings.ReplaceAll(template, "yyyy", fmt.Sprint(t.Year()))
	// result = strings.Re

	return result
}
