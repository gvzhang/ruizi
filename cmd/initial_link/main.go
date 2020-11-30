package main

import (
	"fmt"
	"ruizi/internal"
	"ruizi/internal/dao"
	"ruizi/internal/service"
	"ruizi/pkg/logger"
)

func main() {
	internal.InitConfig()
	logger.Init(internal.GetConfig().RootPath)
	dao.InitLink()
	sl := new(service.Link)
	mainUrls := make([]string, 0)
	mainUrls = append(mainUrls, "https://www.qqsgjy.com/")
	for _, url := range mainUrls {
		err := sl.Add([]byte(url))
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("bootstrap success")
}
