package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	// SQL injection phổ biến
	sqlInjectRe = regexp.MustCompile(`(?i)(union\s+select|select\s+.*from|insert\s+into|drop\s+table|--|;|/\*|\*/|sleep\(|benchmark\(|@@)`)

	// XSS
	xssRe = regexp.MustCompile(`(?i)(<script|javascript:|onerror=|onload=)`)

	// Path traversal
	pathTravRe = regexp.MustCompile(`\.\./|\.\.\\`)
)

// InjectionGuard middleware cho reverse proxy
func InjectionGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1️⃣ Chặn method lạ
		switch r.Method {
		case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
		default:
			http.Error(w, "Phương pháp không được phép", http.StatusMethodNotAllowed)
			return
		}

		// 2️⃣ Path traversal
		if pathTravRe.MatchString(r.URL.Path) {
			http.Error(w, "Đường dẫn không hợp lệ", http.StatusBadRequest)
			return
		}

		// 3️⃣ Query injection
		if hasInjectionInQuery(r.URL) {
			http.Error(w, "Đã phát hiện truy vấn độc hại", http.StatusBadRequest)
			return
		}

		// 4️⃣ Body injection (JSON / form)
		if r.Body != nil && r.ContentLength > 0 {
			bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // max 1MB
			if err != nil {
				http.Error(w, "Nội dung không hợp lệ", http.StatusBadRequest)
				return
			}

			bodyStr := strings.ToLower(string(bodyBytes))
			if sqlInjectRe.MatchString(bodyStr) || xssRe.MatchString(bodyStr) {
				http.Error(w, "Đã phát hiện tải trọng độc hại", http.StatusBadRequest)
				return
			}

			// trả body lại cho reverse proxy
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		next.ServeHTTP(w, r)
	})
}

func hasInjectionInQuery(u *url.URL) bool {
	for _, values := range u.Query() {
		for _, v := range values {
			val := strings.ToLower(v)
			if sqlInjectRe.MatchString(val) || xssRe.MatchString(val) {
				return true
			}
		}
	}
	return false
}
