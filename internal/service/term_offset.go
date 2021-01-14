package service

import (
	"ruizi/internal/dao"
	"ruizi/internal/model"
)

type TermOffset struct {
}

func (to *TermOffset) Add(termId uint64, offset int64) error {
	addModel := &model.TermOffset{
		TermId: termId,
		Offset: offset,
	}
	err := dao.TermOffset.Add(addModel)
	if err != nil {
		return err
	}
	return nil
}

func (to *TermOffset) GetOne(beginOffset int64) (*model.TermOffset, error) {
	return dao.TermOffset.GetOne(beginOffset)
}

func (to *TermOffset) GetByTermId(termId uint64) (*model.TermOffset, error) {
	offset := int64(0)
	for {
		termOffsetModel, err := dao.TermOffset.GetOne(offset)
		if err != nil {
			return nil, err
		}
		if termOffsetModel == nil {
			return nil, nil
		}
		if termOffsetModel.TermId == termId {
			return termOffsetModel, nil
		}
		offset = termOffsetModel.NextOffset
	}
}
