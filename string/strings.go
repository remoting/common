package utils

import "strings"

func SplitFunc(s string, f func(rune) bool) []string {
	sa := make([]string, 0)
	if len(s) <= 0 {
		return sa
	} else {
		i := strings.IndexFunc(s, f)
		for i != -1 {
			sa = append(sa, s[0:i])
			s = s[i+1:]
			i = strings.IndexFunc(s, f)
		}
		return append(sa, s)
	}
}
func Split(s string) []string {
	return SplitFunc(s, func(r rune) bool {
		switch r {
		case ';', ',':
			return true
		}
		return false
	})
}
func IsEmpty(s string) bool {
	if len(strings.TrimSpace(s)) <= 0 {
		return true
	} else {
		return false
	}
}
