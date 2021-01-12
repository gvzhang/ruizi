package model

import (
	"strconv"
	"strings"
)

type Index struct {
	TermId     uint64
	DocIdList  []uint64
	NextOffset int64
}

func (i *Index) String() string {
	docIdListString := make([]string, 0)
	for _, v := range i.DocIdList {
		docIdListString = append(docIdListString, strconv.FormatUint(v, 10))
	}
	return strings.Join([]string{
		strconv.FormatUint(i.TermId, 10),
		strings.Join(docIdListString, ","),
		strconv.FormatInt(i.NextOffset, 10),
	}, " ")
}
