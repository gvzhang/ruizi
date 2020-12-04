package dao

import (
	"os"
	"ruizi/internal"
	"sync"
)

type termId struct {
	lock *sync.RWMutex
}

var TermId *termId

func init() {
	TermId = &termId{
		lock: &sync.RWMutex{},
	}
}

func InitTermId() {
	dataPath := internal.GetConfig().TermId.DataPath
	fp, err := os.OpenFile(dataPath, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
}

