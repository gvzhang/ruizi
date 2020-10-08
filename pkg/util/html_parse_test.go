package util

import (
	"net/url"
	"testing"
)

func TestGetLinks(t *testing.T) {
	body := []byte("<body><div><a id=\"ta\" href=\"http://www.test.com/ac?id=123\">ta</a></div>" +
		"<span><a id=\"ta2\" href=\"https://www.rz.com/\">ta2</a></span></body>")
	htmlParse := NewHtmlParse(body)
	links, err := htmlParse.GetLinks()
	if err != nil {
		t.Error(err)
	}

	for _, l := range links {
		t.Log(l)
	}

	if len(links) != 2 {
		t.Errorf("error link count %d", len(links))
	}
}

const MainUrl = "https://www.test.com/goods/get/?id=123"

func TestJoinLink(t *testing.T) {
	parseUrl, err := url.Parse(MainUrl)
	if err != nil {
		t.Error(err)
	}
	mainHost := parseUrl.Scheme + "://" + parseUrl.Host
	mainHostPath := mainHost + parseUrl.Path

	join1 := "/a/b?c=1"
	link1, err := JoinLink(MainUrl, join1)
	if err != nil {
		t.Error(err)
	}
	expectLink1 := mainHost + join1
	if link1 != expectLink1 {
		t.Errorf("link1 join error. expect: %s actualy: %s", expectLink1, link1)
	}

	join2 := "./a/b?c=1"
	link2, err := JoinLink(MainUrl, join2)
	if err != nil {
		t.Error(err)
	}
	expectLink2 := mainHostPath[:len(mainHostPath)-1] + join2[1:]
	if link2 != expectLink2 {
		t.Errorf("link2 join error. expect: %s actualy: %s", expectLink2, link2)
	}

	join3 := "a/b?c=1"
	link3, err := JoinLink(MainUrl, join3)
	if err != nil {
		t.Error(err)
	}
	expectLink3 := mainHostPath + join3
	if link3 != expectLink3 {
		t.Errorf("link3 join error. expect: %s actualy: %s", expectLink3, link3)
	}

	join4 := "../b?c=1"
	link4, err := JoinLink(MainUrl, join4)
	if err != nil {
		t.Error(err)
	}
	expectLink4 := mainHost + "/goods" + join4[2:]
	if link4 != expectLink4 {
		t.Errorf("link4 join error. expect: %s actualy: %s", expectLink4, link4)
	}

	join5 := "../../b?c=1"
	link5, err := JoinLink(MainUrl, join5)
	if err != nil {
		t.Error(err)
	}
	expectLink5 := mainHost + join5[5:]
	if link5 != expectLink5 {
		t.Errorf("link5 join error. expect: %s actualy: %s", expectLink5, link5)
	}

	join6 := "../../../b?c=1"
	link6, err := JoinLink(MainUrl, join6)
	if err != nil {
		t.Error(err)
	}
	expectLink6 := mainHost + join6[8:]
	if link6 != expectLink6 {
		t.Errorf("link6 join error. expect: %s actualy: %s", expectLink6, link6)
	}

	mainUrl2 := "http://127.0.0.1:8080"
	join7 := "link3.html"
	link7, err := JoinLink(mainUrl2, join7)
	if err != nil {
		t.Error(err)
	}
	expectLink7 := mainUrl2 + "/" + join7
	if link7 != expectLink7 {
		t.Errorf("link7 join error. expect: %s actualy: %s", expectLink7, link7)
	}
}

func TestJoinLinkWithSuffix(t *testing.T) {
	mainUrl := "http://127.0.0.1:8080/link1.html"
	parseUrl, err := url.Parse(mainUrl)
	if err != nil {
		t.Error(err)
	}
	mainHost := parseUrl.Scheme + "://" + parseUrl.Host

	join := "link2.html"
	link, err := JoinLink(mainUrl, join)
	if err != nil {
		t.Error(err)
	}
	expectLink := mainHost + "/" + join
	if link != expectLink {
		t.Errorf("join error. expect: %s actualy: %s", expectLink, link)
	}
}

func TestJoinLinkWithWellName(t *testing.T) {
	mainUrl := "http://127.0.0.1:8080/link1.html"
	join := "#234ae"
	link, err := JoinLink(mainUrl, join)
	if err != nil {
		t.Error(err)
	}
	expectLink := mainUrl
	if link != expectLink {
		t.Errorf("join error. expect: %s actualy: %s", expectLink, link)
	}
}
