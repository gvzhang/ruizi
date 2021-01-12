package dao

import (
	"bytes"
	"encoding/binary"
	"errors"
	"ruizi/internal/model"
	"strings"
	"testing"
)

func TestTmpIndexDoAdd(t *testing.T) {
	buffer := new(bytes.Buffer)
	termId := uint64(100)
	docId := uint64(201)
	tmpIndexModel := &model.TmpIndex{
		TermId: termId,
		DocId:  docId,
	}
	err := TmpIndex.doAdd(buffer, tmpIndexModel)
	if err != nil {
		t.Fatal(err)
	}

	expectBuffer := new(bytes.Buffer)
	binary.Write(expectBuffer, binary.LittleEndian, tmpIndexModel.TermId)
	binary.Write(expectBuffer, binary.LittleEndian, tmpIndexModel.DocId)

	t.Log("actually: " + buffer.String())
	t.Log("  expect: " + expectBuffer.String())
	if strings.Compare(buffer.String(), expectBuffer.String()) != 0 {
		t.Fatal(errors.New("expert buffer is not same"))
	}
}

func TestTmpIndexDoGetOne(t *testing.T) {
	buffer := new(bytes.Buffer)
	termId := uint64(100)
	docId := uint64(201)
	expectIndexModel := &model.TmpIndex{
		TermId:     termId,
		DocId:      docId,
		NextOffset: 16,
	}
	binary.Write(buffer, binary.LittleEndian, expectIndexModel.TermId)
	binary.Write(buffer, binary.LittleEndian, expectIndexModel.DocId)
	rs := bytes.NewReader(buffer.Bytes())

	offset := int64(0)
	tmpIndexModel, err := TmpIndex.doGetOne(rs, offset)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("actually: " + tmpIndexModel.String())
	t.Log("  expect: " + expectIndexModel.String())
	if strings.Compare(tmpIndexModel.String(), expectIndexModel.String()) != 0 {
		t.Fatal("expect model is not same")
	}
}

func TestTmpIndexGetNext(t *testing.T) {
	rows := [][]uint64{
		{101, 201},
		{203, 304},
	}

	buffer := new(bytes.Buffer)
	for _, row := range rows {
		model := &model.TmpIndex{
			TermId: row[0],
			DocId:  row[1],
		}
		binary.Write(buffer, binary.LittleEndian, model.TermId)
		binary.Write(buffer, binary.LittleEndian, model.DocId)
	}
	rs := bytes.NewReader(buffer.Bytes())

	offset := int64(0)
	for idx, row := range rows {
		expectModel := &model.TmpIndex{
			TermId:     row[0],
			DocId:      row[1],
			NextOffset: int64(idx+1) * 16,
		}
		tmpIndexModel, err := TmpIndex.doGetOne(rs, offset)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("actually: " + tmpIndexModel.String())
		t.Log("  expect: " + expectModel.String())
		if strings.Compare(tmpIndexModel.String(), expectModel.String()) != 0 {
			t.Fatalf("expect %s model is not same\n", expectModel.String())
		}
		offset = tmpIndexModel.NextOffset
	}
}
