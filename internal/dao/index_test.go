package dao

import (
	"bytes"
	"encoding/binary"
	"errors"
	"ruizi/internal/model"
	"strings"
	"testing"
)

func TestIndexDoAdd(t *testing.T) {
	buffer := new(bytes.Buffer)
	docIdList := []uint64{11, 34, 56, 908}
	indexModel := &model.Index{
		TermId:    1,
		DocIdList: docIdList,
	}
	err := Index.doAdd(buffer, indexModel)
	if err != nil {
		t.Fatal(err)
	}

	termIdLen := uint32(8)
	totalLen := int64(termIdLen + uint32(len(docIdList)*8))

	expectBuffer := new(bytes.Buffer)
	binary.Write(expectBuffer, binary.LittleEndian, totalLen)
	binary.Write(expectBuffer, binary.LittleEndian, indexModel.TermId)
	for _, v := range indexModel.DocIdList {
		binary.Write(expectBuffer, binary.LittleEndian, v)
	}

	t.Log("actually: " + buffer.String())
	t.Log("  expect: " + expectBuffer.String())
	if strings.Compare(buffer.String(), expectBuffer.String()) != 0 {
		t.Fatal(errors.New("expert buffer is not same"))
	}
}

func TestIndexDoGetOne(t *testing.T) {
	buffer := new(bytes.Buffer)
	docIdList := []uint64{11, 34, 56, 908}
	indexModel := &model.Index{
		TermId:    1,
		DocIdList: docIdList,
	}
	termIdLen := uint32(8)
	totalLen := int64(termIdLen + uint32(len(docIdList)*8))
	binary.Write(buffer, binary.LittleEndian, totalLen)
	binary.Write(buffer, binary.LittleEndian, indexModel.TermId)
	for _, v := range indexModel.DocIdList {
		binary.Write(buffer, binary.LittleEndian, v)
	}

	rs := bytes.NewReader(buffer.Bytes())

	offset := int64(0)
	indexModel, err := Index.doGetOne(rs, offset)
	if err != nil {
		t.Fatal(err)
	}

	expectModel := &model.Index{
		TermId:     1,
		DocIdList:  docIdList,
		NextOffset: offset + 8 + totalLen,
	}
	t.Log("actually: " + indexModel.String())
	t.Log("  expect: " + expectModel.String())
	if strings.Compare(indexModel.String(), expectModel.String()) != 0 {
		t.Fatal("expect model is not same")
	}
}

func TestIndexGetNext(t *testing.T) {
	docIdLists := [][]uint64{
		{11, 34, 56, 908},
		{67, 876, 234, 123},
	}
	termIdLen := uint32(8)

	buffer := new(bytes.Buffer)
	for idx, docIdList := range docIdLists {
		indexModel := &model.Index{
			TermId:    uint64(idx + 1),
			DocIdList: docIdList,
		}

		totalLen := int64(termIdLen + uint32(len(docIdList)*8))
		binary.Write(buffer, binary.LittleEndian, totalLen)
		binary.Write(buffer, binary.LittleEndian, indexModel.TermId)
		for _, v := range indexModel.DocIdList {
			binary.Write(buffer, binary.LittleEndian, v)
		}
	}
	rs := bytes.NewReader(buffer.Bytes())

	offset := int64(0)
	for idx, docIdList := range docIdLists {
		indexModel, err := Index.doGetOne(rs, offset)
		if err != nil {
			t.Fatal(err)
		}
		totalLen := int64(termIdLen + uint32(len(docIdList)*8))
		expectModel := &model.Index{
			TermId:     uint64(idx + 1),
			DocIdList:  docIdList,
			NextOffset: offset + 8 + totalLen,
		}
		t.Log("actually: " + indexModel.String())
		t.Log("  expect: " + expectModel.String())
		if strings.Compare(indexModel.String(), expectModel.String()) != 0 {
			t.Fatalf("expect %d %v model is not same\n", idx, docIdList)
		}
		offset = indexModel.NextOffset
	}
}
