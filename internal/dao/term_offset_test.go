package dao

import (
	"bytes"
	"encoding/binary"
	"errors"
	"ruizi/internal/model"
	"strings"
	"testing"
)

func TestTermOffsetDoAdd(t *testing.T) {
	buffer := new(bytes.Buffer)
	termId := uint64(100)
	offset := int64(12054)
	termOffsetModel := &model.TermOffset{
		TermId: termId,
		Offset: offset,
	}
	err := TermOffset.doAdd(buffer, termOffsetModel)
	if err != nil {
		t.Fatal(err)
	}

	expectBuffer := new(bytes.Buffer)
	binary.Write(expectBuffer, binary.LittleEndian, termOffsetModel.TermId)
	binary.Write(expectBuffer, binary.LittleEndian, termOffsetModel.Offset)

	t.Log("actually: " + buffer.String())
	t.Log("  expect: " + expectBuffer.String())
	if strings.Compare(buffer.String(), expectBuffer.String()) != 0 {
		t.Fatal(errors.New("expert buffer is not same"))
	}
}

func TestTermOffsetDoGetOne(t *testing.T) {
	buffer := new(bytes.Buffer)
	termId := uint64(100)
	offset := int64(12054)
	expectIndexModel := &model.TermOffset{
		TermId:     termId,
		Offset:     offset,
		NextOffset: 16,
	}
	binary.Write(buffer, binary.LittleEndian, expectIndexModel.TermId)
	binary.Write(buffer, binary.LittleEndian, expectIndexModel.Offset)
	rs := bytes.NewReader(buffer.Bytes())

	beginOffset := int64(0)
	termOffsetModel, err := TermOffset.doGetOne(rs, beginOffset)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("actually: " + termOffsetModel.String())
	t.Log("  expect: " + expectIndexModel.String())
	if strings.Compare(termOffsetModel.String(), expectIndexModel.String()) != 0 {
		t.Fatal("expect model is not same")
	}
}

func TestTermOffsetGetNext(t *testing.T) {
	rows := [][]uint64{
		{101, 156421},
		{203, 345787},
	}

	buffer := new(bytes.Buffer)
	for _, row := range rows {
		termOffset := &model.TermOffset{
			TermId: row[0],
			Offset: int64(row[1]),
		}
		binary.Write(buffer, binary.LittleEndian, termOffset.TermId)
		binary.Write(buffer, binary.LittleEndian, termOffset.Offset)
	}
	rs := bytes.NewReader(buffer.Bytes())

	offset := int64(0)
	for idx, row := range rows {
		expectModel := &model.TermOffset{
			TermId:     row[0],
			Offset:     int64(row[1]),
			NextOffset: int64(idx+1) * 16,
		}
		termOffsetModel, err := TermOffset.doGetOne(rs, offset)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("actually: " + termOffsetModel.String())
		t.Log("  expect: " + expectModel.String())
		if strings.Compare(termOffsetModel.String(), expectModel.String()) != 0 {
			t.Fatalf("expect %s model is not same\n", expectModel.String())
		}
		offset = termOffsetModel.NextOffset
	}
}
