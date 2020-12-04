package dao

import (
	"os"
	"ruizi/internal"
	"sync"
)

type tmpIndex struct {
	lock *sync.RWMutex
}

var TmpIndex *tmpIndex

func init() {
	TmpIndex = &tmpIndex{
		lock: &sync.RWMutex{},
	}
}

func InitTmpIndex() {
	dataPath := internal.GetConfig().TmpIndex.DataPath
	fp, err := os.OpenFile(dataPath, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
}
