package dao

import (
	"io"
	"os"
	"ruizi/internal"
	"ruizi/pkg/util"
	"sync"
)

type docId struct {
	lock *sync.RWMutex
}

var DocId *docId

func init() {
	DocId = &docId{
		lock: &sync.RWMutex{},
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

	di.lock.RLock()
	defer di.lock.RUnlock()

	var id uint64
	_, err = fp.Seek(0, os.SEEK_SET)
	if err != nil {
		return 0, err
	}
	err = util.ReadBinary(fp, 8, &id)
	if err != nil && err != io.EOF {
		return 0, err
	}
	return id, nil
}

func (di *docId) Incr() (uint64, error) {
	dataPath := internal.GetConfig().DocId.DataPath
	fp, err := os.OpenFile(dataPath, os.O_RDWR, 0)
	if err != nil {
		return 0, err
	}
	defer fp.Close()

	return di.doIncr(fp)
}

func (di *docId) doIncr(handle io.ReadWriteSeeker) (uint64, error) {
	di.lock.Lock()
	defer di.lock.Unlock()

	var maxId uint64

	_, err := handle.Seek(0, os.SEEK_SET)
	if err != nil {
		return 0, err
	}
	err = util.ReadBinary(handle, 8, &maxId)
	if err != nil && err != io.EOF {
		return 0, err
	}

	_, err = handle.Seek(0, os.SEEK_SET)
	if err != nil {
		return 0, err
	}
	maxId += 1
	err = util.WriteBinary(handle, maxId)
	if err != nil {
		return 0, err
	}

	return maxId, nil
}
