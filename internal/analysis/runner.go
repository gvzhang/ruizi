package analysis

import (
	"ruizi/internal/dao"
	"ruizi/internal/service"
	"ruizi/pkg/util"
)

type Runner struct {
	termService       *service.Term
	tmpIndexService   *service.TmpIndex
	indexService      *service.Index
	termOffsetService *service.TermOffset
	terms             map[string]uint64
}

func NewRunner() (*Runner, error) {
	r := new(Runner)
	r.termService = new(service.Term)
	r.tmpIndexService = new(service.TmpIndex)
	r.indexService = new(service.Index)
	r.termOffsetService = new(service.TermOffset)
	r.terms = make(map[string]uint64)
	err := r.initTerms()
	return r, err
}

func (r *Runner) Start() error {
	var err error
	err = r.MatchWords()
	if err != nil {
		return err
	}

	err = r.MakeIndexBin()
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) MatchWords() error {
	offset := int64(0)

	mw, err := service.NewMatchWord()
	if err != nil {
		return err
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

		// tw单词列表存入单词编号文件（term_id.bin）
		var tid uint64
		var ok bool
		for _, word := range tw {
			tid, ok = r.terms[word]
			if ok == false {
				tid, err = r.termService.Add([]byte(word))
				if err != nil {
					return err
				}
				r.terms[word] = tid
			}

			// 单词id生成临时索引文件（term_id=>doc_id, tmp_index.bin）
			err = r.tmpIndexService.Add(tid, docModel.Id)
			if err != nil {
				return err
			}
		}

		offset = docModel.NextOffset
	}
}

// 临时索引文件生成倒排索引（term_id=>doc_id1, doc_id2...，index.bin）
func (r *Runner) MakeIndexBin() error {
	tmpIndexData, err := r.tmpIndexService.GetAll(0)
	if err != nil {
		return err
	}
	tmpIndexData = r.tmpIndexService.Sort(tmpIndexData)
	indexData := make(map[uint64][]uint64)
	for _, tmpIndex := range tmpIndexData {
		indexData[tmpIndex.TermId] = append(indexData[tmpIndex.TermId], tmpIndex.DocId)
	}
	for termId, docIdList := range indexData {
		offset, err := r.indexService.Add(termId, docIdList)
		if err != nil {
			return err
		}

		// 记录单词编号在索引文件中的偏移位置的文件（term_offset.bin）
		// 这个文件的作用是，帮助我们快速地查找某个单词编号在倒排索引中存储的位置，进而快速地从倒排索引中读取单词编号对应的网页编号列表。
		err = r.termOffsetService.Add(termId, offset)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) initTerms() error {
	offset := int64(0)
	for {
		term, err := r.termService.GetOne(offset)
		if err != nil {
			return err
		}
		if term == nil {
			return nil
		}
		r.terms[string(term.Txt)] = term.Id
		offset = term.NextOffset
	}
}

func (r *Runner) Stop() error {
	return nil
}
