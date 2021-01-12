package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"ruizi/pkg/logger"
	"strconv"

	"go.uber.org/zap/zapcore"
	"gopkg.in/ini.v1"
)

var (
	configIni *ini.File
	config    *Config
)

func GetConfig() *Config {
	return config
}

type Config struct {
	Base
	Link
	BloomFilter
	Doc
	DocLink
	DocId
	WordLib
	TermId
	Term
	TmpIndex
	Index
	TermOffset
}

type Base struct {
	Env      string
	RootPath string
}

type Link struct {
	DataPath string
}

type BloomFilter struct {
	DataPath string
}

type Doc struct {
	DataPath string
}

type DocLink struct {
	DataPath string
}

type DocId struct {
	DataPath string
}

type WordLib struct {
	DataPath string
}

type Term struct {
	DataPath string
}

type TermId struct {
	DataPath string
}

type TmpIndex struct {
	DataPath string
}

type Index struct {
	DataPath string
}

type TermOffset struct {
	DataPath string
}

func InitConfig() {
	env := os.Getenv("RUIZI_ENV")
	rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("get root dir fail:" + err.Error())
	}
	configFile := rootPath + "/config.ini"
	fmt.Println("Load Config File " + configFile)
	configIni, err = ini.Load(configFile)
	if err != nil {
		panic(err.Error())
	}

	base := Base{
		Env:      env,
		RootPath: rootPath,
	}

	link := Link{
		DataPath: GetRaw("data_path_link"),
	}

	bloomFilter := BloomFilter{
		DataPath: GetRaw("data_path_bloom_filter"),
	}

	doc := Doc{
		DataPath: GetRaw("data_path_doc"),
	}

	docLink := DocLink{
		DataPath: GetRaw("data_path_doc_link"),
	}

	docId := DocId{
		DataPath: GetRaw("data_path_doc_id"),
	}

	wordLib := WordLib{
		DataPath: GetRaw("data_path_word_lib"),
	}

	termId := TermId{
		DataPath: GetRaw("data_path_term_id"),
	}

	term := Term{
		DataPath: GetRaw("data_path_term"),
	}

	tmpIndex := TmpIndex{
		DataPath: GetRaw("data_path_tmp_index"),
	}

	index := Index{
		DataPath: GetRaw("data_path_index"),
	}

	termOffset := TermOffset{
		DataPath: GetRaw("data_path_term_offset"),
	}

	config = &Config{
		base,
		link,
		bloomFilter,
		doc,
		docLink,
		docId,
		wordLib,
		termId,
		term,
		tmpIndex,
		index,
		termOffset,
	}

	logLevel, err := strconv.Atoi(GetRaw("log_level"))
	if err != nil {
		panic(err.Error())
	}
	logger.Leavel = zapcore.Level(logLevel)
	logger.Target = GetRaw("log_target")
}

func GetRaw(key string) string {
	return configIni.Section("").Key(key).String()
}
