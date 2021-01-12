package service

import (
	"ruizi/internal/dao"
	"ruizi/internal/model"
)

type Term struct {
}

func (t *Term) Add(txt []byte) (uint64, error) {
	addModel := &model.Term{
		Status: model.TermStatusEnable,
		Txt:    txt,
	}
	err := dao.Term.Add(addModel)
	if err != nil {
		return 0, err
	}
	return addModel.Id, nil
}

func (t *Term) GetOne(beginOffset int64) (*model.Term, error) {
	return dao.Term.GetOne(beginOffset)
}
