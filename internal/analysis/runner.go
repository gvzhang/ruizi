package analysis

import (
	"ruizi/internal/dao"
	"ruizi/internal/service"
	"ruizi/pkg/util"
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

		// 过滤JS、CSS、HTML标签（需优化性能）
		hf := util.NewHtmlFilter(docModel.Raw)
		hf.Css().Js().Html()
		docRaw := hf.GetBody()

		// 文本分词
		tw, err := mw.Search(docRaw)
		if err != nil {
			return err
		}
		tw = util.Uniq(tw)

		offset = docModel.NextOffset
	}
}

func (r *Runner) Stop() error {
	return nil
}
