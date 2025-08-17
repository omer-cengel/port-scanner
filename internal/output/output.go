package output

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"port-scanner/internal/types"
	"port-scanner/internal/utils"
	"strings"
	"time"
)

const (
	headerPort          = "Port"
	headerStatus        = "Status"
	dateFormat          = "2006-01-02_15:04:05"
	outputDirectory     = "/output"
	directoryPermission = 0755
	filePermission      = 0644
)

var (
	writeFileError = errors.New("failed to write file")
)

func Export(results []types.Result, cfg types.Config) error {
	format, err := ParseFormat(cfg.Format)
	if err != nil {
		format = FormatTxt
	}

	output, err := formatResults(results, format)
	if err != nil {
		return err
	}

	outputPath := generateOutputPath(cfg.Output, format.Extension())
	fmt.Println(outputPath)
	err = writeToFile(outputPath, output)
	if err != nil {
		return err
	}

	return nil
}

func formatResults(results []types.Result, format Format) (string, error) {
	switch format {
	case FormatCsv:
		return toCSV(results)
	case FormatJson:
		return toJSON(results)
	case FormatTxt:
		return toTXT(results), nil
	default:
		return toTXT(results), nil
	}
}

func toJSON(results []types.Result) (string, error) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", writeFileError
	}
	return string(data), nil
}

func toCSV(results []types.Result) (string, error) {
	var sb strings.Builder
	writer := csv.NewWriter(&sb)

	err := writer.Write([]string{headerPort, headerStatus})
	if err != nil {
		return "", writeFileError
	}

	for _, r := range results {
		err = writer.Write([]string{
			fmt.Sprintf("%d", r.Port),
			fmt.Sprintf("%t", r.Status),
		})
		if err != nil {
			return "", writeFileError
		}
	}

	writer.Flush()
	err = writer.Error()
	if err != nil {
		return "", writeFileError
	}

	content := sb.String()
	return strings.TrimSuffix(content, "\n"), nil
}

func toTXT(results []types.Result) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-6s %-6s\n", headerPort, headerStatus))

	for _, result := range results {
		sb.WriteString(fmt.Sprintf("%-6d %-6t\n", result.Port, result.Status))
	}

	return sb.String()
}

func generateOutputPath(output, extension string) string {
	if utils.IsDockerized() {
		return generateDockerOutputPath(output, extension)
	}
	return generateLocalOutputPath(output, extension)
}

func generateDockerOutputPath(output, extension string) string {
	if output == "" {
		fileName := generateFileName()
		return filepath.Join(outputDirectory, fileName+extension)
	}

	base := filepath.Base(output)
	if strings.HasSuffix(base, extension) {
		return filepath.Join(outputDirectory, base)
	}

	return filepath.Join(outputDirectory, base+extension)
}

func generateLocalOutputPath(output, extension string) string {
	if output == "" {
		fileName := generateFileName()
		return fileName + extension
	}

	if strings.HasSuffix(output, string(os.PathSeparator)) {
		fileName := generateFileName()
		return filepath.Join(output, fileName+extension)
	}

	base := filepath.Base(output)
	dir := filepath.Dir(output)

	if strings.HasSuffix(base, extension) {
		return filepath.Join(dir, base)
	}

	return filepath.Join(dir, base+extension)
}

func generateFileName() string {
	return time.Now().Format(dateFormat)
}

func writeToFile(filePath, content string) error {
	dir := filepath.Dir(filePath)

	err := os.MkdirAll(dir, directoryPermission)
	if err != nil {
		return writeFileError
	}

	err = os.WriteFile(filePath, []byte(content), filePermission)
	if err != nil {
		return writeFileError
	}

	return nil
}
