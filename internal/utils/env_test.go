package utils

import (
	"testing"
)

func TestIsDockerized(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		setEnv   bool
		expected bool
	}{
		{
			name:     "environment variable set to 'true'",
			envValue: "true",
			setEnv:   true,
			expected: true,
		},
		{
			name:     "environment variable set to 'TRUE'",
			envValue: "TRUE",
			setEnv:   true,
			expected: true,
		},
		{
			name:     "environment variable set to 'True'",
			envValue: "True",
			setEnv:   true,
			expected: true,
		},
		{
			name:     "environment variable set to '1'",
			envValue: "1",
			setEnv:   true,
			expected: true,
		},
		{
			name:     "environment variable set to 'false'",
			envValue: "false",
			setEnv:   true,
			expected: false,
		},
		{
			name:     "environment variable set to 'FALSE'",
			envValue: "FALSE",
			setEnv:   true,
			expected: false,
		},
		{
			name:     "environment variable set to 'False'",
			envValue: "False",
			setEnv:   true,
			expected: false,
		},
		{
			name:     "environment variable set to '0'",
			envValue: "0",
			setEnv:   true,
			expected: false,
		},
		{
			name:     "environment variable set to invalid value",
			envValue: "invalid",
			setEnv:   true,
			expected: false,
		},
		{
			name:     "environment variable set to empty string",
			envValue: "",
			setEnv:   true,
			expected: false,
		},
		{
			name:     "environment variable set to whitespace",
			envValue: "   ",
			setEnv:   true,
			expected: false,
		},
		{
			name:     "environment variable not set",
			setEnv:   false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv(envDockerized, tt.envValue)
			}

			result := IsDockerized()
			if result != tt.expected {
				t.Errorf("IsDockerized() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
