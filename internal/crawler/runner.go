package crawler

import (
	"math/rand"
	"net/url"
	"ruizi/internal"
	"ruizi/internal/model"
	"ruizi/internal/service"
	"ruizi/pkg/logger"
	"ruizi/pkg/util"
	"time"
)

type Runner struct {
}

func (r *Runner) Start() error {
	var err error
	var linkModel *model.Link
	var exists, cjln bool
	var docId uint64
	var body []byte
	var bodyLinks []string
	var jln string

	linkService := &service.Link{}
	docService := &service.Doc{}
	docLinkService := &service.DocLink{}
	bloomSavePath := internal.GetConfig().BloomFilter.DataPath

	// 初始化bloom
	logger.Sugar.Info("bloom init")
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
bf:
	// todo 事务原子性问题
	for {
		logger.Sugar.Infof("crawler get link from %d", offset)
		linkModel, err = linkService.GetOne(offset)
		if err != nil {
			break bf
		}

		if linkModel == nil {
			break bf
		}
		offset = linkModel.NextOffset

		if linkModel.Status == model.LinkStatusDone {
			logger.Sugar.Infof("%s crawler is done", linkModel.Url)
			continue
		}

		url := linkModel.Url

		// 布隆过滤器过滤爬取过的url
		logger.Sugar.Infof("%s crawler check bloom", url)
		exists, err = bloom.Get(string(url))
		if err != nil {
			break bf
		}
		// 会有误差，但可接受
		if exists == true {
			logger.Sugar.Infof("bloom find url %s", url)
			err = linkService.FinishCrawler(linkModel)
			if err != nil {
				break bf
			}
			continue
		}

		// 开始爬取
		logger.Sugar.Infof("%s crawler begin", url)
		body, err = util.RetryGet(string(url))
		if err != nil {
			logger.Sugar.Infof("%s can not crawler %s", url, err.Error())
			err = linkService.FinishCrawler(linkModel)
			if err != nil {
				break bf
			}
		}

		// 保存原始网页
		logger.Sugar.Infof("%s crawler save doc", url)
		docId, err = docService.Add(body)
		if err != nil {
			break bf
		}

		// 保存网页id对应链接
		logger.Sugar.Infof("%s crawler save dock_link", url)
		err = docLinkService.Add(docId, url)
		if err != nil {
			break bf
		}

		// 分析url
		logger.Sugar.Infof("%s crawler html parse", url)
		htmlParse := util.NewHtmlParse(body)
		bodyLinks, err = htmlParse.GetLinks()
		if err != nil {
			break bf
		}
		if len(bodyLinks) != 0 {
			// 保存待爬库
			logger.Sugar.Infof("%s crawler save links %d", url, len(bodyLinks))
			for _, ln := range bodyLinks {
				jln, err = util.JoinLink(string(url), ln)
				if err != nil {
					break bf
				}
				cjln, err = checkJoinUrl(string(url), jln)
				if err != nil {
					break bf
				}
				if cjln == false {
					continue
				}
				err = linkService.Add([]byte(jln))
				if err != nil {
					break bf
				}
			}
		}

		// 爬虫完成，更新布隆过滤器
		logger.Sugar.Infof("%s crawler update bloom", url)
		err = bloom.Set(string(url))
		if err != nil {
			break bf
		}

		// bloom持久化
		logger.Sugar.Infof("%s crawler persistence bloom", url)
		util.BloomPersistence(bloom.OutputDB(), bloomSavePath)

		// 更新待爬数据状态
		logger.Sugar.Infof("%s crawler finish", url)
		err = linkService.FinishCrawler(linkModel)
		if err != nil {
			break bf
		}

		sleepSec := rand.Intn(1000) + 500
		logger.Sugar.Infof("%s crawler sleep %d", url, sleepSec)
		time.Sleep(time.Duration(sleepSec) * time.Millisecond)
	}

	if linkModel != nil && err != nil {
		logger.Sugar.Infof("%s crawler error %w", linkModel.Url, err)
	}
	return err
}

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
