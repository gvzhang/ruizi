package dao

import (
	"bytes"
	"testing"
)

func TestDoGet(t *testing.T) {
	buf := new(bytes.Buffer)
	var i uint64
	for i = 1; i <= 100; i++ {
		id, err := DocId.doGet(buf)
		if err != nil {
			t.Error(err)
		}
		if id != i {
			t.Errorf("error id %d", id)
		}
	}
}
