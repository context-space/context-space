package utils

import "strings"

var sb strings.Builder

func StringsBuilder(strs ...string) string {
	defer sb.Reset()
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
