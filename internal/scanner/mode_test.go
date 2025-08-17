package scanner

import (
	"testing"
	"time"
)

func TestParseMode(t *testing.T) {
	tests := []struct {
		input       string
		wantMode    Mode
		expectError bool
	}{
		{"stealth", ModeStealth, false},
		{"default", ModeDefault, false},
		{"rapid", ModeRapid, false},
		{"STEALTH", ModeStealth, false}, // case-insensitive
		{"invalid", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseMode(tt.input)
			if (err != nil) != tt.expectError {
				t.Fatalf("ParseMode(%q) error = %v, wantErr %v", tt.input, err, tt.expectError)
			}
			if got != tt.wantMode {
				t.Errorf("ParseMode(%q) = %v, want %v", tt.input, got, tt.wantMode)
			}
		})
	}
}

func TestModeWorkerCount(t *testing.T) {
	tests := []struct {
		mode     Mode
		expected int
	}{
		{ModeStealth, 10},
		{ModeDefault, 100},
		{ModeRapid, 1000},
	}

	for _, tt := range tests {
		t.Run(string(tt.mode), func(t *testing.T) {
			if got := tt.mode.WorkerCount(); got != tt.expected {
				t.Errorf("%v.WorkerCount() = %d; want %d", tt.mode, got, tt.expected)
			}
		})
	}
}

func TestModeTimeout(t *testing.T) {
	tests := []struct {
		mode     Mode
		expected time.Duration
	}{
		{ModeStealth, 5 * time.Second},
		{ModeDefault, 1 * time.Second},
		{ModeRapid, 500 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(string(tt.mode), func(t *testing.T) {
			if got := tt.mode.Timeout(); got != tt.expected {
				t.Errorf("%v.Timeout() = %v; want %v", tt.mode, got, tt.expected)
			}
		})
	}
}
