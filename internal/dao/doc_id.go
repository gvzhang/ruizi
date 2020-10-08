package dao

import (
	"io"
	"os"
	"ruizi/internal"
	"ruizi/pkg/util"
	"sync"
)

type docId struct {
	lock *sync.Mutex
}

var DocId *docId

func init() {
	DocId = &docId{
		lock: &sync.Mutex{},
	}
}

func InitDocId() {
	dataPath := internal.GetConfig().DocId.DataPath
	fp, err := os.OpenFile(dataPath, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
}

func (di *docId) Get() (uint64, error) {
	dataPath := internal.GetConfig().DocId.DataPath
	fp, err := os.OpenFile(dataPath, os.O_RDWR, 0)
	if err != nil {
		return 0, err
	}
	defer fp.Close()

	return di.doGet(fp)
}

func (di *docId) doGet(handle io.ReadWriter) (uint64, error) {
	di.lock.Lock()
	defer di.lock.Unlock()

	var maxId uint64
	err := util.ReadBinary(handle, 8, &maxId)
	if err != nil && err != io.EOF {
		return 0, err
	}

	maxId += 1
	err = util.WriteBinary(handle, maxId)
	if err != nil {
		return 0, err
	}

	return maxId, nil
}
