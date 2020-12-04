package util

import (
	"strings"
	"testing"
)

const HtmlBody string = `<html>
<head>
<meta charset="UTF-8"></meta>
<script type="text/javascript" src="https://cdn.bootcss.com/jquery/1.12.4/jquery.min.js"></script>
<style type="text/css"></style>
</head>
<body>
<iframe id="mf" width=1000 height=700 src="http://www.phptest56.local/"></iframe>
<input type="button" id="btn" value="测试" />
<script type="text/javascript">
	$("#btn").click(function(){
		var a = document.getElementById("mf").contentWindow.document.body.outerHTML;
		alert(a);
	});
</script>
<script>
	var _hmt = _hmt || [];
	(function() {
		var hm = document.createElement("script");
		hm.src = "https://hm.baidu.com/hm.js?082cf24f462676424893181d7123400e";
		var s = document.getElementsByTagName("script")[0]; 
		s.parentNode.insertBefore(hm, s);
	})();
</script>
<h1>段落1</h1>
<p id="article" style="font-size:red">正文...</p>
</body>
</html>`

func getHtmlBody() []byte {
	return []byte(HtmlBody)
}

func TestJsFilter(t *testing.T) {
	body := getHtmlBody()
	hf := NewHtmlFilter(body)
	hf.Js()
	isContainJs(hf, t)
}

func isContainJs(hf *htmlFilter, t *testing.T) {
	if strings.Contains(string(hf.body), "scripts") {
		t.Error("body contain js")
	}
}

func TestCssFilter(t *testing.T) {
	body := getHtmlBody()
	hf := NewHtmlFilter(body)
	hf.Js()
	isContainCss(hf, t)
}

func isContainCss(hf *htmlFilter, t *testing.T) {
	if strings.Contains(string(hf.body), "scripts") {
		t.Error("body contain css")
	}
}

func TestHtmlFilter(t *testing.T) {
	body := getHtmlBody()
	hf := NewHtmlFilter(body)
	hf.Html()
	isContainHtml(hf, t)
}

func isContainHtml(hf *htmlFilter, t *testing.T) {
	hfb := string(hf.body)
	for _, v := range hf.htmlTags {
		if strings.Contains(hfb, "<"+v+">") || strings.Contains(hfb, "<"+v+" ") || strings.Contains(hfb, "</"+v+">") {
			t.Error("body contain html " + v)
			t.Log(hfb)
		}
	}
}

func TestJsCssHtmlFilter(t *testing.T) {
	body := getHtmlBody()
	hf := NewHtmlFilter(body)
	hf.Js().Css().Html()
	isContainJs(hf, t)
	isContainCss(hf, t)
	isContainHtml(hf, t)
	t.Log(strings.Trim(string(hf.body), "\r\n"))
}
