package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/pubgo/errors"
)

func InitRouterV1(r *gin.RouterGroup) {
	//qd:=ctx.Request.URL.Query()
	//ctx.Param("id")
	//ctx.GetHeader()
	//ctx.Errors.JSON()
	//ctx.Error()
	//ctx.ClientIP()
	//ctx.AbortWithError()
	//ctx.GetRawData()
	//var dt=make(map[string]interface{})
	//ctx.ShouldBindJSON()

	// 添加数据
	r.POST("/rds/:db/:tbs", func(ctx *gin.Context) {

	})

	// 批量删除
	r.DELETE("/rds/:db/:tbs", func(ctx *gin.Context) {
		//	 从query获取数据
		// 从body获取数据

		var dt = make(map[string]interface{})
		errors.Wrap(ctx.ShouldBindJSON(dt), "")
		//_where := dt["where"]

		qd := ctx.Request.URL.Query()
		for k, v := range qd {
			dt[k] = errors.If(len(v) == 0, "", v[0])
		}
	})

	// 批量修改
	r.PUT("/rds/:db/:tbs", nil)
	// 批量获取
	r.GET("/rds/:db/:tbs", nil)
	// 分页查询加上简单统计
	r.GET("/rds/:db/:tbs/page", nil)
	r.GET("/rds/:db/:tbs/count", nil)

	// 修改单个
	r.PUT("/rd/:db/:tb/:id", nil)
	// 获取单个
	r.GET("/rd/:db/:tb/:id", nil)
	// 删除单个
	r.DELETE("/rd/:db/:tb/:id", func(ctx *gin.Context) {
		//id := ctx.Param("id")
		//assert.T(id == "", "")
	})

	// 耗时的大数据量查询
	r.POST("/queries/:db/:tb/", nil)
	r.GET("/queries/:id", nil)

	// 自定义脚本执行
	r.POST("/script/:db/:name", nil)
	r.GET("/script/:name/:id", nil) // 结果缓存时间预定

	// sql语句执行
	r.POST("/sql/:db/:name", nil)
	r.GET("/sql/:name/:id", nil)
}
