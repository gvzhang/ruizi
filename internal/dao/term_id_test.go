package dao

import (
	"ruizi/pkg/util"
	"sync"
	"testing"
)

func TestDoTermIdIncr(t *testing.T) {
	buf := util.NewSeekableBuffer()
	var i uint64
	for i = 1; i <= 100; i++ {
		id, err := TermId.doIncr(buf)
		if err != nil {
			t.Error(err)
		}
		if id != i {
			t.Errorf("error id %d", id)
		}
	}
}

func TestDoTermIdGetConcurrent(t *testing.T) {
	buf := util.NewSeekableBuffer()
	var i uint64
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i = 1; i <= 100; i++ {
		go func() {
			_, err := TermId.doIncr(buf)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	docId, err := TermId.doIncr(buf)
	if err != nil {
		t.Error(err)
	}
	if docId != 101 {
		t.Errorf("concurrent result value error %d", docId)
	}
}
