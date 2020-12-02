package util

import (
	"bytes"
	"sync"
)

type htmlFilter struct {
	body     []byte
	jsTag    string
	cssTag   string
	htmlTags []string
	lock     *sync.RWMutex
}

func NewHtmlFilter(b []byte) *htmlFilter {
	hf := &htmlFilter{
		body: b,
	}
	hf.jsTag = "script"
	hf.cssTag = "style"
	hf.htmlTags = []string{"!DOCTYPE", "a", "abbr", "acronym", "address", "applet", "area", "article", "aside",
		"audio", "b", "base", "basefont", "bdi", "bdo", "big", "blockquote", "body", "br", "button", "canvas",
		"caption", "center", "cite", "code", "col", "colgroup", "command", "datalist", "dd", "del", "details",
		"dir", "details", "dir", "div", "dfn", "dialog", "dl", "dt", "em", "embed", "fieldset", "figcaption",
		"figure", "font", "footer", "form", "frame", "frameset", "h1", "h2", "h3", "h4", "h5", "h6", "head",
		"header", "hr", "html", "i", "iframe", "img", "input", "input", "ins", "kdb", "keygen", "label",
		"legend", "li", "link", "map", "mark", "menu", "menuitem", "meta", "meter", "nav", "noframes",
		"noscript", "object", "ol", "optgroup", "option", "output", "p", "param", "pre", "progress", "q",
		"rp", "rt", "ruby", "s", "samp", "section", "select", "small", "source", "span", "strike", "strong",
		"sub", "summary", "sup", "table", "tbody", "td", "textarea", "tfoot", "th", "thead", "time", "title",
		"tr", "track", "tt", "u", "ul", "var", "video", "wbr"}
	hf.lock = new(sync.RWMutex)
	return hf
}

func (hf *htmlFilter) tagEndPos(tag []byte, content []byte, sp int) int {
	tl := len(tag)
	cl := len(content)
	pj := tag[0 : tl-1]
	if sp+tl < cl && bytes.Equal(content[sp:sp+tl-1], pj) {
		for x := sp; x < cl; x++ {
			if content[x] == tag[tl-1] {
				return x
			}
		}
	}
	return 0
}

func (hf *htmlFilter) filterTagWithContent(tag string) *htmlFilter {
	hf.lock.Lock()
	defer hf.lock.Unlock()

	jtb, jte := []byte("<"+tag+">"), []byte("</"+tag+">")
	jel, bl := len(jte), len(hf.body)
	cb := hf.body

	i, j := 0, 0
	for ; j < bl; i, j = i+1, j+1 {
		mb := hf.tagEndPos(jtb, cb, j)
		if mb != 0 {
			j = mb
			for ; j < bl-1; j++ {
				if j+jel >= bl {
					continue
				}
				if bytes.Equal(cb[j:j+jel], jte) {
					// +1换行符
					j = j + jel + 1
					break
				}
			}
		}
		cb[i] = cb[j]
	}
	hf.body = cb[0:i]
	return hf
}

func (hf *htmlFilter) Js() *htmlFilter {
	return hf.filterTagWithContent(hf.jsTag)
}

func (hf *htmlFilter) Css() *htmlFilter {
	return hf.filterTagWithContent(hf.cssTag)
}

func (hf *htmlFilter) Html() *htmlFilter {
	hf.lock.Lock()
	defer hf.lock.Unlock()

	cb := hf.body
	bl := len(hf.body)
	hTags := make([][]byte, len(hf.htmlTags)*2)
	var k int
	for _, h := range hf.htmlTags {
		hTags[k] = []byte("<" + h + ">")
		hTags[k+1] = []byte("</" + h + ">")
		k = k + 2
	}
	for _, t := range hTags {
		i, j := 0, 0
		for ; j < bl; i, j = i+1, j+1 {
			mb := hf.tagEndPos(t, cb, j)
			if mb != 0 {
				j = mb
			}
			cb[i] = cb[j]
		}
		cb = cb[0:i]
	}
	hf.body = cb
	return hf
}
