package main

import (
	"os"
	"os/signal"
	"ruizi/internal"
	"ruizi/internal/analysis"
	"ruizi/pkg/logger"
	"ruizi/pkg/util"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	internal.InitConfig()
	logger.Init(internal.GetConfig().RootPath)

	runner := analysis.Runner{}
	go func() {
		defer util.RecoverPanic()
		runner.Start()
	}()

	c := make(chan os.Signal, 0)
	signal.Notify(c, syscall.SIGTERM)

	s := <-c
	logger.Logger.Info("analysis stop start", zap.Any("signal", s))
	err := runner.Stop()
	if err != nil {
		logger.Logger.Error("analysis stop error", zap.Error(err))
	}
	logger.Logger.Info("analysis stop end")

}
