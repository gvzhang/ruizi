package model

import (
	"strconv"
	"strings"
)

type TermOffset struct {
	TermId     uint64
	Offset     int64
	NextOffset int64
}

func (to *TermOffset) String() string {
	return strings.Join([]string{
		strconv.FormatUint(to.TermId, 10),
		strconv.FormatInt(to.Offset, 10),
		strconv.FormatInt(to.NextOffset, 10),
	}, " ")
}
