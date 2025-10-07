package timeformat

import (
	"fmt"
	"strings"
	"time"
)

func Format(t time.Time, format string) string {
	formatReplacers := []struct {
		token string
		To    func() string
	}{
		// Order is important. YYYY before YY to make sure we dont set short year twice, etc
		{"YYYY", func() string { return fmt.Sprint(t.Year()) }},
		{"YY", func() string { return fmt.Sprint(t.Year())[0:2] }},
		{"MM", func() string { return fmt.Sprintf("%02d", t.Month()) }},
		{"M", func() string { return fmt.Sprintf("%d", t.Month()) }},
		{"DD", func() string { return fmt.Sprintf("%02d", t.Day()) }},
		{"D", func() string { return fmt.Sprintf("%d", t.Day()) }},

		{"hh", func() string { return fmt.Sprintf("%02d", t.Hour()) }},
		{"h", func() string { return fmt.Sprintf("%d", t.Hour()) }},

		{"mm", func() string { return fmt.Sprintf("%02d", t.Minute()) }},
		{"m", func() string { return fmt.Sprintf("%d", t.Minute()) }},

		{"ss", func() string { return fmt.Sprintf("%02d", t.Second()) }},
		{"s", func() string { return fmt.Sprintf("%d", t.Second()) }},

		{"ww", func() string { _, week := t.ISOWeek(); return fmt.Sprintf("%02d", week) }},
		{"w", func() string { _, week := t.ISOWeek(); return fmt.Sprintf("%d", week) }},
	}

	for _, replacer := range formatReplacers {
		// we only need to run a replacer when the a token is present
		if strings.Contains(format, replacer.token) {
			format = strings.ReplaceAll(format, replacer.token, replacer.To())
		}
	}

	return format
}
