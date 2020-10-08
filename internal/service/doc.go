package service

import (
	"ruizi/internal/dao"
	"ruizi/internal/model"
)

type Doc struct {
}

func (d *Doc) Add(raw []byte) (uint64, error) {
	addModel := &model.Doc{
		Status: model.DocStatusWait,
		Raw:    raw,
	}
	err := dao.Doc.Add(addModel)
	if err != nil {
		return 0, err
	}
	return addModel.Id, nil
}

func (d *Doc) GetOne(beginOffset int64) (*model.Doc, error) {
	return dao.Doc.GetOne(beginOffset)
}

func (d *Doc) FinishAnalysis(docModel *model.Doc) error {
	err := dao.Doc.UpdateStatus(docModel, model.DocStatusAnalysis)
	if err != nil {
		return err
	}
	return nil
}
