package command

import (
	"fmt"
	"os"
	"port-scanner/internal/output"
	"port-scanner/internal/scanner"
	"port-scanner/internal/types"
	"port-scanner/internal/utils"
	"strings"

	"github.com/spf13/cobra"
)

var (
	cfg     types.Config
	rootCmd = &cobra.Command{
		Use:     "port-scanner",
		Short:   "Port Scanner",
		Example: getExamples(),
		Version: "1.0.0",
		RunE:    run,
	}
)

func init() {
	rootCmd.Flags().StringVarP(&cfg.Address, "address", "a", "", "domain or ip address")
	rootCmd.Flags().StringVarP(&cfg.Ports, "ports", "p", "1-65535", "range: 1-1024 or list: 80,443")
	rootCmd.Flags().StringVarP(&cfg.Mode, "mode", "m", "default", "stealth, default, rapid")
	rootCmd.Flags().StringVarP(&cfg.Output, "output", "o", "", "output file name")
	rootCmd.Flags().StringVarP(&cfg.Format, "format", "f", "txt", "txt, json, csv")
	rootCmd.Flags().IntVarP(&cfg.Timeout, "timeout", "t", 0, "timeout per port in milliseconds")
	_ = rootCmd.MarkFlagRequired("address")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func getExamples() string {
	if utils.IsDockerized() {
		return strings.Join([]string{
			"docker run --rm -v /path/to/your/output:/output port-scanner -a 192.168.1.134",
			"docker run --rm -v /path/to/your/output:/output port-scanner -a 192.168.1.134 -p 1-1024 -m stealth",
			"docker run --rm -v /path/to/your/output:/output port-scanner -a 192.168.1.134 -p 80,443 -o results -f json",
		}, "\n")
	}

	return strings.Join([]string{
		"port-scanner -a 192.168.1.134",
		"port-scanner -a 192.168.1.134 -p 1-1024 -m stealth",
		"port-scanner -a 192.168.1.134 -p 80,443 -o results -f json",
	}, "\n")
}

func run(_ *cobra.Command, _ []string) error {
	results, err := scanner.Scan(cfg)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	err = output.Export(results, cfg)
	if err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	return nil
}
