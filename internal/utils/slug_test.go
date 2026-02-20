package utils

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"  Leading and Trailing  ", "leading-and-trailing"},
		{"UPPERCASE STRING", "uppercase-string"},
		{"special!@#$%chars", "specialchars"},
		{"multiple   spaces", "multiple-spaces"},
		{"already-slugified", "already-slugified"},
		{"Mixed CASE with-Dashes", "mixed-case-with-dashes"},
		{"", ""},
		{"   ", ""},
		{"Product Name (v2)", "product-name-v2"},
		{"café-latte", "caf-latte"},
		{"one--two---three", "one-two-three"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Slugify(tt.input)
			if result != tt.expected {
				t.Fatalf("Slugify(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
