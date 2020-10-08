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

	config = &Config{
		base,
		link,
		bloomFilter,
		doc,
		docLink,
		docId,
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
