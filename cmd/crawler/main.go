package main

import (
	"os"
	"os/signal"
	"ruizi/internal"
	"ruizi/internal/crawler"
	"ruizi/internal/dao"
	"ruizi/pkg/logger"
	"ruizi/pkg/util"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	internal.InitConfig()
	logger.Init(internal.GetConfig().RootPath)
	dao.InitLink()
	dao.InitDocId()
	dao.InitDoc()
	dao.InitDocLink()

	finishCh := make(chan struct{}, 0)
	runner := crawler.NewRunner()
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
		logger.Logger.Info("crawler stop start", zap.Any("signal", s))
		err := runner.Stop()
		if err != nil {
			logger.Logger.Error("crawler stop error", zap.Error(err))
		}
		logger.Logger.Info("crawler stop end")
	case <-finishCh:
		logger.Sugar.Info("crawler finish")
	}
}
