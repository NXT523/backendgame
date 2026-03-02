package utils

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// GetenvString trả về giá trị env dạng string.
// Nếu key không tồn tại hoặc rỗng → trả về def.
func GetenvString(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// GetenvStrings trả về slice string từ env dạng CSV: "a,b,c".
// Tự trim space, bỏ phần tử rỗng.
// Nếu kết quả rỗng → trả về def.
func GetenvStrings(key string, def []string) []string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}

	if len(out) == 0 {
		return def
	}
	return out
}

// GetenvInt trả về giá trị biến môi trường dạng int.
// Nếu không tồn tại, rỗng hoặc parse lỗi, trả về def.
func GetenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

// GetenvInt64 trả về giá trị biến môi trường dạng int64.
// Nếu không tồn tại hoặc parse lỗi, trả về def.
func GetenvInt64(key string, def int64) int64 {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return def
}

// GetenvFloat trả về giá trị biến môi trường dạng float64.
// Nếu không tồn tại hoặc parse lỗi, trả về def.
func GetenvFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}

// GetenvBool trả về giá trị biến môi trường dạng bool.
// Chấp nhận: 1, t, T, TRUE, true, True, 0, f, false...
// Nếu không tồn tại hoặc parse lỗi, trả về def.
func GetenvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

// GetenvDuration trả về giá trị biến môi trường dạng time.Duration.
// Ví dụ: 500ms, 2s, 1m.
// Nếu không tồn tại hoặc parse lỗi, trả về def.
func GetenvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
