package util

import (
	"regexp"
	"strings"
)

func RemoveHTMLTags(description string) string {
	re := regexp.MustCompile(`<a href=".*?">(.*?)</a>`)
	description = re.ReplaceAllString(description, "$1")

	description = strings.ReplaceAll(description, "<code>", "*")
	description = strings.ReplaceAll(description, "</code>", "*")

	description = strings.ReplaceAll(description, "<em>", "*")
	description = strings.ReplaceAll(description, "</em>", "*")

	return description
}
