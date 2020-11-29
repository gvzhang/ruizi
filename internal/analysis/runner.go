package analysis

import (
	"fmt"
	"ruizi/internal/dao"
	"ruizi/internal/service"
)

type Runner struct {
}

func (r *Runner) Start() {
	offset := int64(0)
	docModel, err := dao.Doc.GetOne(offset)
	if err != nil {
		panic(err)
	}
	mw, err := service.NewMatchWord()
	if err != nil {
		panic(err)
	}
	tw, err := mw.Search(docModel.Raw)
	if err != nil {
		panic(err)
	}
	fmt.Println(tw)
}

func (r *Runner) Stop() error {
	return nil
}
