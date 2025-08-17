package output

import (
	"fmt"
	"strings"
)

type Format string

const (
	FormatCsv  Format = "csv"
	FormatJson Format = "json"
	FormatTxt  Format = "txt"
)

type metadata struct {
	extension string
}

var metadataMap = map[Format]metadata{
	FormatCsv: {
		extension: ".csv",
	},
	FormatJson: {
		extension: ".json",
	},
	FormatTxt: {
		extension: ".txt",
	},
}

func (f Format) Extension() string {
	return metadataMap[f].extension
}

func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "csv":
		return FormatCsv, nil
	case "json":
		return FormatJson, nil
	case "txt":
		return FormatTxt, nil
	default:
		return "", fmt.Errorf("invalid format: %q", s)
	}
}
