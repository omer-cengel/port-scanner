package scanner

import (
	"errors"
	"fmt"
	"net"
	"port-scanner/internal/types"
	"sync"
	"testing"
	"time"

	"github.com/vbauerster/mpb"
)

func TestScan(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create test listener: %v", err)
	}
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)
	openPort := listener.Addr().(*net.TCPAddr).Port

	tests := []struct {
		name      string
		config    types.Config
		expectErr bool
	}{
		{
			name: "valid config with single port",
			config: types.Config{
				Address: "127.0.0.1",
				Ports:   fmt.Sprintf("%d", openPort),
				Mode:    "default",
				Timeout: 0,
			},
			expectErr: false,
		},
		{
			name: "valid config with port range",
			config: types.Config{
				Address: "127.0.0.1",
				Ports:   fmt.Sprintf("%d-%d", openPort, openPort),
				Mode:    "default",
				Timeout: 100,
			},
			expectErr: false,
		},
		{
			name: "invalid port format",
			config: types.Config{
				Address: "127.0.0.1",
				Ports:   "abc",
				Mode:    "default",
				Timeout: 0,
			},
			expectErr: true,
		},
		{
			name: "invalid mode falls back to default",
			config: types.Config{
				Address: "127.0.0.1",
				Ports:   fmt.Sprintf("%d", openPort),
				Mode:    "invalid",
				Timeout: 0,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Scan(tt.config)

			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}

			if results == nil {
				t.Error("Expected results but got nil")
				return
			}

			if len(results) == 0 {
				t.Error("Expected at least one result")
			}
		})
	}
}

func TestParsePorts(t *testing.T) {
	tests := []struct {
		name      string
		ports     string
		expected  []int
		expectErr error
	}{
		{
			name:      "valid port range",
			ports:     "80-82",
			expected:  []int{80, 81, 82},
			expectErr: nil,
		},
		{
			name:      "single port range",
			ports:     "80-80",
			expected:  []int{80},
			expectErr: nil,
		},
		{
			name:      "invalid port range",
			ports:     "82-80",
			expected:  nil,
			expectErr: invalidPortRangeError,
		},
		{
			name:      "single port",
			ports:     "80",
			expected:  []int{80},
			expectErr: nil,
		},
		{
			name:      "multiple ports",
			ports:     "80,443,22",
			expected:  []int{80, 443, 22},
			expectErr: nil,
		},
		{
			name:      "invalid port in list",
			ports:     "80,abc,22",
			expected:  nil,
			expectErr: invalidPortFormatError,
		},
		{
			name:      "empty string",
			ports:     "",
			expected:  nil,
			expectErr: invalidPortFormatError,
		},
		{
			name:      "port list with dash in middle",
			ports:     "80,443-445,22",
			expected:  []int{80, 443, 445, 22},
			expectErr: invalidPortFormatError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parsePorts(tt.ports)

			if tt.expectErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectErr)
				} else if !errors.Is(err, tt.expectErr) {
					t.Errorf("Expected error %v, got %v", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d ports, got %d", len(tt.expected), len(result))
				return
			}

			for i, expectedPort := range tt.expected {
				if result[i] != expectedPort {
					t.Errorf("Expected port[%d] = %d, got %d", i, expectedPort, result[i])
				}
			}
		})
	}
}

func TestParsePortRange(t *testing.T) {
	tests := []struct {
		name      string
		portRange string
		expected  []int
		expectErr error
	}{
		{
			name:      "valid small range",
			portRange: "80-82",
			expected:  []int{80, 81, 82},
			expectErr: nil,
		},
		{
			name:      "single port range",
			portRange: "80-80",
			expected:  []int{80},
			expectErr: nil,
		},
		{
			name:      "range with spaces",
			portRange: "80 - 82",
			expected:  []int{80, 81, 82},
			expectErr: nil,
		},
		{
			name:      "valid edge range",
			portRange: "1-3",
			expected:  []int{1, 2, 3},
			expectErr: nil,
		},
		{
			name:      "start greater than end",
			portRange: "82-80",
			expected:  nil,
			expectErr: invalidPortRangeError,
		},
		{
			name:      "start below minimum",
			portRange: "0-80",
			expected:  nil,
			expectErr: invalidPortRangeError,
		},
		{
			name:      "end above maximum",
			portRange: "80-65536",
			expected:  nil,
			expectErr: invalidPortRangeError,
		},
		{
			name:      "invalid start format",
			portRange: "abc-80",
			expected:  nil,
			expectErr: invalidPortFormatError,
		},
		{
			name:      "invalid end format",
			portRange: "80-xyz",
			expected:  nil,
			expectErr: invalidPortFormatError,
		},
		{
			name:      "empty string",
			portRange: "",
			expected:  nil,
			expectErr: invalidPortFormatError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parsePortRange(tt.portRange)

			if tt.expectErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectErr)
				} else if !errors.Is(err, tt.expectErr) {
					t.Errorf("Expected error %v, got %v", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d ports, got %d", len(tt.expected), len(result))
				return
			}

			for i, expectedPort := range tt.expected {
				if result[i] != expectedPort {
					t.Errorf("Expected port[%d] = %d, got %d", i, expectedPort, result[i])
				}
			}
		})
	}
}

