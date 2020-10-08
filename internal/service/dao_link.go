package service

import (
	"ruizi/internal/dao"
	"ruizi/internal/model"
)

type DocLink struct {
}

func (d *DocLink) Add(docId uint64, url []byte) error {
	addModel := &model.DocLink{
		DocId: docId,
		Url:   url,
	}
	return dao.DocLink.Add(addModel)
}

func (d *DocLink) GetOne(docId uint64) (*model.DocLink, error) {
	return dao.DocLink.GetOne(docId)
}
