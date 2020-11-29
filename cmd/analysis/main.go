package main

import (
	"ruizi/internal"
	"ruizi/internal/analysis"
	"ruizi/pkg/logger"
)

func main() {
	internal.InitConfig()
	logger.Init(internal.GetConfig().RootPath)

	runner := analysis.Runner{}
	runner.Start()
}
