package model

import (
	"strconv"
	"strings"
)

const TermStatusEnable = 1
const TermStatusDisable = 2

type Term struct {
	Id         uint64
	Txt        []byte
	Status     uint8
	Offset     int64
	NextOffset int64
}

func (d *Term) String() string {
	return strings.Join([]string{
		strconv.FormatUint(d.Id, 10),
		string(d.Txt),
		strconv.Itoa(int(d.Status)),
		strconv.Itoa(int(d.Offset)),
		strconv.FormatInt(d.NextOffset, 10),
	}, " ")
}
