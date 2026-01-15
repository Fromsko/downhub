package main

import (
	"github.com/Fromsko/downhub/cmd"
	"github.com/Fromsko/downhub/common"
	"github.com/Fromsko/downhub/config"
	"github.com/Fromsko/downhub/handler"
	"github.com/Fromsko/downhub/logs"
	"fmt"
	"os"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("downhub.yaml")
	if err != nil {
		// If config file doesn't exist, create default config
		cfg = config.GetDefaultConfig()
		config.SaveConfig(cfg, "downhub.yaml")
	}

	// Pass config to all modules
	cmd.SetConfig(cfg)
	common.SetConfig(cfg)
	handler.SetConfig(cfg)
	logs.SetConfig(cfg)

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
