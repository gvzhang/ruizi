package util

type TrieNode struct {
	data         byte
	children     map[byte]*TrieNode
	isEndingWord bool
}

type Trie struct {
	root *TrieNode
}

func (t *Trie) Insert(text []byte) {
	var ok bool
	p := t.root
	for k, v := range text {
		_, ok = p.children[v]
		if ok == false {
			nn := new(TrieNode)
			nn.data = text[k]
			nn.children = make(map[byte]*TrieNode, 0)
			p.children[v] = nn
		}
		p = p.children[v]
	}
	p.isEndingWord = true
}

func (t *Trie) Find(pattern []byte) bool {
	var ok bool
	p := t.root
	for _, v := range pattern {
		p, ok = p.children[v]
		if ok == false {
			return false
		}
	}
	if p.isEndingWord == false {
		return false
	}
	return true
}

func (t *Trie) Search(content []byte) ([]string, error) {
	var ok bool
	swl := make([]string, 0)
	for i := range content {
		j := i
		p := t.root
		swls := make([]byte, 0)
		for {
			swls = append(swls, content[j])
			p, ok = p.children[content[j]]
			if ok == false {
				break
			}
			if p.isEndingWord == true {
				swl = append(swl, string(swls))
				break
			}
			j++
		}
	}
	return swl, nil
}

func NewTrie() *Trie {
	tn := new(TrieNode)
	tn.data = '/'
	tn.children = make(map[byte]*TrieNode, 0)
	t := new(Trie)
	t.root = tn
	return t
}
