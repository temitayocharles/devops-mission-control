package main

import (
	"github.com/yourusername/devops-mission-control/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
