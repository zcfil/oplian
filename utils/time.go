package utils

import (
	"strings"
	"time"
)

func WholeSecond(du time.Duration) string {
	str := du.String()
	if strings.HasSuffix(str, "s") {
		for i, v := range str {
			if v == '.' {
				return str[:i] + "s"
			}
		}
	}
	return str
}
