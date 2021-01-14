package api

import (
	"fmt"
	"net/http"
	"ruizi/internal"
	"ruizi/internal/service"
	"ruizi/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type SearchHttpServer struct {
	Router *gin.Engine
}

func (shs *SearchHttpServer) InitRouter() error {
	shs.Router = gin.New()
	// 将数值解析为json.Number的interface，而不是一个float64
	gin.EnableJsonDecoderUseNumber()
	shs.Router.LoadHTMLGlob(internal.GetConfig().Base.RootPath + "/templates/*")
	shs.Router.Any("/index.html", shs.index)
	shs.Router.GET("/page/raw", shs.rawPage)
	shs.Router.GET("/search/words", shs.searchWord)
	shs.Router.GET("/search/index", shs.searchIndex)
	return nil
}

type doc struct {
	DocId uint64
	Url   string
}

func (shs *SearchHttpServer) index(c *gin.Context) {
	ginH := gin.H{}
	ginH["title"] = "Ruizi fruit Search"

	searchWord := c.PostForm("wd")
	ginH["wd"] = searchWord
	if searchWord != "" {
		docs, err := doSearch(searchWord)
		ginH["Data"] = docs
		if err != nil {
			c.HTML(http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}

	c.HTML(http.StatusOK, "index.tmpl", ginH)
}

// 当用户在搜索框中，输入某个查询文本的时候，我们先对用户输入的文本进行分词处理。假设分词之后，我们得到 k 个单词。
// 我们拿这 k 个单词，去 term_id.bin 对应的散列表中，查找对应的单词编号。经过这个查询之后，我们得到了这 k 个单词对应的单词编号。
// 我们拿这 k 个单词编号，去 term_offset.bin 对应的散列表中，查找每个单词编号在倒排索引文件中的偏移位置。经过这个查询之后，我们得到了 k 个偏移位置。
// 我们拿这 k 个偏移位置，去倒排索引（index.bin）中，查找 k 个单词对应的包含它的网页编号列表。经过这一步查询之后，我们得到了 k 个网页编号列表。
// 我们针对这 k 个网页编号列表，统计每个网页编号出现的次数。具体到实现层面，我们可以借助散列表来进行统计。
// 统计得到的结果，我们按照出现次数的多少，从小到大排序。出现次数越多，说明包含越多的用户查询单词（用户输入的搜索文本，经过分词之后的单词）。
func doSearch(searchWord string) ([]doc, error) {
	docs := make([]doc, 0)
	termService := new(service.Term)
	term, err := termService.GetByWord([]byte(searchWord))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find term error by %s", searchWord))
	}
	if term == nil {
		return docs, nil
	}
	logger.Sugar.Infof("doSearch1 TermId %d", term.Id)

	termOffsetService := new(service.TermOffset)
	termOffset, err := termOffsetService.GetByTermId(term.Id)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find termOffset error by %d", term.Id))
	}
	if termOffset == nil {
		return docs, nil
	}
	logger.Sugar.Infof("doSearch2 TermOffset %d", termOffset.Offset)

	indexService := new(service.Index)
	index, err := indexService.GetOne(termOffset.Offset)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find index error by %d", termOffset.Offset))
	}
	logger.Sugar.Infof("doSearch3 Index %+v", index)

	docLinkModel := new(service.DocLink)
	for _, v := range index.DocIdList {
		docLink, err := docLinkModel.GetOne(v)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("find doc_link error by %d", v))
		}
		docs = append(docs, doc{DocId: docLink.DocId, Url: string(docLink.Url)})
	}
	return docs, nil
}

func (shs *SearchHttpServer) rawPage(c *gin.Context) {
	ginH := gin.H{}
	ginH["title"] = "Ruizi fruit Search"

	qid, _ := c.GetQuery("id")
	docId, err := strconv.Atoi(qid)
	if err != nil {
		c.HTML(http.StatusInternalServerError, err.Error(), nil)
		return
	}
	docService := new(service.Doc)
	doc, err := docService.GetById(uint64(docId))
	if err != nil {
		c.HTML(http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if doc != nil {
		docLinkModel := new(service.DocLink)
		docLink, err := docLinkModel.GetOne(doc.Id)
		if err != nil {
			c.HTML(http.StatusInternalServerError, err.Error(), nil)
			return
		}
		ginH["raw"] = string(doc.Raw)
		ginH["link"] = string(docLink.Url)
	}

	c.HTML(http.StatusOK, "raw.tmpl", ginH)
}

type term struct {
	Id  uint64
	Txt string
}

func (shs *SearchHttpServer) searchWord(c *gin.Context) {
	termService := new(service.Term)
	terms, err := termService.GetAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, err.Error(), nil)
		return
	}
	data := make([]term, 0)
	for _, tm := range terms {
		data = append(data, term{
			Id:  tm.Id,
			Txt: string(tm.Txt),
		})
	}
	c.HTML(http.StatusOK, "debug.tmpl", gin.H{
		"Data": data,
	})
}

func (shs *SearchHttpServer) searchIndex(c *gin.Context) {
	tid, _ := c.GetQuery("id")
	termId, err := strconv.Atoi(tid)
	if err != nil {
		c.HTML(http.StatusInternalServerError, err.Error(), nil)
		return
	}
	indexService := new(service.Index)
	index, err := indexService.GetByTermId(uint64(termId))
	if err != nil {
		c.HTML(http.StatusInternalServerError, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, index.DocIdList)
}
