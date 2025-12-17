package mapping

import (
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Alias (string -> int cap_bac)
var Alias2Cap = map[string]int{
	// tiếng Việt
	"thuong":      1,
	"hiem":        2,
	"su thi":      3,
	"huyen thoai": 4,
	"thien thoai": 5,

	// tiếng Anh
	"common":    1,
	"rare":      2,
	"epic":      3,
	"legend":    4,
	"legendary": 4,
	"mythic":    5,
}

// Chuẩn hoá tên về dạng “canonical”: bỏ dấu, thường hoá, gọn khoảng trắng, đổi đ/Đ→d
func CanonicalizeName(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	t := transform.Chain(
		norm.NFD,
		transform.RemoveFunc(func(r rune) bool { return unicode.Is(unicode.Mn, r) }),
		norm.NFC,
	)
	out, _, _ := transform.String(t, s)
	out = strings.ToLower(out)
	out = strings.ReplaceAll(out, "đ", "d")
	out = strings.ReplaceAll(out, "Đ", "d")
	out = strings.Join(strings.Fields(out), " ")
	return out
}

// MapStringToCap: nhận chuỗi (có/không dấu) → cấp_bậc
func MapStringToCap(raw string) (int, bool) {
	key := CanonicalizeName(raw)
	if v, ok := Alias2Cap[key]; ok {
		return v, true
	}
	return 0, false
}

//  Parse “cấp bậc” từ text người dùng

type CapBacMode int

const (
	CapBacNone CapBacMode = iota
	CapBacEQ
	CapBacIN
	CapBacRange
)

type CapBacFilter struct {
	Mode   CapBacMode
	Values []int // cho EQ (len=1) & IN
	Min    int
	Max    int
}

// ParseCapBacText: "1,3,5" | "2-5" | "4" -> filter có cấu trúc
func ParseCapBacText(raw string) (CapBacFilter, bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return CapBacFilter{Mode: CapBacNone}, false
	}
	noSpace := strings.ReplaceAll(s, " ", "")

	// Danh sách: "1,2,5"
	if strings.Contains(noSpace, ",") {
		parts := strings.Split(noSpace, ",")
		ints := make([]int, 0, len(parts))
		for _, p := range parts {
			if v, err := strconv.Atoi(p); err == nil {
				ints = append(ints, v)
			}
		}
		if len(ints) > 0 {
			return CapBacFilter{Mode: CapBacIN, Values: ints}, true
		}
		return CapBacFilter{Mode: CapBacNone}, false
	}

	// Khoảng: "a-b"
	if strings.Contains(noSpace, "-") {
		pp := strings.SplitN(noSpace, "-", 2)
		if len(pp) == 2 {
			if l, err1 := strconv.Atoi(pp[0]); err1 == nil {
				if r, err2 := strconv.Atoi(pp[1]); err2 == nil && l <= r {
					return CapBacFilter{Mode: CapBacRange, Min: l, Max: r}, true
				}
			}
		}
		return CapBacFilter{Mode: CapBacNone}, false
	}

	// Số đơn: "2"
	if v, err := strconv.Atoi(noSpace); err == nil {
		return CapBacFilter{Mode: CapBacEQ, Values: []int{v}}, true
	}

	return CapBacFilter{Mode: CapBacNone}, false
}