func TestParsePortList(t *testing.T) {
	tests := []struct {
		name      string
		portList  string
		expected  []int
		expectErr error
	}{
		{
			name:      "single port",
			portList:  "80",
			expected:  []int{80},
			expectErr: nil,
		},
		{
			name:      "multiple ports",
			portList:  "80,443,22",
			expected:  []int{80, 443, 22},
			expectErr: nil,
		},
		{
			name:      "ports with spaces",
			portList:  "80, 443 , 22",
			expected:  []int{80, 443, 22},
			expectErr: nil,
		},
		{
			name:      "invalid port format",
			portList:  "80,abc,22",
			expected:  nil,
			expectErr: invalidPortFormatError,
		},
		{
			name:      "port below minimum",
			portList:  "80,0,22",
			expected:  nil,
			expectErr: invalidPortRangeError,
		},
		{
			name:      "port above maximum",
			portList:  "80,65536,22",
			expected:  nil,
			expectErr: invalidPortRangeError,
		},
		{
			name:      "empty string",
			portList:  "",
			expected:  nil,
			expectErr: invalidPortFormatError,
		},
		{
			name:      "valid edge ports",
			portList:  "1,65535",
			expected:  []int{1, 65535},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parsePortList(tt.portList)

			if tt.expectErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectErr)
				} else if !errors.Is(err, tt.expectErr) {
					t.Errorf("Expected error %v, got %v", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d ports, got %d", len(tt.expected), len(result))
				return
			}

			for i, expectedPort := range tt.expected {
				if result[i] != expectedPort {
					t.Errorf("Expected port[%d] = %d, got %d", i, expectedPort, result[i])
				}
			}
		})
	}
}

func TestScanPorts(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create test listener: %v", err)
	}
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)
	openPort := listener.Addr().(*net.TCPAddr).Port

	tests := []struct {
		name        string
		portList    []int
		workerCount int
		expected    []bool
	}{
		{
			name:        "single open port",
			portList:    []int{openPort},
			workerCount: 1,
			expected:    []bool{true},
		},
		{
			name:        "mixed ports",
			portList:    []int{openPort, 99999},
			workerCount: 2,
			expected:    []bool{true, false},
		},
		{
			name:        "multiple closed ports",
			portList:    []int{99999, 99998},
			workerCount: 1,
			expected:    []bool{false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := scanPorts("127.0.0.1", tt.portList, time.Millisecond*100, tt.workerCount)

			if len(results) != len(tt.portList) {
				t.Errorf("Expected %d results, got %d", len(tt.portList), len(results))
			}

			for i, expected := range tt.expected {
				if i < len(results) {
					if results[i].Status != expected {
						t.Errorf("Result[%d].Status = %v, want %v", i, results[i].Status, expected)
					}
					if results[i].Port != tt.portList[i] {
						t.Errorf("Result[%d].Port = %d, want %d", i, results[i].Port, tt.portList[i])
					}
				}
			}
		})
	}
}

func TestCreateScanTasks(t *testing.T) {
	tests := []struct {
		name     string
		portList []int
	}{
		{
			name:     "empty port list",
			portList: []int{},
		},
		{
			name:     "single port",
			portList: []int{80},
		},
		{
			name:     "multiple ports",
			portList: []int{80, 443, 22, 21},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks := createScanTasks(tt.portList)

			if tasks == nil {
				t.Error("Tasks channel should not be nil")
			}

			var receivedTasks []types.Task
			for task := range tasks {
				receivedTasks = append(receivedTasks, task)
			}

			if len(receivedTasks) != len(tt.portList) {
				t.Errorf("Expected %d tasks, got %d", len(tt.portList), len(receivedTasks))
			}

			for i, expectedPort := range tt.portList {
				if i < len(receivedTasks) {
					if receivedTasks[i].Index != i {
						t.Errorf("Task[%d].Index = %d, want %d", i, receivedTasks[i].Index, i)
					}
					if receivedTasks[i].Port != expectedPort {
						t.Errorf("Task[%d].Port = %d, want %d", i, receivedTasks[i].Port, expectedPort)
					}
				}
			}
		})
	}
}

