package gossi

import (
	"strings"
)

func echo(s string) string {
	return s
}

func uppercase(s string) string {
	return strings.ToUpper(s)
}

func lowercase(s string) string {
	return strings.ToLower(s)
}

const sign = " by gossi"

func Run(program string, data string) string {
	switch program {
	case "echo":
		return echo(data) + sign

	case "uppercase":
		return uppercase(data) + sign

	case "lowercase":
		return lowercase(data) + sign

	default:
		return "" + sign

	}
}
