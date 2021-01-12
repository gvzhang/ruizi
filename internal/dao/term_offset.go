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

type termOffset struct {
	lock *sync.RWMutex
}

var TermOffset *termOffset

func init() {
	TermOffset = &termOffset{
		lock: &sync.RWMutex{},
	}
}

func InitTermOffset() {
	dataPath := internal.GetConfig().TermOffset.DataPath
	fp, err := os.OpenFile(dataPath, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
}

func (dl *termOffset) Add(tim *model.TermOffset) error {
	dataPath := internal.GetConfig().TermOffset.DataPath
	fp, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer fp.Close()
	return dl.doAdd(fp, tim)
}

func (dl *termOffset) doAdd(handle io.Writer, tim *model.TermOffset) error {
	dl.lock.Lock()
	defer dl.lock.Unlock()

	// 内存一并写入实现原子操作
	writeBuffer := new(bytes.Buffer)
	binary.Write(writeBuffer, binary.LittleEndian, tim.TermId)
	binary.Write(writeBuffer, binary.LittleEndian, tim.Offset)
	err := util.WriteBinary(handle, writeBuffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (dl *termOffset) GetOne(beginOffset int64) (*model.TermOffset, error) {
	fp, err := os.Open(internal.GetConfig().TermOffset.DataPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return dl.doGetOne(fp, beginOffset)
}

func (dl *termOffset) doGetOne(handle io.ReadSeeker, beginOffset int64) (*model.TermOffset, error) {
	dl.lock.RLock()
	defer dl.lock.RUnlock()

	// o(n)查找,使用二叉树或加索引优化性能
	_, err := handle.Seek(beginOffset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	dataByte := make([]byte, 16)
	err = util.ReadBinary(handle, 16, &dataByte)
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

	var offset int64
	err = binary.Read(bytes.NewBuffer(dataByte[8:16]), binary.LittleEndian, &offset)
	if err != nil {
		return nil, err
	}

	curOffset, err := handle.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	return &model.TermOffset{
		TermId:     termId,
		Offset:      offset,
		NextOffset: curOffset,
	}, nil
}
