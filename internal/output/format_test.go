package output

import (
	"testing"
)

func TestFormat_Extension(t *testing.T) {
	tests := []struct {
		name     string
		format   Format
		expected string
	}{
		{"CSV format", FormatCsv, ".csv"},
		{"JSON format", FormatJson, ".json"},
		{"TXT format", FormatTxt, ".txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.format.Extension()
			if got != tt.expected {
				t.Errorf("Extension() = %q; want %q", got, tt.expected)
			}
		})
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input       string
		expected    Format
		expectError bool
	}{
		{"csv", FormatCsv, false},
		{"CSV", FormatCsv, false},
		{"json", FormatJson, false},
		{"JSON", FormatJson, false},
		{"txt", FormatTxt, false},
		{"TXT", FormatTxt, false},
		{"xml", "", true},
		{"", "", true},
		{"unknown", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("ParseFormat(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseFormat(%q) unexpected error: %v", tt.input, err)
				}
				if got != tt.expected {
					t.Errorf("ParseFormat(%q) = %v; want %v", tt.input, got, tt.expected)
				}
			}
		})
	}
}
