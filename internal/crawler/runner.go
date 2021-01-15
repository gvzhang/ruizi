package crawler

import (
	"fmt"
	"math/rand"
	"net/url"
	"ruizi/internal"
	"ruizi/internal/model"
	"ruizi/internal/service"
	"ruizi/pkg/logger"
	"ruizi/pkg/util"
	"time"
)

type LinkError struct {
	url []byte
	err error
}

func (l *LinkError) Error() string {
	return fmt.Sprintf("%s crawler error %s", l.url, l.err)
}

type Runner struct {
	linkService    *service.Link
	docService     *service.Doc
	docLinkService *service.DocLink
}

func NewRunner() *Runner {
	r := new(Runner)
	r.linkService = new(service.Link)
	r.docService = new(service.Doc)
	r.docLinkService = new(service.DocLink)
	return r
}

func (r *Runner) Start() error {
	// 初始化bloom
	logger.Sugar.Info("bloom init")
	bloomSavePath := internal.GetConfig().BloomFilter.DataPath
	seeds := []int8{4, 9, 16, 22, 31}
	bloomData, err := util.BloomFileData(bloomSavePath)
	if err != nil {
		return err
	}
	bloom := util.NewBloom(1<<24, seeds)
	if len(bloomData) != 0 {
		logger.Sugar.Infof("ImportDB %d", len(bloomData))
		bloom.ImportDB(bloomData)
	}

	offset := int64(0)
	// todo 事务原子性问题, 并发提高处理速度
	for {
		logger.Sugar.Infof("crawler get link from %d", offset)
		linkModel, err := r.linkService.GetOne(offset)
		if err != nil {
			return err
		}

		if linkModel == nil {
			return nil
		}
		offset = linkModel.NextOffset

		if linkModel.Status == model.LinkStatusDone {
			logger.Sugar.Infof("%s crawler is done", linkModel.Url)
			continue
		}

		lu := linkModel.Url
		le := new(LinkError)
		le.url = lu

		// 布隆过滤器过滤爬取过的url
		logger.Sugar.Infof("%s crawler check bloom", lu)
		exists, err := bloom.Get(string(lu))
		if err != nil {
			le.err = err
			return le
		}
		// 会有误差，但可接受
		if exists == true {
			logger.Sugar.Infof("bloom find url %s", lu)
			err = r.linkService.FinishCrawler(linkModel)
			if err != nil {
				le.err = err
				return le
			}
			continue
		}

		// 开始爬取
		logger.Sugar.Infof("%s crawler begin", lu)
		body, err := util.RetryGet(string(lu))
		if err != nil {
			logger.Sugar.Infof("%s can not crawler %s", lu, err.Error())
			err = r.linkService.FinishCrawler(linkModel)
			if err != nil {
				le.err = err
				return le
			}
		}

		// 保存原始网页
		logger.Sugar.Infof("%s crawler save doc", lu)
		docId, err := r.docService.Add(body)
		if err != nil {
			le.err = err
			return le
		}

		// 保存网页id对应链接
		logger.Sugar.Infof("%s crawler save dock_link", lu)
		err = r.docLinkService.Add(docId, lu)
		if err != nil {
			le.err = err
			return le
		}

		// 分析url更新待爬库
		err = r.processUrl(string(lu), body)
		if err != nil {
			le.err = err
			return le
		}

		// 爬虫完成，更新布隆过滤器
		logger.Sugar.Infof("%s crawler update bloom", lu)
		err = bloom.Set(string(lu))
		if err != nil {
			le.err = err
			return le
		}

		// bloom持久化
		logger.Sugar.Infof("%s crawler persistence bloom", lu)
		_ = util.BloomPersistence(bloom.OutputDB(), bloomSavePath)

		// 更新待爬数据状态
		logger.Sugar.Infof("%s crawler finish", lu)
		err = r.linkService.FinishCrawler(linkModel)
		if err != nil {
			le.err = err
			return le
		}

		wait := rand.Intn(1000) + 500
		logger.Sugar.Infof("%s crawler sleep %d", lu, wait)
		time.Sleep(time.Duration(wait) * time.Millisecond)
	}
}

// 分析url更新待爬库，广度优先爬取
func (r *Runner) processUrl(lu string, body []byte) error {
	logger.Sugar.Infof("%s crawler html parse", lu)
	htmlParse := util.NewHtmlParse(body)
	bodyLinks, err := htmlParse.GetLinks()
	if err != nil {
		return err
	}
	if len(bodyLinks) != 0 {
		// 保存待爬库
		logger.Sugar.Infof("%s crawler save links %d", lu, len(bodyLinks))
		for _, ln := range bodyLinks {
			jln, err := util.JoinLink(lu, ln)
			if err != nil {
				return err
			}
			cjln, err := checkJoinUrl(lu, jln)
			if err != nil {
				return err
			}
			if cjln == false {
				continue
			}
			err = r.linkService.Add([]byte(jln))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 检查是否爬虫地址
func checkJoinUrl(mainUrl string, joinUrl string) (bool, error) {
	if joinUrl == "" {
		return false, nil
	}

	jp, err := url.Parse(joinUrl)
	if err != nil {
		return false, err
	}
	if jp.Host == "" {
		return false, nil
	}

	mp, err := url.Parse(mainUrl)
	if err != nil {
		return false, err
	}
	if mp.Host == "" {
		return false, nil
	}

	jpHost := jp.Host
	if len(jpHost) > 4 && jpHost[:4] == "www." {
		jpHost = jpHost[4:]
	}

	mpHost := mp.Host
	if len(mpHost) > 4 && mpHost[:4] == "www." {
		mpHost = mpHost[4:]
	}

	return jpHost == mpHost, nil
}

func (r *Runner) Stop() error {
	return nil
}
