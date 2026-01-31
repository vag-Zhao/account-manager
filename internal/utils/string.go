package utils

import (
	"strings"
	"unicode"
)

// TrimSpace removes leading and trailing whitespace
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// IsEmpty checks if a string is empty or contains only whitespace
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty checks if a string is not empty
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// Truncate truncates a string to the specified length
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Contains checks if a string contains a substring (case-insensitive)
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// EqualsIgnoreCase checks if two strings are equal (case-insensitive)
func EqualsIgnoreCase(s1, s2 string) bool {
	return strings.EqualFold(s1, s2)
}

// SanitizeString removes non-printable characters from a string
func SanitizeString(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, s)
}

// FirstN returns the first n characters of a string
func FirstN(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// LastN returns the last n characters of a string
func LastN(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

// ReplaceMultiple replaces multiple old strings with new strings
func ReplaceMultiple(s string, replacements map[string]string) string {
	result := s
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	return result
}

// SplitAndTrim splits a string by delimiter and trims each part
func SplitAndTrim(s, delimiter string) []string {
	parts := strings.Split(s, delimiter)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// JoinNonEmpty joins non-empty strings with a delimiter
func JoinNonEmpty(delimiter string, parts ...string) string {
	nonEmpty := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}
	return strings.Join(nonEmpty, delimiter)
}

// MaskString masks part of a string (useful for passwords, emails)
func MaskString(s string, visibleStart, visibleEnd int) string {
	if len(s) <= visibleStart+visibleEnd {
		return s
	}
	start := s[:visibleStart]
	end := s[len(s)-visibleEnd:]
	masked := strings.Repeat("*", len(s)-visibleStart-visibleEnd)
	return start + masked + end
}

// DefaultIfEmpty returns the default value if the string is empty
func DefaultIfEmpty(s, defaultValue string) string {
	if IsEmpty(s) {
		return defaultValue
	}
	return s
}
