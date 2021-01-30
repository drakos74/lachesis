package bench

import (
	"fmt"
	"strings"
)

func Format(attributes map[string]int, labels ...string) string {
	var sb strings.Builder
	values := make([]interface{}, 0)

	// make sure we start with the right delimiter to get rid of the default benchmark prefix
	sb.WriteString("|")
	for _, label := range labels {
		sb.WriteString(fmt.Sprintf("%s|", label))
	}
	for label, value := range attributes {
		sb.WriteString(fmt.Sprintf("%s:%d|", label, value))
		values = append(values, value)
	}
	str := sb.String()
	return str
}
