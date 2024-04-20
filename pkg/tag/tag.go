package tag

import "strings"

func ParseJsonTag(tag string) string {
	if strings.Index(tag, ",") != -1 {
		split := strings.Split(tag, ",")
		tag = split[0]
	}
	return tag
}
