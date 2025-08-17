package scanner

import (
	"fmt"
	"strings"
	"time"
)

type Mode string

const (
	ModeStealth = "stealth"
	ModeDefault = "default"
	ModeRapid   = "rapid"
)

type metadata struct {
	workerCount int
	timeout     time.Duration
}

var metadataMap = map[Mode]metadata{
	ModeStealth: {
		workerCount: 10,
		timeout:     5 * time.Second,
	},
	ModeDefault: {
		workerCount: 100,
		timeout:     1 * time.Second,
	},
	ModeRapid: {
		workerCount: 1000,
		timeout:     500 * time.Millisecond,
	},
}

func (m Mode) WorkerCount() int {
	return metadataMap[m].workerCount
}

func (m Mode) Timeout() time.Duration {
	return metadataMap[m].timeout
}

func ParseMode(s string) (Mode, error) {
	switch strings.ToLower(s) {
	case "stealth":
		return ModeStealth, nil
	case "default":
		return ModeDefault, nil
	case "rapid":
		return ModeRapid, nil
	default:
		return "", fmt.Errorf("invalid mode: %q", s)
	}
}
