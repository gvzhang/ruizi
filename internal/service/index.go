package service

import (
	"ruizi/internal/dao"
	"ruizi/internal/model"
)

type Index struct {
}

func (i *Index) Add(termId uint64, docIdList []uint64) (int64, error) {
	addModel := &model.Index{
		TermId:    termId,
		DocIdList: docIdList,
	}
	offset, err := dao.Index.Add(addModel)
	if err != nil {
		return 0, err
	}
	return offset, nil
}

func (i *Index) GetOne(beginOffset int64) (*model.Index, error) {
	return dao.Index.GetOne(beginOffset)
}
