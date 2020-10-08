package model

import (
	"strconv"
	"strings"
)

const LinkStatusWait = 1
const LinkStatusDone = 2

type Link struct {
	Status     uint8
	Url        []byte
	Offset     int64
	NextOffset int64
}

func (l *Link) String() string {
	return strings.Join([]string{
		string(l.Url),
		strconv.Itoa(int(l.Status)),
		strconv.FormatInt(l.Offset, 10),
		strconv.FormatInt(l.NextOffset, 10),
	}, " ")
}
