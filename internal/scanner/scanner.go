package scanner

import (
	"errors"
	"fmt"
	"net"
	"port-scanner/internal/types"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

const (
	minPortNumber      = 1
	maxPortNumber      = 65535
	portRangeDelimiter = "-"
	portListSeparator  = ","
	networkTCP         = "tcp"
	addressFormat      = "%s:%d"
)

var (
	invalidPortFormatError = errors.New("invalid port format: expected range: '1-1024' or list: '80,443'")
	invalidPortRangeError  = errors.New("invalid port range: expected range between 1 and 65535")
)

func Scan(cfg types.Config) ([]types.Result, error) {
	ports, err := parsePorts(cfg.Ports)
	if err != nil {
		return nil, err
	}

	mode, err := ParseMode(cfg.Mode)
	if err != nil {
		mode = ModeDefault
	}

	timeout := mode.Timeout()
	if cfg.Timeout > 0 {
		timeout = time.Duration(cfg.Timeout) * time.Millisecond
	}

	results := scanPorts(cfg.Address, ports, timeout, mode.WorkerCount())

	return results, nil
}

func parsePorts(ports string) ([]int, error) {
	if strings.Contains(ports, portRangeDelimiter) {
		return parsePortRange(ports)
	} else {
		return parsePortList(ports)
	}
}

func parsePortRange(portRange string) ([]int, error) {
	parts := strings.SplitN(portRange, portRangeDelimiter, 2)

	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, invalidPortFormatError
	}

	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, invalidPortFormatError
	}

	if start > end || start < minPortNumber || end > maxPortNumber {
		return nil, invalidPortRangeError
	}

	ports := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		ports = append(ports, i)
	}

	return ports, nil
}

func parsePortList(portList string) ([]int, error) {
	parts := strings.Split(portList, portListSeparator)
	ports := make([]int, 0, len(parts))
	seen := make(map[int]bool)

	for _, part := range parts {
		port, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, invalidPortFormatError
		}

		if port < minPortNumber || port > maxPortNumber {
			return nil, invalidPortRangeError
		}

		if !seen[port] {
			seen[port] = true
			ports = append(ports, port)
		}
	}

	return ports, nil
}

func scanPorts(host string, portList []int, timeout time.Duration, workerCount int) []types.Result {
	tasks := createScanTasks(portList)
	results := make([]types.Result, len(portList))
	progress, bar := buildProgressBar(portList)

	var wg sync.WaitGroup
	startScanWorkers(tasks, results, host, timeout, workerCount, bar, &wg)

	wg.Wait()
	progress.Wait()
	return results
}

func createScanTasks(portList []int) chan types.Task {
	tasks := make(chan types.Task, len(portList))
	for i, port := range portList {
		tasks <- types.Task{Index: i, Port: port}
	}
	close(tasks)
	return tasks
}

func buildProgressBar(portList []int) (*mpb.Progress, *mpb.Bar) {
	p := mpb.New(mpb.WithWidth(60))
	b := p.AddBar(int64(len(portList)),
		mpb.PrependDecorators(
			decor.Name("Scanning "),
			decor.CountersNoUnit("%d / %d"),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)
	return p, b
}

func startScanWorkers(
	tasks chan types.Task,
	results []types.Result,
	host string,
	timeout time.Duration,
	workerCount int,
	bar *mpb.Bar,
	wg *sync.WaitGroup,
) {
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runScanWorker(tasks, results, host, timeout, bar)
		}()
	}
}

func runScanWorker(
	tasks chan types.Task,
	results []types.Result,
	host string,
	timeout time.Duration,
	bar *mpb.Bar,
) {
	for task := range tasks {
		open := scanPort(host, task.Port, timeout)
		results[task.Index] = types.Result{Port: task.Port, Status: open}
		bar.Increment()
	}
}

func scanPort(host string, port int, timeout time.Duration) bool {
	address := fmt.Sprintf(addressFormat, host, port)
	conn, err := net.DialTimeout(networkTCP, address, timeout)
	if err != nil {
		return false
	}
	defer func() {
		_ = conn.Close()
	}()
	return true
}
