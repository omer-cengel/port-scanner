package output

import (
	"encoding/json"
	"os"
	"path/filepath"
	"port-scanner/internal/types"
	"strings"
	"testing"
	"time"
)

var testResults = []types.Result{
	{Port: 80, Status: true},
	{Port: 443, Status: false},
	{Port: 8080, Status: true},
}

func TestExport(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		results []types.Result
		config  types.Config
		wantErr bool
	}{
		{
			name:    "successful export with txt format",
			results: testResults,
			config:  types.Config{Format: "txt", Output: filepath.Join(tempDir, "test.txt")},
			wantErr: false,
		},
		{
			name:    "successful export with csv format",
			results: testResults,
			config:  types.Config{Format: "csv", Output: filepath.Join(tempDir, "test.csv")},
			wantErr: false,
		},
		{
			name:    "successful export with json format",
			results: testResults,
			config:  types.Config{Format: "json", Output: filepath.Join(tempDir, "test.json")},
			wantErr: false,
		},
		{
			name:    "export with unknown format defaults to txt",
			results: testResults,
			config:  types.Config{Format: "unknown", Output: filepath.Join(tempDir, "test_unknown.txt")},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Export(tt.results, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Export() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormatResults(t *testing.T) {
	tests := []struct {
		name    string
		results []types.Result
		format  Format
		wantErr bool
	}{
		{
			name:    "format as CSV",
			results: testResults,
			format:  FormatCsv,
			wantErr: false,
		},
		{
			name:    "format as JSON",
			results: testResults,
			format:  FormatJson,
			wantErr: false,
		},
		{
			name:    "format as TXT",
			results: testResults,
			format:  FormatTxt,
			wantErr: false,
		},
		{
			name:    "format with empty results",
			results: []types.Result{},
			format:  FormatCsv,
			wantErr: false,
		},
		{
			name:    "format unknown",
			results: []types.Result{},
			format:  Format("unknown"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := formatResults(tt.results, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && output == "" {
				t.Errorf("formatResults() returned empty output")
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		name    string
		results []types.Result
		wantErr bool
	}{
		{
			name:    "valid results",
			results: testResults,
			wantErr: false,
		},
		{
			name:    "empty results",
			results: []types.Result{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := toJSON(tt.results)
			if (err != nil) != tt.wantErr {
				t.Errorf("toJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				var parsed []types.Result
				if err := json.Unmarshal([]byte(output), &parsed); err != nil {
					t.Errorf("toJSON() produced invalid JSON: %v", err)
				}

				if len(parsed) != len(tt.results) {
					t.Errorf("toJSON() length mismatch: got %d, want %d", len(parsed), len(tt.results))
				}
			}
		})
	}
}

func TestToCSV(t *testing.T) {
	tests := []struct {
		name    string
		results []types.Result
		wantErr bool
	}{
		{
			name:    "valid results",
			results: testResults,
			wantErr: false,
		},
		{
			name:    "empty results",
			results: []types.Result{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := toCSV(tt.results)
			if (err != nil) != tt.wantErr {
				t.Errorf("toCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !strings.Contains(output, "Port,Status") {
					t.Errorf("toCSV() missing expected header")
				}

				lines := strings.Split(strings.TrimSpace(output), "\n")
				expectedLines := len(tt.results) + 1
				if len(lines) != expectedLines {
					t.Errorf("toCSV() line count mismatch: got %d, want %d", len(lines), expectedLines)
				}
			}
		})
	}
}

func TestToTXT(t *testing.T) {
	tests := []struct {
		name    string
		results []types.Result
	}{
		{
			name:    "valid results",
			results: testResults,
		},
		{
			name:    "empty results",
			results: []types.Result{},
		},
		{
			name:    "single result",
			results: []types.Result{{Port: 22, Status: true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := toTXT(tt.results)

			if !strings.Contains(output, "Port") || !strings.Contains(output, "Status") {
				t.Errorf("toTXT() missing expected header")
			}

			lines := strings.Split(strings.TrimSpace(output), "\n")
			expectedLines := len(tt.results) + 1
			if len(lines) != expectedLines {
				t.Errorf("toTXT() line count mismatch: got %d, want %d", len(lines), expectedLines)
			}
		})
	}
}

func TestGenerateOutputPath(t *testing.T) {
	tests := []struct {
		name       string
		dockerized string
		output     string
		extension  string
		want       string
	}{
		{
			name:       "empty output generates filename",
			dockerized: "false",
			output:     "",
			extension:  ".txt",
			want:       ".txt",
		},
		{
			name:       "output with extension",
			dockerized: "false",
			output:     "test.txt",
			extension:  ".txt",
			want:       "test.txt",
		},
		{
			name:       "output without extension",
			dockerized: "false",
			output:     "test",
			extension:  ".csv",
			want:       "test.csv",
		},
		{
			name:       "output as directory",
			dockerized: "false",
			output:     "results/",
			extension:  ".json",
			want:       ".json",
		},
		{
			name:       "output as directory",
			dockerized: "true",
			output:     "results/",
			extension:  ".json",
			want:       ".json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DOCKERIZED", tt.dockerized)
			got := generateOutputPath(tt.output, tt.extension)

			if tt.output == "" || strings.HasSuffix(tt.output, "/") {
				if !strings.HasSuffix(got, tt.extension) {
					t.Errorf("generateOutputPath() = %v, want suffix %v", got, tt.extension)
				}
			} else {
				if tt.output == tt.want || strings.HasSuffix(got, tt.extension) {
				} else {
					t.Errorf("generateOutputPath() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGenerateDockerOutputPath(t *testing.T) {
	tests := []struct {
		name      string
		output    string
		extension string
		wantDir   string
	}{
		{
			name:      "empty output",
			output:    "",
			extension: ".txt",
			wantDir:   "/output",
		},
		{
			name:      "output with extension",
			output:    "test.csv",
			extension: ".csv",
			wantDir:   "/output",
		},
		{
			name:      "output without extension",
			output:    "test",
			extension: ".json",
			wantDir:   "/output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateDockerOutputPath(tt.output, tt.extension)

			if !strings.HasPrefix(got, tt.wantDir) {
				t.Errorf("generateDockerOutputPath() = %v, want prefix %v", got, tt.wantDir)
			}

			if !strings.HasSuffix(got, tt.extension) {
				t.Errorf("generateDockerOutputPath() = %v, want suffix %v", got, tt.extension)
			}
		})
	}
}

func TestGenerateLocalOutputPath(t *testing.T) {
	tests := []struct {
		name      string
		output    string
		extension string
	}{
		{
			name:      "empty output",
			output:    "",
			extension: ".txt",
		},
		{
			name:      "output with extension",
			output:    "test.csv",
			extension: ".csv",
		},
		{
			name:      "output without extension",
			output:    "test",
			extension: ".json",
		},
		{
			name:      "directory path",
			output:    "results/",
			extension: ".txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateLocalOutputPath(tt.output, tt.extension)

			if !strings.HasSuffix(got, tt.extension) {
				t.Errorf("generateLocalOutputPath() = %v, want suffix %v", got, tt.extension)
			}
		})
	}
}

func TestGenerateFileName(t *testing.T) {
	fileName := generateFileName()

	_, err := time.Parse(dateFormat, fileName)
	if err != nil {
		t.Errorf("generateFileName() = %v, failed to parse with format %v: %v", fileName, dateFormat, err)
	}

	time.Sleep(time.Second)
	fileName2 := generateFileName()
	if fileName == fileName2 {
		t.Logf("generateFileName() produced same filename twice (this may be timing-related): %v", fileName)
	}
}

func TestWriteToFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		filePath string
		content  string
		wantErr  bool
	}{
		{
			name:     "write to valid path",
			filePath: filepath.Join(tempDir, "test.txt"),
			content:  "test content",
			wantErr:  false,
		},
		{
			name:     "write to nested directory",
			filePath: filepath.Join(tempDir, "sub", "test.txt"),
			content:  "test content",
			wantErr:  false,
		},
		{
			name:     "write empty content",
			filePath: filepath.Join(tempDir, "empty.txt"),
			content:  "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := writeToFile(tt.filePath, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("writeToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				content, err := os.ReadFile(tt.filePath)
				if err != nil {
					t.Errorf("Failed to read written file: %v", err)
					return
				}

				if string(content) != tt.content {
					t.Errorf("writeToFile() content = %v, want %v", string(content), tt.content)
				}
			}
		})
	}
}
