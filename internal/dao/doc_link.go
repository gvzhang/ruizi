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

type docLink struct {
	lock *sync.RWMutex
}

var DocLink *docLink

func init() {
	DocLink = &docLink{
		lock: &sync.RWMutex{},
	}
}

func InitDocLink() {
	dataPath := internal.GetConfig().DocLink.DataPath
	fp, err := os.OpenFile(dataPath, os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
}

func (dl *docLink) Add(dlm *model.DocLink) error {
	dataPath := internal.GetConfig().DocLink.DataPath
	fp, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer fp.Close()
	return dl.doAdd(fp, dlm)
}

func (dl *docLink) doAdd(handle io.Writer, dlm *model.DocLink) error {
	dl.lock.Lock()
	defer dl.lock.Unlock()

	docIdLen := 8
	uLen := len(dlm.Url)
	if uLen > (int(internal.UINT16_MAX) - docIdLen) {
		return errors.New("url len can not excedded uint16 max")
	}

	// 内存一并写入实现原子操作
	writeBuffer := new(bytes.Buffer)
	binary.Write(writeBuffer, binary.LittleEndian, uint16(docIdLen+uLen))
	binary.Write(writeBuffer, binary.LittleEndian, dlm.DocId)
	binary.Write(writeBuffer, binary.LittleEndian, dlm.Url)
	err := util.WriteBinary(handle, writeBuffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (dl *docLink) GetOne(docId uint64) (*model.DocLink, error) {
	fp, err := os.Open(internal.GetConfig().DocLink.DataPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return dl.doGetOne(fp, docId)
}

func (dl *docLink) doGetOne(handle io.ReadSeeker, docId uint64) (*model.DocLink, error) {
	dl.lock.RLock()
	defer dl.lock.RUnlock()

	// o(n)查找,使用二叉树或加索引优化性能
	docIdLen := 8
	offset := int64(0)
	for {
		_, err := handle.Seek(offset, os.SEEK_SET)
		if err != nil {
			return nil, err
		}

		var dataLen uint16
		err = util.ReadBinary(handle, 2, &dataLen)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		dataByte := make([]byte, dataLen)
		err = util.ReadBinary(handle, uint64(dataLen), &dataByte)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		var did uint64
		err = binary.Read(bytes.NewBuffer(dataByte[:docIdLen]), binary.LittleEndian, &did)
		if err != nil {
			return nil, err
		}

		url := make([]byte, len(dataByte)-docIdLen)
		err = binary.Read(bytes.NewBuffer(dataByte[docIdLen:]), binary.LittleEndian, &url)
		if err != nil {
			return nil, err
		}

		offset, err = handle.Seek(0, os.SEEK_CUR)
		if err != nil {
			return nil, err
		}

		if did == docId {
			return &model.DocLink{
				DocId: docId,
				Url:   url,
			}, nil
		}
	}
	return nil, nil
}
