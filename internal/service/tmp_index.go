package service

import (
	"ruizi/internal/dao"
	"ruizi/internal/model"
)

type TmpIndex struct {
}

func (ti *TmpIndex) Add(termId uint64, docId uint64) error {
	addModel := &model.TmpIndex{
		TermId: termId,
		DocId:  docId,
	}
	err := dao.TmpIndex.Add(addModel)
	if err != nil {
		return err
	}
	return nil
}

func (ti *TmpIndex) GetOne(beginOffset int64) (*model.TmpIndex, error) {
	return dao.TmpIndex.GetOne(beginOffset)
}

func (ti *TmpIndex) GetAll(offset int64) ([]*model.TmpIndex, error) {
	tmpIndexList := make([]*model.TmpIndex, 0)
	for {
		row, err := ti.GetOne(offset)
		if err != nil {
			return nil, err
		}
		if row == nil {
			break
		}
		tmpIndexList = append(tmpIndexList, row)
		offset = row.NextOffset
	}
	return tmpIndexList, nil
}

func (ti *TmpIndex) Sort(arr []*model.TmpIndex) []*model.TmpIndex {
	if len(arr) < 2 {
		return arr
	}

	i := len(arr) / 2
	left := ti.Sort(arr[0:i])
	right := ti.Sort(arr[i:])

	return ti.merge(left, right)
}

func (ti *TmpIndex) merge(l []*model.TmpIndex, r []*model.TmpIndex) []*model.TmpIndex {
	i, j := 0, 0
	m, n := len(l), len(r)
	var res []*model.TmpIndex
	for i < m && j < n {
		if l[i].TermId < r[j].TermId {
			res = append(res, l[i])
			i++
		} else {
			res = append(res, r[j])
			j++
		}
	}
	res = append(res, l[i:]...)
	res = append(res, r[j:]...)
	return res
}
