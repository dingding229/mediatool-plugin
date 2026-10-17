package contracts

import (
	"strings"
	"unicode/utf8"
)

// NameHint is an untrusted recognition hint, never a media identity or file path.
type NameHint struct {
	Title     string `json:"title"`
	Year      string `json:"year"`
	MediaType string `json:"mediaType"`
	Season    int    `json:"season"`
	Episodes  []int  `json:"episodes"`
}
type NameRecognition struct {
	Name     string   `json:"name"`
	Hint     NameHint `json:"hint"`
	Continue bool     `json:"continue"`
}

func ValidateName(v NameRecognition) error {
	if v.Name == "" || len(v.Name) > 2048 || strings.ContainsAny(v.Name, "/\\\x00\r\n") || !utf8.ValidString(v.Name) || v.Hint.Title == "" || len(v.Hint.Title) > 512 || !utf8.ValidString(v.Hint.Title) || strings.ContainsAny(v.Hint.Title, "\x00\r\n") || v.Hint.Season < 0 || v.Hint.Season > 999 || len(v.Hint.Episodes) > 100 {
		return NewError(InvalidExtensionResult)
	}
	if v.Hint.MediaType != "movie" && v.Hint.MediaType != "series" && v.Hint.MediaType != "unknown" {
		return NewError(InvalidExtensionResult)
	}
	if len(v.Hint.Year) != 0 {
		if len(v.Hint.Year) != 4 {
			return NewError(InvalidExtensionResult)
		}
		for _, c := range v.Hint.Year {
			if c < '0' || c > '9' {
				return NewError(InvalidExtensionResult)
			}
		}
	}
	for _, n := range v.Hint.Episodes {
		if n < 1 || n > 9999 {
			return NewError(InvalidExtensionResult)
		}
	}
	return nil
}
