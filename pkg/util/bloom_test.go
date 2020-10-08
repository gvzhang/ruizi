package util

import (
	"errors"
	"testing"
)

const BloomSize = 2 << 24

var seeds []int8

func init() {
	seeds = []int8{4, 9, 14, 18, 23}
}

func TestSimpleHash(t *testing.T) {
	url := "http://www.test.com"
	for _, seed := range seeds {
		sh := &simpleHash{
			size: BloomSize,
			seed: seed,
		}
		sv := sh.Sum(url)
		if sv == 0 {
			t.Error(errors.New("Sum can not zero"))
		}
	}
}

func TestSet(t *testing.T) {
	blm := NewBloom(BloomSize, seeds)
	blm.Set("http://www.test.com")

	hfns := []int{6267865, 22626720, 7267963, 9668951, 3660614}
	for _, num := range hfns {
		bi := num / 8
		br := num % 8
		if blm.db[bi]&(1<<br) == 0 {
			t.Error(errors.New("hash position is zero"))
		}
	}
}

func TestGet(t *testing.T) {
	url := "http://www.test.com"
	blm := NewBloom(BloomSize, seeds)
	blm.Set(url)
	find, err := blm.Get(url)
	if err != nil {
		t.Error(err)
	}
	if find == false {
		t.Error("can find url")
	}
}
