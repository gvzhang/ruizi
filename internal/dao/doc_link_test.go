package dao

import (
	"bytes"
	"encoding/binary"
	"errors"
	"ruizi/internal/model"
	"strings"
	"testing"
)

func TestDocLinkDoAdd(t *testing.T) {
	buffer := new(bytes.Buffer)
	docLinkModel := &model.DocLink{
		DocId: 1,
		Url:   []byte("http://www.test.com"),
	}
	err := DocLink.doAdd(buffer, docLinkModel)
	if err != nil {
		t.Fatal(err)
	}

	expectBuffer := new(bytes.Buffer)
	binary.Write(expectBuffer, binary.LittleEndian, uint16(len(docLinkModel.Url)+8))
	binary.Write(expectBuffer, binary.LittleEndian, docLinkModel.DocId)
	binary.Write(expectBuffer, binary.LittleEndian, docLinkModel.Url)

	t.Log("actually: " + buffer.String())
	t.Log("  expect: " + expectBuffer.String())
	if strings.Compare(buffer.String(), expectBuffer.String()) != 0 {
		t.Fatal(errors.New("expert buffer is not same"))
	}
}

func TestDocLinkDoGetOne(t *testing.T) {
	docId := uint64(1)
	url := []byte("http://www.test.com")
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, uint16(len(url)+8))
	binary.Write(buffer, binary.LittleEndian, docId)
	binary.Write(buffer, binary.LittleEndian, url)
	rs := bytes.NewReader(buffer.Bytes())

	bytes.NewBuffer()

	docLinkModel, err := DocLink.doGetOne(rs, docId)
	if err != nil {
		t.Fatal(err)
	}

	expectModel := &model.DocLink{
		DocId: docId,
		Url:   url,
	}
	t.Log("actually: " + docLinkModel.String())
	t.Log("  expect: " + expectModel.String())
	if strings.Compare(docLinkModel.String(), expectModel.String()) != 0 {
		t.Fatal("expect model is not same")
	}
}
