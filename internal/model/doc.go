package model

import (
	"strconv"
	"strings"
)

const DocStatusWait = 1
const DocStatusAnalysis = 2

type Doc struct {
	Id         uint64
	Status     uint8
	Size       uint32
	Raw        []byte
	Offset     int64
	NextOffset int64
}

func (d *Doc) String() string {
	return strings.Join([]string{
		strconv.FormatUint(d.Id, 10),
		strconv.Itoa(int(d.Status)),
		strconv.Itoa(int(d.Size)),
		string(d.Raw),
		strconv.FormatInt(d.Offset, 10),
		strconv.FormatInt(d.NextOffset, 10),
	}, " ")
}
