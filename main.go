package main

import (
	"os"

	"github.com/dotpilot/cmd"
	"github.com/dotpilot/utils"
)

func main() {
	if err := cmd.Execute(); err != nil {
		utils.Logger.Error().Err(err).Msg("Error executing command")
		os.Exit(1)
	}
}
