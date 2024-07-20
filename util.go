package main

import "strings"

// edit path to be accepted by OS
func sanitizePath(path string) string {
	sanitized := strings.ReplaceAll(strings.ReplaceAll(path, "<b>", ""), "</b>", "") // remove <b> tags present in aviaplanner chart data
	sanitized = strings.ReplaceAll(strings.ReplaceAll(sanitized, "/", ""), `\`, "")  // remove all slashes to avoid trying to open nonexistent directories
	return sanitized
}
