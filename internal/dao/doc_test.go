package dao

import (
	"bytes"
	"encoding/binary"
	"errors"
	"ruizi/internal/model"
	"strings"
	"testing"

	"github.com/orcaman/writerseeker"
)

const DefaultRaw = "<body><a href=\"http://www.rz.com\">Hello World!</a></body>"

func TestDoAdd(t *testing.T) {
	buffer := new(bytes.Buffer)
	rawSize := uint32(len(DefaultRaw))
	docModel := &model.Doc{
		Id:     1,
		Status: model.DocStatusWait,
		Raw:    []byte(DefaultRaw),
		Size:   rawSize,
	}
	err := Doc.doAdd(buffer, docModel)
	if err != nil {
		t.Fatal(err)
	}

	docIdLen := uint32(8)
	statusLen := uint32(1)
	rawSizeLen := uint32(4)
	totalLen := uint64(docIdLen + statusLen + rawSizeLen + rawSize)

	expectBuffer := new(bytes.Buffer)
	binary.Write(expectBuffer, binary.LittleEndian, totalLen)
	binary.Write(expectBuffer, binary.LittleEndian, docModel.Id)
	binary.Write(expectBuffer, binary.LittleEndian, docModel.Status)
	binary.Write(expectBuffer, binary.LittleEndian, docModel.Size)
	binary.Write(expectBuffer, binary.LittleEndian, docModel.Raw)

	t.Log("actually: " + buffer.String())
	t.Log("  expect: " + expectBuffer.String())
	if strings.Compare(buffer.String(), expectBuffer.String()) != 0 {
		t.Fatal(errors.New("expert buffer is not same"))
	}
}

func TestDoGetOne(t *testing.T) {
	buffer := new(bytes.Buffer)
	rawSize := uint32(len(DefaultRaw))
	docModel := &model.Doc{
		Id:     1,
		Status: model.DocStatusWait,
		Raw:    []byte(DefaultRaw),
		Size:   rawSize,
	}
	docIdLen := uint32(8)
	statusLen := uint32(1)
	rawSizeLen := uint32(4)
	totalLen := int64(docIdLen + statusLen + rawSizeLen + rawSize)

	binary.Write(buffer, binary.LittleEndian, totalLen)
	binary.Write(buffer, binary.LittleEndian, docModel.Id)
	binary.Write(buffer, binary.LittleEndian, docModel.Status)
	binary.Write(buffer, binary.LittleEndian, docModel.Size)
	binary.Write(buffer, binary.LittleEndian, docModel.Raw)

	rs := bytes.NewReader(buffer.Bytes())

	offset := int64(0)
	docModel, err := Doc.doGetOne(rs, offset)
	if err != nil {
		t.Fatal(err)
	}

	expectModel := &model.Doc{
		Id:         1,
		Status:     model.DocStatusWait,
		Raw:        []byte(DefaultRaw),
		Size:       rawSize,
		Offset:     offset,
		NextOffset: offset + int64(8+totalLen),
	}
	t.Log("actually: " + docModel.String())
	t.Log("  expect: " + expectModel.String())
	if strings.Compare(docModel.String(), expectModel.String()) != 0 {
		t.Fatal("expect model is not same")
	}
}

func TestGetNext(t *testing.T) {
	raws := [][]byte{
		[]byte(DefaultRaw),
		[]byte("<p><a href=\"http://www.ruizi.com\">ruizi</a></p>"),
	}
	docIdLen := uint32(8)
	statusLen := uint32(1)
	rawSizeLen := uint32(4)

	buffer := new(bytes.Buffer)
	for idx, raw := range raws {
		rawSize := uint32(len(raw))
		docModel := &model.Doc{
			Id:     uint64(idx + 1),
			Status: model.DocStatusWait,
			Raw:    []byte(raw),
			Size:   rawSize,
		}

		totalLen := int64(docIdLen + statusLen + rawSizeLen + rawSize)
		binary.Write(buffer, binary.LittleEndian, totalLen)
		binary.Write(buffer, binary.LittleEndian, docModel.Id)
		binary.Write(buffer, binary.LittleEndian, docModel.Status)
		binary.Write(buffer, binary.LittleEndian, docModel.Size)
		binary.Write(buffer, binary.LittleEndian, docModel.Raw)
	}
	rs := bytes.NewReader(buffer.Bytes())

	offset := int64(0)
	for idx, raw := range raws {
		rawSize := uint32(len(raw))
		docModel, err := Doc.doGetOne(rs, offset)
		if err != nil {
			t.Fatal(err)
		}
		totalLen := int64(docIdLen + statusLen + rawSizeLen + rawSize)
		expectModel := &model.Doc{
			Id:         uint64(idx + 1),
			Status:     model.DocStatusWait,
			Raw:        []byte(raw),
			Size:       rawSize,
			Offset:     offset,
			NextOffset: offset + int64(8+totalLen),
		}
		t.Log("actually: " + docModel.String())
		t.Log("  expect: " + expectModel.String())
		if strings.Compare(docModel.String(), expectModel.String()) != 0 {
			t.Fatalf("expect %s model is not same\n", string(raw))
		}
		offset = docModel.NextOffset
	}
}

func TestUpdateStatus(t *testing.T) {
	buffer := new(bytes.Buffer)
	rawSize := uint32(len(DefaultRaw))
	addModel := &model.Doc{
		Id:     1,
		Status: model.DocStatusWait,
		Raw:    []byte(DefaultRaw),
		Size:   rawSize,
	}
	err := Doc.doAdd(buffer, addModel)
	if err != nil {
		t.Fatal(err)
	}

	offset := int64(0)
	rs := bytes.NewReader(buffer.Bytes())
	docModel, err := Doc.doGetOne(rs, offset)
	if err != nil {
		t.Fatal(err)
	}

	wr := &writerseeker.WriterSeeker{}
	wr.Write(buffer.Bytes())
	err = Doc.doUpdateStatus(wr, docModel, model.DocStatusAnalysis)
	if err != nil {
		t.Fatal(err)
	}

	actuallyModel, err := Doc.doGetOne(wr.BytesReader(), offset)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("actually: " + actuallyModel.String())
	t.Log("  expect: " + docModel.String())
	if strings.Compare(actuallyModel.String(), docModel.String()) != 0 {
		t.Fatal(errors.New("expect model is not same"))
	}
}
