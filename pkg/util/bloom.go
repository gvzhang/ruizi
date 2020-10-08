package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

type simpleHash struct {
	size int
	seed int8
}

func (sh *simpleHash) Sum(val string) int {
	var result int
	for _, c := range val {
		result = int(sh.seed)*result + int(c)
	}
	// 注意该哈希算法对size的设置有限制
	// 如2^次方+1，则永远返回相同值
	return (sh.size - 1) & result
}

type bloom struct {
	db  []byte
	hfs []*simpleHash
}

func NewBloom(sz int, sds []int8) *bloom {
	blm := &bloom{
		db: make([]byte, sz/8+1),
	}
	shf := make([]*simpleHash, len(sds))
	for si, sd := range sds {
		shf[si] = &simpleHash{
			size: sz,
			seed: sd,
		}
	}
	blm.hfs = shf
	return blm
}

func (b *bloom) ImportDB(idb []byte) {
	b.db = idb
}

func (b *bloom) OutputDB() []byte {
	return b.db
}

func (b *bloom) Set(s string) error {
	if len(b.hfs) == 0 {
		return errors.New("hashFunc can not empty")
	}
	for _, hf := range b.hfs {
		hv := hf.Sum(s)
		bi := hv / 8
		br := hv % 8
		b.db[bi] |= 1 << br
	}

	return nil
}

func (b *bloom) Get(s string) (bool, error) {
	if len(b.hfs) == 0 {
		return false, errors.New("hashFunc can not empty")
	}
	for _, hf := range b.hfs {
		hv := hf.Sum(s)
		bi := hv / 8
		br := hv % 8
		if b.db[bi]&(1<<br) == 0 {
			return false, nil
		}
	}
	return true, nil
}

func BloomPersistence(db []byte, filename string) error {
	fp, err := os.OpenFile(filename, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer fp.Close()
	writeBuffer := new(bytes.Buffer)
	binary.Write(writeBuffer, binary.LittleEndian, db)
	err = WriteBinary(fp, writeBuffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func BloomFileData(filename string) ([]byte, error) {
	fp, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	stat, err := fp.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := stat.Size()
	dataByte := make([]byte, fileSize)
	err = ReadBinary(fp, uint64(fileSize), &dataByte)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}
	return dataByte, nil
}
