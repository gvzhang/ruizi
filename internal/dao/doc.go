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

type doc struct {
	lock *sync.RWMutex
}

var Doc *doc

func init() {
	Doc = &doc{
		lock: &sync.RWMutex{},
	}
}

func InitDoc() {
	dataPath := internal.GetConfig().Doc.DataPath
	fp, err := os.OpenFile(dataPath, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
}

func (d *doc) Add(dm *model.Doc) error {
	dataPath := internal.GetConfig().Doc.DataPath
	fp, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer fp.Close()

	docId, err := DocId.Get()
	if err != nil {
		return err
	}
	dm.Id = docId
	return d.doAdd(fp, dm)
}

func (d *doc) doAdd(handle io.Writer, dm *model.Doc) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	rawSize := len(dm.Raw)
	if int64(rawSize) > int64(internal.UINT32_MAX) {
		return errors.New("raw len can not excedded uint32 max")
	}
	dm.Size = uint32(rawSize)
	rawSizeLen := 4
	docIdLen := 8

	// 内存一并写入实现原子操作
	statusLen := 1
	totalLen := int64(docIdLen + statusLen + rawSizeLen + rawSize)

	writeBuffer := new(bytes.Buffer)
	binary.Write(writeBuffer, binary.LittleEndian, totalLen)
	binary.Write(writeBuffer, binary.LittleEndian, dm.Id)
	binary.Write(writeBuffer, binary.LittleEndian, dm.Status)
	binary.Write(writeBuffer, binary.LittleEndian, dm.Size)
	binary.Write(writeBuffer, binary.LittleEndian, dm.Raw)
	err := util.WriteBinary(handle, writeBuffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (d *doc) GetOne(beginOffset int64) (*model.Doc, error) {
	fp, err := os.Open(internal.GetConfig().Doc.DataPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return d.doGetOne(fp, beginOffset)
}

func (d *doc) doGetOne(handle io.ReadSeeker, beginOffset int64) (*model.Doc, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	_, err := handle.Seek(beginOffset, os.SEEK_SET)
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

	var docId uint64
	err = binary.Read(bytes.NewBuffer(dataByte[:8]), binary.LittleEndian, &docId)
	if err != nil {
		return nil, err
	}

	var status uint8
	err = binary.Read(bytes.NewBuffer(dataByte[8:9]), binary.LittleEndian, &status)
	if err != nil {
		return nil, err
	}

	var size uint32
	err = binary.Read(bytes.NewBuffer(dataByte[9:13]), binary.LittleEndian, &size)
	if err != nil {
		return nil, err
	}

	raw := make([]byte, len(dataByte)-8-1-4)
	err = binary.Read(bytes.NewBuffer(dataByte[13:]), binary.LittleEndian, &raw)
	if err != nil {
		return nil, err
	}

	curOffset, err := handle.Seek(0, os.SEEK_CUR)
	if err != nil {
		return nil, err
	}

	return &model.Doc{
		Id:         docId,
		Status:     status,
		Size:       size,
		Raw:        raw,
		Offset:     beginOffset,
		NextOffset: curOffset,
	}, nil
}

func (d *doc) UpdateStatus(docModel *model.Doc, status uint8) error {
	dataPath := internal.GetConfig().Doc.DataPath
	fp, err := os.OpenFile(dataPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer fp.Close()
	return d.doUpdateStatus(fp, docModel, status)
}

func (d *doc) doUpdateStatus(handle io.WriteSeeker, docModel *model.Doc, status uint8) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	_, err := handle.Seek(docModel.Offset+8+8, os.SEEK_SET)
	if err != nil {
		return err
	}

	err = util.WriteBinary(handle, status)
	if err != nil {
		return err
	}
	docModel.Status = status
	return nil
}
