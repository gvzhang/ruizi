package analysis

import (
	"fmt"
	"ruizi/internal/dao"
	"ruizi/internal/service"
)

type Runner struct {
}

func NewRunner() *Runner {
	r := new(Runner)
	return r
}

func (r *Runner) Start() error {
	offset := int64(0)

	mw, err := service.NewMatchWord()
	if err != nil {
		panic(err)
	}
	
	for {
		docModel, err := dao.Doc.GetOne(offset)
		if err != nil {
			return err
		}
		if docModel == nil {
			return nil
		}
		tw, err := mw.Search(docModel.Raw)
		if err != nil {
			return err
		}
		fmt.Println(docModel.Id, tw)
		offset = docModel.NextOffset
	}

	return nil
}

func (r *Runner) Stop() error {
	return nil
}
