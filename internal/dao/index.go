package dao

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"ruizi/internal"
	"ruizi/internal/model"
	"ruizi/pkg/util"
	"sync"
)

type index struct {
	lock *sync.RWMutex
}

var Index *index

func init() {
	Index = &index{
		lock: &sync.RWMutex{},
	}
}

func InitIndex() {
	dataPath := internal.GetConfig().Index.DataPath
	fp, err := os.OpenFile(dataPath, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
}

func (t *index) Add(tm *model.Index) (int64, error) {
	dataPath := internal.GetConfig().Index.DataPath
	fp, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return 0, err
	}
	defer fp.Close()
	return t.doAdd(fp, tm)
}

func (t *index) doAdd(handle io.WriteSeeker, tm *model.Index) (int64, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	// 内存一并写入实现原子操作
	termIdLen := 8
	totalLen := int64(termIdLen + len(tm.DocIdList)*8)

	writeBuffer := new(bytes.Buffer)
	binary.Write(writeBuffer, binary.LittleEndian, totalLen)
	binary.Write(writeBuffer, binary.LittleEndian, tm.TermId)
	for _, v := range tm.DocIdList {
		binary.Write(writeBuffer, binary.LittleEndian, v)
	}
	err := util.WriteBinary(handle, writeBuffer.Bytes())
	if err != nil {
		return 0, err
	}
	curOffset, err := handle.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	return curOffset, nil
}

func (t *index) GetOne(beginOffset int64) (*model.Index, error) {
	fp, err := os.Open(internal.GetConfig().Index.DataPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return t.doGetOne(fp, beginOffset)
}

func (t *index) doGetOne(handle io.ReadSeeker, beginOffset int64) (*model.Index, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	_, err := handle.Seek(beginOffset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	var dataLen int64
	err = util.ReadBinary(handle, 8, &dataLen)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	dataByte := make([]byte, dataLen)
	err = util.ReadBinary(handle, uint64(dataLen), &dataByte)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	var termId uint64
	err = binary.Read(bytes.NewBuffer(dataByte[:8]), binary.LittleEndian, &termId)
	if err != nil {
		return nil, err
	}

	var docIdList []uint64
	var i int64
	docIdListLen := dataLen - 8
	for i = 8; i <= docIdListLen; i += 8 {
		var docId uint64
		err = binary.Read(bytes.NewBuffer(dataByte[i:i+8]), binary.LittleEndian, &docId)
		if err != nil {
			return nil, err
		}
		docIdList = append(docIdList, docId)
	}

	curOffset, err := handle.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	return &model.Index{
		TermId:     termId,
		DocIdList:  docIdList,
		NextOffset: curOffset,
	}, nil
}
