package model

import (
	"strconv"
	"strings"
)

type TmpIndex struct {
	TermId     uint64
	DocId      uint64
	NextOffset int64
}

func (d *TmpIndex) String() string {
	return strings.Join([]string{
		strconv.FormatUint(d.TermId, 10),
		strconv.FormatUint(d.DocId, 10),
		strconv.FormatInt(d.NextOffset, 10),
	}, " ")
}
