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


	// 查询
	// 当用户在搜索框中，输入某个查询文本的时候，我们先对用户输入的文本进行分词处理。假设分词之后，我们得到 k 个单词。
	// 我们拿这 k 个单词，去 term_id.bin 对应的散列表中，查找对应的单词编号。经过这个查询之后，我们得到了这 k 个单词对应的单词编号。
	// 我们拿这 k 个单词编号，去 term_offset.bin 对应的散列表中，查找每个单词编号在倒排索引文件中的偏移位置。经过这个查询之后，我们得到了 k 个偏移位置。
	// 我们拿这 k 个偏移位置，去倒排索引（index.bin）中，查找 k 个单词对应的包含它的网页编号列表。经过这一步查询之后，我们得到了 k 个网页编号列表。
	// 我们针对这 k 个网页编号列表，统计每个网页编号出现的次数。具体到实现层面，我们可以借助散列表来进行统计。
	// 统计得到的结果，我们按照出现次数的多少，从小到大排序。出现次数越多，说明包含越多的用户查询单词（用户输入的搜索文本，经过分词之后的单词）。
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
