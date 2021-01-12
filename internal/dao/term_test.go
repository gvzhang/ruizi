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

func TestTermDoAdd(t *testing.T) {
	txt := []byte("苹果")
	buffer := new(bytes.Buffer)
	rawSize := uint16(len(txt))
	termModel := &model.Term{
		Id:  1,
		Txt: txt,
	}
	err := Term.doAdd(buffer, termModel)
	if err != nil {
		t.Fatal(err)
	}

	termIdLen := uint32(8)
	statusLen := uint32(1)
	totalLen := uint64(termIdLen + statusLen + uint32(rawSize))

	expectBuffer := new(bytes.Buffer)
	binary.Write(expectBuffer, binary.LittleEndian, totalLen)
	binary.Write(expectBuffer, binary.LittleEndian, termModel.Id)
	binary.Write(expectBuffer, binary.LittleEndian, termModel.Status)
	binary.Write(expectBuffer, binary.LittleEndian, termModel.Txt)

	t.Log("actually: " + buffer.String())
	t.Log("  expect: " + expectBuffer.String())
	if strings.Compare(buffer.String(), expectBuffer.String()) != 0 {
		t.Fatal(errors.New("expert buffer is not same"))
	}
}

func TestTermDoGetOne(t *testing.T) {
	txt := []byte("苹果")
	buffer := new(bytes.Buffer)
	rawSize := uint16(len(txt))
	termModel := &model.Term{
		Id:     1,
		Status: model.TermStatusWait,
		Txt:    []byte(txt),
	}
	termIdLen := uint32(8)
	statusLen := uint32(1)
	totalLen := int64(termIdLen + statusLen + uint32(rawSize))

	binary.Write(buffer, binary.LittleEndian, totalLen)
	binary.Write(buffer, binary.LittleEndian, termModel.Id)
	binary.Write(buffer, binary.LittleEndian, termModel.Status)
	binary.Write(buffer, binary.LittleEndian, termModel.Txt)

	rs := bytes.NewReader(buffer.Bytes())

	offset := int64(0)
	termModel, err := Term.doGetOne(rs, offset)
	if err != nil {
		t.Fatal(err)
	}

	expectModel := &model.Term{
		Id:         1,
		Status:     model.TermStatusWait,
		Txt:        []byte(txt),
		Offset:     offset,
		NextOffset: offset + int64(8+totalLen),
	}
	t.Log("actually: " + termModel.String())
	t.Log("  expect: " + expectModel.String())
	if strings.Compare(termModel.String(), expectModel.String()) != 0 {
		t.Fatal("expect model is not same")
	}
}

func TestTermGetNext(t *testing.T) {
	txts := [][]byte{
		[]byte("苹果"),
		[]byte("火龙果"),
	}
	termIdLen := uint32(8)
	statusLen := uint32(1)

	buffer := new(bytes.Buffer)
	for idx, txt := range txts {
		rawSize := uint16(len(txt))
		termModel := &model.Term{
			Id:     uint64(idx + 1),
			Status: model.TermStatusWait,
			Txt:    txt,
		}

		totalLen := int64(termIdLen + statusLen + uint32(rawSize))
		binary.Write(buffer, binary.LittleEndian, totalLen)
		binary.Write(buffer, binary.LittleEndian, termModel.Id)
		binary.Write(buffer, binary.LittleEndian, termModel.Status)
		binary.Write(buffer, binary.LittleEndian, termModel.Txt)
	}
	rs := bytes.NewReader(buffer.Bytes())

	offset := int64(0)
	for idx, txt := range txts {
		rawSize := uint16(len(txt))
		termModel, err := Term.doGetOne(rs, offset)
		if err != nil {
			t.Fatal(err)
		}
		totalLen := int64(termIdLen + statusLen + uint32(rawSize))
		expectModel := &model.Term{
			Id:         uint64(idx + 1),
			Status:     model.TermStatusWait,
			Txt:        []byte(txt),
			Offset:     offset,
			NextOffset: offset + int64(8+totalLen),
		}
		t.Log("actually: " + termModel.String())
		t.Log("  expect: " + expectModel.String())
		if strings.Compare(termModel.String(), expectModel.String()) != 0 {
			t.Fatalf("expect %s model is not same\n", string(txt))
		}
		offset = termModel.NextOffset
	}
}

func TestTermUpdateStatus(t *testing.T) {
	txt := []byte("苹果")
	buffer := new(bytes.Buffer)
	addModel := &model.Term{
		Id:     1,
		Status: model.TermStatusWait,
		Txt:    []byte(txt),
	}
	err := Term.doAdd(buffer, addModel)
	if err != nil {
		t.Fatal(err)
	}

	offset := int64(0)
	rs := bytes.NewReader(buffer.Bytes())
	docModel, err := Term.doGetOne(rs, offset)
	if err != nil {
		t.Fatal(err)
	}

	wr := &writerseeker.WriterSeeker{}
	wr.Write(buffer.Bytes())
	err = Term.doUpdateStatus(wr, docModel, model.TermStatusAnalysis)
	if err != nil {
		t.Fatal(err)
	}

	actuallyModel, err := Term.doGetOne(wr.BytesReader(), offset)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("actually: " + actuallyModel.String())
	t.Log("  expect: " + docModel.String())
	if strings.Compare(actuallyModel.String(), docModel.String()) != 0 {
		t.Fatal(errors.New("expect model is not same"))
	}
}
