package main

import (
	"fmt"
	"os"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/alert"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfgPath := "portwatch.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	notifier := alert.NewConsoleNotifier(os.Stdout)

	m := monitor.New(cfg, notifier)

	fmt.Printf("portwatch starting — watching %s on %d port(s)\n",
		cfg.Host, len(cfg.Ports))

	if err := m.Run(); err != nil {
		return fmt.Errorf("monitor exited: %w", err)
	}
	return nil
}
