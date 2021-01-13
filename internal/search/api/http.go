package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ruizi/internal"
)

type SearchHttpServer struct {
	Router *gin.Engine
}

func (shs *SearchHttpServer) InitRouter() error {
	shs.Router = gin.New()
	// 将数值解析为json.Number的interface，而不是一个float64
	gin.EnableJsonDecoderUseNumber()
	shs.Router.LoadHTMLGlob(internal.GetConfig().Base.RootPath + "/templates/*")
	shs.Router.GET("/index.html", shs.index)
	return nil
}

func (shs *SearchHttpServer) index(c *gin.Context) {
	data := make([]uint64, 0)
	searchWord := c.PostForm("wd")
	if searchWord != "" {
		data = []uint64{1, 2, 3, 4}
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Ruizi fruit Search",
		"data":  data,
	})
}
