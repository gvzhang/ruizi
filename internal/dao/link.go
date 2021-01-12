package dao

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"ruizi/internal"
	"ruizi/internal/model"
	"ruizi/pkg/util"
	"sync"
)

type link struct {
	lock *sync.RWMutex
}

var Link *link

func init() {
	Link = &link{
		lock: &sync.RWMutex{},
	}
}

func InitLink() {
	dataPath := internal.GetConfig().Link.DataPath
	fp, err := os.OpenFile(dataPath, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
}

func (l *link) Add(lm *model.Link) error {
	dataPath := internal.GetConfig().Link.DataPath
	fp, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer fp.Close()
	return l.doAdd(fp, lm)
}

func (l *link) doAdd(handle io.Writer, lm *model.Link) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	statusLen := 1
	uLen := len(lm.Url)
	if uLen > (int(internal.UINT16_MAX) - statusLen) {
		return errors.New("url len can not excedded uint16 max")
	}

	// 内存一并写入实现原子操作
	writeBuffer := new(bytes.Buffer)
	binary.Write(writeBuffer, binary.LittleEndian, uint16(uLen+statusLen))
	binary.Write(writeBuffer, binary.LittleEndian, lm.Status)
	binary.Write(writeBuffer, binary.LittleEndian, lm.Url)
	err := util.WriteBinary(handle, writeBuffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (l *link) GetOne(beginOffset int64) (*model.Link, error) {
	fp, err := os.Open(internal.GetConfig().Link.DataPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return l.doGetOne(fp, beginOffset)
}

func (l *link) doGetOne(handle io.ReadSeeker, beginOffset int64) (*model.Link, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	_, err := handle.Seek(beginOffset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	var dataLen uint16
	err = util.ReadBinary(handle, 2, &dataLen)
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

	var status uint8
	err = binary.Read(bytes.NewBuffer(dataByte[:1]), binary.LittleEndian, &status)
	if err != nil {
		return nil, err
	}

	url := make([]byte, len(dataByte)-1)
	err = binary.Read(bytes.NewBuffer(dataByte[1:]), binary.LittleEndian, &url)
	if err != nil {
		return nil, err
	}

	curOffset, err := handle.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	return &model.Link{
		Status:     status,
		Url:        url,
		Offset:     beginOffset,
		NextOffset: curOffset,
	}, nil
}

func (l *link) UpdateStatus(linkModel *model.Link, status uint8) error {
	dataPath := internal.GetConfig().Link.DataPath
	fp, err := os.OpenFile(dataPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer fp.Close()
	return l.doUpdateStatus(fp, linkModel, status)
}

func (l *link) doUpdateStatus(handle io.WriteSeeker, linkModel *model.Link, status uint8) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	_, err := handle.Seek(linkModel.Offset+2, io.SeekStart)
	if err != nil {
		return err
	}

	err = util.WriteBinary(handle, status)
	if err != nil {
		return err
	}
	linkModel.Status = status
	return nil
}
