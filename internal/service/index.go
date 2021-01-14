package service

import (
	"ruizi/internal/dao"
	"ruizi/internal/model"
	"ruizi/pkg/logger"
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

func (i *Index) GetByTermId(termId uint64) (*model.Index, error) {
	offset := int64(0)
	for {
		indexModel, err := dao.Index.GetOne(offset)
		if err != nil {
			return nil, err
		}
		if indexModel == nil {
			return nil, nil
		}
		if indexModel.TermId == termId {
			logger.Sugar.Infof("Index.GetByTermId offset %d", offset)
			return indexModel, nil
		}
		offset = indexModel.NextOffset
	}
}
