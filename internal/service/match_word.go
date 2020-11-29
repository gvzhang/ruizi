package service

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"ruizi/internal"
	"ruizi/pkg/util"

	"github.com/pkg/errors"
)

type MatchWord struct {
	trie *util.Trie
}

func NewMatchWord() (*MatchWord, error) {
	t := util.NewTrie()
	dataPath := internal.GetConfig().WordLib.DataPath
	files, err := ioutil.ReadDir(dataPath)
	if err != nil {
		return nil, errors.Wrap(err, "MatchWord ini error")
	}
	for _, f := range files {
		fo, err := os.Open(path.Join(dataPath, f.Name()))
		if err != nil {
			return nil, errors.Wrap(err, "open "+f.Name()+" fail")
		}
		br := bufio.NewReader(fo)
		sc := bufio.NewScanner(br)
		sc.Split(bufio.ScanLines)
		for sc.Scan() {
			t.Insert(sc.Bytes())
		}
	}

	mw := new(MatchWord)
	mw.trie = t
	return mw, nil
}

func (mw *MatchWord) Search(content []byte) ([]string, error) {
	return mw.trie.Search(content)
}
