package utils

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)[A-Z][a-z]+")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func CamelToSnake(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

// - Trim space đầu/cuối Gộp nhiều space thành 1
// - Chuyển chữ thường
// Dùng cho: key redis, lock key, so sánh logic
func NormalizeString(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	text = strings.Join(strings.Fields(text), " ")
	return strings.ToLower(text)
}

// - Trim space đầu/cuối. Gộp nhiều space thành 1
// - Không chuyển chữ thường
// Dùng cho: lưu DB, hiển thị
func NormalizeKeepCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	return strings.Join(strings.Fields(s), " ")
}
