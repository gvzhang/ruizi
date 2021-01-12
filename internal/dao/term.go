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

type term struct {
	lock *sync.RWMutex
}

var Term *term

func init() {
	Term = &term{
		lock: &sync.RWMutex{},
	}
}

func InitTerm() {
	dataPath := internal.GetConfig().Term.DataPath
	fp, err := os.OpenFile(dataPath, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
}

func (t *term) Add(tm *model.Term) error {
	dataPath := internal.GetConfig().Term.DataPath
	fp, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer fp.Close()

	term, err := TermId.Incr()
	if err != nil {
		return err
	}
	tm.Id = term
	return t.doAdd(fp, tm)
}

func (t *term) doAdd(handle io.Writer, tm *model.Term) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	txtLen := len(tm.Txt)
	if txtLen > int(internal.UINT16_MAX) {
		return errors.New("txt len can not excedded uint16 max")
	}
	termIdLen := 8

	// 内存一并写入实现原子操作
	statusLen := 1
	totalLen := int64(termIdLen + statusLen + txtLen)

	writeBuffer := new(bytes.Buffer)
	binary.Write(writeBuffer, binary.LittleEndian, totalLen)
	binary.Write(writeBuffer, binary.LittleEndian, tm.Id)
	binary.Write(writeBuffer, binary.LittleEndian, tm.Status)
	binary.Write(writeBuffer, binary.LittleEndian, tm.Txt)
	err := util.WriteBinary(handle, writeBuffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (t *term) GetOne(beginOffset int64) (*model.Term, error) {
	fp, err := os.Open(internal.GetConfig().Term.DataPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return t.doGetOne(fp, beginOffset)
}

func (t *term) doGetOne(handle io.ReadSeeker, beginOffset int64) (*model.Term, error) {
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

	var status uint8
	err = binary.Read(bytes.NewBuffer(dataByte[8:9]), binary.LittleEndian, &status)
	if err != nil {
		return nil, err
	}

	txt := make([]byte, len(dataByte)-8-1)
	err = binary.Read(bytes.NewBuffer(dataByte[9:]), binary.LittleEndian, &txt)
	if err != nil {
		return nil, err
	}

	curOffset, err := handle.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	return &model.Term{
		Id:         termId,
		Txt:        txt,
		Status:     status,
		Offset:     beginOffset,
		NextOffset: curOffset,
	}, nil
}

func (t *term) UpdateStatus(termModel *model.Term, status uint8) error {
	dataPath := internal.GetConfig().Term.DataPath
	fp, err := os.OpenFile(dataPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer fp.Close()
	return t.doUpdateStatus(fp, termModel, status)
}

func (t *term) doUpdateStatus(handle io.WriteSeeker, termModel *model.Term, status uint8) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	_, err := handle.Seek(termModel.Offset+8+8, io.SeekStart)
	if err != nil {
		return err
	}

	err = util.WriteBinary(handle, status)
	if err != nil {
		return err
	}
	termModel.Status = status
	return nil
}
