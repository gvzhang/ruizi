package util

import (
	"net/url"
	"path/filepath"
	"strings"
)

type htmlParse struct {
	body []byte
}

func NewHtmlParse(b []byte) *htmlParse {
	hp := &htmlParse{
		body: b,
	}
	return hp
}

func (hp *htmlParse) GetLinks() ([]string, error) {
	links := make([]string, 0)

	// 使用BF算法，也就是暴力匹配;-)
	fa := false
	fas := false
	var link []byte
	bodyLen := len(hp.body)
	for i := 0; i < bodyLen; i++ {
		if fas == true {
			if hp.body[i] == '"' || hp.body[i] == '\'' {
				links = append(links, string(link))
				fa = false
				fas = false
				link = link[0:0]
				continue
			}
			link = append(link, hp.body[i])
			continue
		}

		if fa == true && (bodyLen >= i+5) && strings.Compare(string(hp.body[i:i+5]), "href=") == 0 {
			fas = true
			i = i + 5
			continue
		}

		if (bodyLen >= i+2) && strings.Compare(string(hp.body[i:i+2]), "<a") == 0 {
			fa = true
			i = i + 1
			continue
		}
	}

	return links, nil
}

func JoinLink(mainUrl string, joinUrl string) (string, error) {
	jp, err := url.Parse(joinUrl)
	if err != nil {
		return "", err
	}
	if jp.Scheme != "" {
		if (jp.Scheme == "http") || (jp.Scheme == "https") {
			return joinUrl, nil
		} else {
			return "", nil
		}
	}

	parseUrl, err := url.Parse(mainUrl)
	if err != nil {
		return "", err
	}

	resLink := parseUrl.Scheme + "://" + parseUrl.Host
	if joinUrl == "" {
		return resLink, nil
	}
	if joinUrl[:1] == "#" {
		return mainUrl, nil
	}

	if joinUrl[:1] != "/" {
		urlPath := parseUrl.Path
		if urlPath == "" {
			urlPath = "/"
		}
		if filepath.Ext(urlPath) != "" {
			urlPath = filepath.Dir(urlPath)
		}

		if len(joinUrl) > 1 && joinUrl[:2] == "./" {
			joinUrl = joinUrl[2:]
		}

		if len(joinUrl) > 2 && joinUrl[:3] == "../" {
			// 处理层级问题
			if urlPath[len(urlPath)-1:] == "/" {
				urlPath = filepath.Dir(urlPath)
			}
			for joinUrl[:3] == "../" {
				urlPath = filepath.Dir(urlPath)
				joinUrl = joinUrl[3:]
			}
		}

		if urlPath[len(urlPath)-1:] != "/" {
			urlPath += "/"
		}
		resLink += urlPath
	}

	resLink += joinUrl
	return resLink, nil
}
