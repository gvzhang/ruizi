package service

import (
	"ruizi/internal/dao"
	"ruizi/internal/model"
	"strings"
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

func (t *Term) GetByWord(word []byte) (*model.Term, error) {
	offset := int64(0)
	for {
		termModel, err := dao.Term.GetOne(offset)
		if err != nil {
			return nil, err
		}
		if termModel == nil {
			return nil, nil
		}
		if strings.Compare(string(termModel.Txt), string(word)) == 0 {
			return termModel, nil
		}
		offset = termModel.NextOffset
	}
}

func (t *Term) GetAll() ([]*model.Term, error) {
	result := make([]*model.Term, 0)
	offset := int64(0)
	for {
		termModel, err := dao.Term.GetOne(offset)
		if err != nil {
			return nil, err
		}
		if termModel == nil {
			break
		}
		result = append(result, termModel)
		offset = termModel.NextOffset
	}
	return result, nil
}
