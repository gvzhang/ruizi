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
<h1>段落1</h1>
<p id="article" style="font-size:red">正文...</p>
</body>
</html>`

func TestJsFilter(t *testing.T) {
	HtmlBodyB := []byte(HtmlBody)
	hf := NewHtmlFilter(HtmlBodyB)
	hf.Js()
	if strings.Contains(string(hf.body), "scripts") {
		t.Error("body contain js")
	}
}

func TestCssFilter(t *testing.T) {
	HtmlBodyB := []byte(HtmlBody)
	hf := NewHtmlFilter(HtmlBodyB)
	hf.Js()
	if strings.Contains(string(hf.body), "scripts") {
		t.Error("body contain css")
	}
}

func TestHtmlFilter(t *testing.T) {
	HtmlBodyB := []byte(HtmlBody)
	hf := NewHtmlFilter(HtmlBodyB)
	hf.Html()
	t.Log(string(hf.body))
}

func TestJsCssHtmlFilter(t *testing.T) {

}
