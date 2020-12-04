package main

import (
	"go.uber.org/zap"
	"os"
	"os/signal"
	"ruizi/internal"
	"ruizi/internal/analysis"
	"ruizi/internal/dao"
	"ruizi/pkg/logger"
	"ruizi/pkg/util"
	"syscall"
)

func main() {
	internal.InitConfig()
	logger.Init(internal.GetConfig().RootPath)
	dao.InitTermId()
	dao.InitTmpIndex()

	finishCh := make(chan struct{}, 0)
	runner := analysis.NewRunner()
	go func(finishCh chan struct{}) {
		defer util.RecoverPanic()
		err := runner.Start()
		if err != nil {
			logger.Sugar.Error(err)
		}
		finishCh <- struct{}{}
	}(finishCh)

	c := make(chan os.Signal, 0)
	signal.Notify(c, syscall.SIGTERM)

	select {
	case s := <-c:
		logger.Logger.Info("analysis stop start", zap.Any("signal", s))
		err := runner.Stop()
		if err != nil {
			logger.Logger.Error("analysis stop error", zap.Error(err))
		}
		logger.Logger.Info("analysis stop end")
	case <-finishCh:
		logger.Sugar.Info("analysis finish")
	}
}
