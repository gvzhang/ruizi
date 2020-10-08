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

func TestDoAdd(t *testing.T) {
	buffer := new(bytes.Buffer)
	linkModel := &model.Link{
		Url:    []byte("http://www.test.com"),
		Status: model.LinkStatusWait,
	}
	err := Link.doAdd(buffer, linkModel)
	if err != nil {
		t.Fatal(err)
	}

	expectBuffer := new(bytes.Buffer)
	binary.Write(expectBuffer, binary.LittleEndian, uint16(len(linkModel.Url)+1))
	binary.Write(expectBuffer, binary.LittleEndian, linkModel.Status)
	binary.Write(expectBuffer, binary.LittleEndian, linkModel.Url)

	t.Log("actually: " + buffer.String())
	t.Log("  expect: " + expectBuffer.String())
	if strings.Compare(buffer.String(), expectBuffer.String()) != 0 {
		t.Fatal(errors.New("expert buffer is not same"))
	}
}

func TestDoGetOne(t *testing.T) {
	url := []byte("http://www.test.com")
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, uint16(len(url)+1))
	binary.Write(buffer, binary.LittleEndian, uint8(1))
	binary.Write(buffer, binary.LittleEndian, url)
	rs := bytes.NewReader(buffer.Bytes())

	offset := int64(0)
	linkModel, err := Link.doGetOne(rs, offset)
	if err != nil {
		t.Fatal(err)
	}

	expectModel := &model.Link{
		Status:     model.LinkStatusWait,
		Url:        url,
		Offset:     offset,
		NextOffset: offset + int64(2+1+len(url)),
	}
	t.Log("actually: " + linkModel.String())
	t.Log("  expect: " + expectModel.String())
	if strings.Compare(linkModel.String(), expectModel.String()) != 0 {
		t.Fatal("expect model is not same")
	}
}

func TestGetNext(t *testing.T) {
	urls := [][]byte{
		[]byte("http://www.test.com"),
		[]byte("http://www.ruizi.com"),
	}

	buffer := new(bytes.Buffer)
	for _, url := range urls {
		binary.Write(buffer, binary.LittleEndian, uint16(len(url)+1))
		binary.Write(buffer, binary.LittleEndian, uint8(1))
		binary.Write(buffer, binary.LittleEndian, url)
	}
	rs := bytes.NewReader(buffer.Bytes())

	offset := int64(0)
	for _, url := range urls {
		linkModel, err := Link.doGetOne(rs, offset)
		if err != nil {
			t.Fatal(err)
		}

		expectModel := &model.Link{
			Status:     model.LinkStatusWait,
			Url:        url,
			Offset:     offset,
			NextOffset: offset + int64(2+1+len(url)),
		}
		t.Log("actually: " + linkModel.String())
		t.Log("  expect: " + expectModel.String())
		if strings.Compare(linkModel.String(), expectModel.String()) != 0 {
			t.Fatalf("expect %s model is not same\n", string(url))
		}
		offset = linkModel.NextOffset
	}
}

func TestUpdateStatus(t *testing.T) {
	buffer := new(bytes.Buffer)
	addModel := &model.Link{
		Status: model.LinkStatusWait,
		Url:    []byte("http://www.test.com"),
	}
	err := Link.doAdd(buffer, addModel)
	if err != nil {
		t.Fatal(err)
	}

	offset := int64(0)
	rs := bytes.NewReader(buffer.Bytes())
	linkModel, err := Link.doGetOne(rs, offset)
	if err != nil {
		t.Fatal(err)
	}

	wr := &writerseeker.WriterSeeker{}
	wr.Write(buffer.Bytes())
	err = Link.doUpdateStatus(wr, linkModel, model.LinkStatusDone)
	if err != nil {
		t.Fatal(err)
	}

	actuallyModel, err := Link.doGetOne(wr.BytesReader(), offset)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("actually: " + actuallyModel.String())
	t.Log("  expect: " + linkModel.String())
	if strings.Compare(actuallyModel.String(), linkModel.String()) != 0 {
		t.Fatal(errors.New("expect model is not same"))
	}
}
