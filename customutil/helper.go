package customutil

import "strings"

func FormatJsonForRedis(strData string) string {
	strData = strings.ReplaceAll(strData, "\n", "")
	strData = strings.ReplaceAll(strData, "\t", "")

	return strData
}
