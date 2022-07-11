package regex

import (
	"fmt"
	"regexp"
)

func CaptureN(r *regexp.Regexp, n int, s string) ([]string, error) {
	m := r.FindAllStringSubmatch(s, -1)
	if len(m) != 1 {
		return nil, fmt.Errorf("expected 1 match, got %d", len(m))
	}
	if len(m[0]) < n {
		return nil, fmt.Errorf("expected at least %d submatches, got %d", n, len(m))
	}
	return m[0][0:n], nil
}

func Capture2(r *regexp.Regexp, s string) (string, string, error) {
	m, err := CaptureN(r, 2, s)
	if err != nil {
		return "", "", err
	}
	return m[0], m[1], nil
}

// Implement Capture3 through Capture8.
func Capture3(r *regexp.Regexp, s string) (string, string, string, error) {
	m, err := CaptureN(r, 3, s)
	if err != nil {
		return "", "", "", err
	}
	return m[0], m[1], m[2], nil
}

func Capture4(r *regexp.Regexp, s string) (string, string, string, string, error) {
	m, err := CaptureN(r, 4, s)
	if err != nil {
		return "", "", "", "", err
	}
	return m[0], m[1], m[2], m[3], nil
}

func Capture5(r *regexp.Regexp, s string) (string, string, string, string, string, error) {
	m, err := CaptureN(r, 5, s)
	if err != nil {
		return "", "", "", "", "", err
	}
	return m[0], m[1], m[2], m[3], m[4], nil
}

func Capture6(r *regexp.Regexp, s string) (string, string, string, string, string, string, error) {
	m, err := CaptureN(r, 6, s)
	if err != nil {
		return "", "", "", "", "", "", err
	}
	return m[0], m[1], m[2], m[3], m[4], m[5], nil
}

func Capture7(r *regexp.Regexp, s string) (string, string, string, string, string, string, string, error) {
	m, err := CaptureN(r, 7, s)
	if err != nil {
		return "", "", "", "", "", "", "", err
	}
	return m[0], m[1], m[2], m[3], m[4], m[5], m[6], nil
}

func Capture8(r *regexp.Regexp, s string) (string, string, string, string, string, string, string, string, error) {
	m, err := CaptureN(r, 7, s)
	if err != nil {
		return "", "", "", "", "", "", "", "", err
	}
	return m[0], m[1], m[2], m[3], m[4], m[5], m[6], m[7], nil
}
