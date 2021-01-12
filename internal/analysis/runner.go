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

		// todo tw单词列表存入单词编号文件（term_id.bin）
		// todo 单词id生成临时索引文件（term_id=>doc_id, tmp_index.bin）
		offset = docModel.NextOffset
	}

	// todo 临时索引文件生成倒排索引（term_id=>doc_id1, doc_id2...，index.bin）

	// todo 记录单词编号在索引文件中的偏移位置的文件（term_offset.bin）
	// 这个文件的作用是，帮助我们快速地查找某个单词编号在倒排索引中存储的位置，进而快速地从倒排索引中读取单词编号对应的网页编号列表。

	// 查询
	// 当用户在搜索框中，输入某个查询文本的时候，我们先对用户输入的文本进行分词处理。假设分词之后，我们得到 k 个单词。
	// 我们拿这 k 个单词，去 term_id.bin 对应的散列表中，查找对应的单词编号。经过这个查询之后，我们得到了这 k 个单词对应的单词编号。
	// 我们拿这 k 个单词编号，去 term_offset.bin 对应的散列表中，查找每个单词编号在倒排索引文件中的偏移位置。经过这个查询之后，我们得到了 k 个偏移位置。
	// 我们拿这 k 个偏移位置，去倒排索引（index.bin）中，查找 k 个单词对应的包含它的网页编号列表。经过这一步查询之后，我们得到了 k 个网页编号列表。
	// 我们针对这 k 个网页编号列表，统计每个网页编号出现的次数。具体到实现层面，我们可以借助散列表来进行统计。
	// 统计得到的结果，我们按照出现次数的多少，从小到大排序。出现次数越多，说明包含越多的用户查询单词（用户输入的搜索文本，经过分词之后的单词）。
}

func (r *Runner) Stop() error {
	return nil
}
