package main

import (
	"github.com/yourusername/ops-tool/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
