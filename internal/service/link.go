package service

import (
	"ruizi/internal/dao"
	"ruizi/internal/model"
)

type Link struct {
}

func (l *Link) Add(url []byte) error {
	addModel := &model.Link{
		Status: model.LinkStatusWait,
		Url:    url,
	}
	return dao.Link.Add(addModel)
}

func (l *Link) GetOne(beginOffset int64) (*model.Link, error) {
	return dao.Link.GetOne(beginOffset)
}

func (l *Link) FinishCrawler(linkModel *model.Link) error {
	err := dao.Link.UpdateStatus(linkModel, model.LinkStatusDone)
	if err != nil {
		return err
	}
	return nil
}
