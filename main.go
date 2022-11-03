package main

import (
	"go.uber.org/zap"
	"os"
)

var (
	logger, _ = zap.NewProduction()
)

func main() {
	err := SearchCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}