func TestBuildProgressBar(t *testing.T) {
	tests := []struct {
		name     string
		portList []int
		expected int64
	}{
		{
			name:     "empty port list",
			portList: []int{},
			expected: 0,
		},
		{
			name:     "single port",
			portList: []int{80},
			expected: 1,
		},
		{
			name:     "multiple ports",
			portList: []int{80, 443, 22, 21, 25},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, b := buildProgressBar(tt.portList)

			if p == nil {
				t.Error("Progress should not be nil")
			}

			if b == nil {
				t.Error("Bar should not be nil")
			}
		})
	}
}

func TestStartScanWorkers(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create test listener: %v", err)
	}
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)
	openPort := listener.Addr().(*net.TCPAddr).Port

	testTasks := []types.Task{
		{Port: openPort, Index: 0},
		{Port: 99999, Index: 1},
		{Port: openPort, Index: 2},
		{Port: 99998, Index: 3},
	}

	tasks := make(chan types.Task, len(testTasks))
	results := make([]types.Result, len(testTasks))
	var wg sync.WaitGroup

	p := mpb.New()
	bar := p.AddBar(int64(len(testTasks)))

	workerCount := 3
	startScanWorkers(tasks, results, "127.0.0.1", time.Millisecond*100, workerCount, bar, &wg)

	for _, task := range testTasks {
		tasks <- task
	}
	close(tasks)

	wg.Wait()
	p.Wait()

	expectedStatuses := []bool{true, false, true, false}
	for i, expected := range expectedStatuses {
		if results[i].Status != expected {
			t.Errorf("Result[%d].Status = %v, want %v", i, results[i].Status, expected)
		}
		if results[i].Port != testTasks[i].Port {
			t.Errorf("Result[%d].Port = %d, want %d", i, results[i].Port, testTasks[i].Port)
		}
	}
}

func TestRunScanWorker(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create test listener: %v", err)
	}
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)
	openPort := listener.Addr().(*net.TCPAddr).Port

	testTasks := []types.Task{
		{Port: openPort, Index: 0},
		{Port: 99999, Index: 1},
		{Port: openPort, Index: 2},
	}

	tasks := make(chan types.Task, len(testTasks))
	results := make([]types.Result, len(testTasks))

	p := mpb.New()
	bar := p.AddBar(int64(len(testTasks)))

	for _, task := range testTasks {
		tasks <- task
	}
	close(tasks)

	runScanWorker(tasks, results, "127.0.0.1", time.Millisecond*100, bar)

	p.Wait()

	expectedResults := []struct {
		port   int
		status bool
	}{
		{openPort, true},
		{99999, false},
		{openPort, true},
	}

	for i, expected := range expectedResults {
		if results[i].Port != expected.port {
			t.Errorf("Result[%d].Port = %d, want %d", i, results[i].Port, expected.port)
		}
		if results[i].Status != expected.status {
			t.Errorf("Result[%d].Status = %v, want %v", i, results[i].Status, expected.status)
		}
	}
}

func TestScanPort(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     int
		timeout  time.Duration
		expected bool
	}{
		{
			name:     "open port",
			host:     "127.0.0.1",
			port:     0,
			timeout:  time.Second,
			expected: true,
		},
		{
			name:     "closed port",
			host:     "127.0.0.1",
			port:     99999,
			timeout:  time.Millisecond * 100,
			expected: false,
		},
		{
			name:     "invalid host",
			host:     "invalid.host",
			port:     80,
			timeout:  time.Millisecond * 100,
			expected: false,
		},
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create test listener: %v", err)
	}
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	tests[0].port = listener.Addr().(*net.TCPAddr).Port

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanPort(tt.host, tt.port, tt.timeout)
			if result != tt.expected {
				t.Errorf("scanPort() = %v, want %v", result, tt.expected)
			}
		})
	}
}
