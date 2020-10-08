package model

import (
	"strconv"
	"strings"
)

type DocLink struct {
	DocId uint64
	Url   []byte
}

func (dl *DocLink) String() string {
	return strings.Join([]string{
		strconv.FormatUint(dl.DocId, 10),
		string(dl.Url),
	}, " ")
}
