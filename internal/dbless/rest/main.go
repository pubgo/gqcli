package rest

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/pubgo/mycli/internal/dbless/rest/v1"
	"github.com/pubgo/mycli/version"
	"net/http"
)

func App() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(logger.SetLogger())

	r.StaticFile("/favicon.ico", "favicon.ico")

	app := r.Group("/db2rest")
	app.GET("/", func(ctx *gin.Context) {
		var rs []gin.H
		for _, rt := range r.Routes() {
			rs = append(rs, gin.H{
				"Method":  rt.Method,
				"Path":    rt.Path,
				"Handler": rt.Handler,
			})
		}
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"routes":  rs,
			"version": version.Version,
			"buildV":  version.BuildV,
			"commitV": version.CommitV,
		})
	})

	app.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
		return
	})

	// 服务端API
	api := app.Group("/api")
	v1.InitRouterV1(api.Group("/v1"))

	return r
}
