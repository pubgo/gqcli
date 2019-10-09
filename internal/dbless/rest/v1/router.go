package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pubgo/g/errors"
	"github.com/pubgo/mycli/pkg/restorm"
	"net/http"
	"strings"
)

func InitRouterV1(r *gin.RouterGroup) {

	// 添加数据
	r.POST("/rds/:db/:tb", func(ctx *gin.Context) {
		defer errors.Resp(func(err *errors.Err) {
			ctx.IndentedJSON(http.StatusBadRequest, err.StackTrace())
		})

		_db := ctx.Param("db")
		_tb := ctx.Param("tb")

		var _dts []map[string]interface{}
		errors.PanicM(ctx.ShouldBindJSON(&_dts), "参数解析失败")

		rtm := restorm.Default()
		errors.PanicM(rtm.ResCreateMany(_db, _tb, _dts...), "数据存储失败")
		ctx.String(http.StatusOK, "")
	})

	// 批量删除
	r.DELETE("/rds/:db/:tb", func(ctx *gin.Context) {
		defer errors.Resp(func(err *errors.Err) {
			ctx.IndentedJSON(http.StatusBadRequest, err.StackTrace())
		})

		_db := ctx.Param("db")
		_tb := ctx.Param("tb")

		var _where []string
		var _whereData []interface{}
		_whereData = append(_whereData, "")
		for k, v := range ctx.Request.URL.Query() {
			if len(v) == 0 {
				continue
			}

			_where = append(_where, fmt.Sprintf("%s=? ", k))
			_whereData = append(_whereData, v[0])
		}
		_whereData[0] = strings.Join(_where, " and ")

		rtm := restorm.Default()
		errors.PanicM(rtm.ResDeleteMany(_db, _tb, _whereData...), "数据存储失败")
		ctx.String(http.StatusOK, "")
	})

	// 批量修改
	r.PUT("/rds/:db/:tb", func(ctx *gin.Context) {
		defer errors.Resp(func(err *errors.Err) {
			ctx.IndentedJSON(http.StatusBadRequest, err.StackTrace())
		})

		_db := ctx.Param("db")
		_tb := ctx.Param("tb")
		var _dt map[string]interface{}
		errors.PanicM(ctx.ShouldBindJSON(&_dt), "参数解析失败")

		var _where []string
		var _whereData []interface{}
		_whereData = append(_whereData, "")
		for k, v := range ctx.Request.URL.Query() {
			if len(v) == 0 {
				continue
			}

			_where = append(_where, fmt.Sprintf("%s=? ", k))
			_whereData = append(_whereData, v[0])
		}
		_whereData[0] = strings.Join(_where, " and ")

		rtm := restorm.Default()
		errors.PanicM(rtm.ResUpdateMany(_db, _tb, _dt, _whereData...), "数据存储失败")
		ctx.String(http.StatusOK, "")
	})

	// 批量获取
	r.GET("/rds/:db/:tb", func(ctx *gin.Context) {
		defer errors.Resp(func(err *errors.Err) {
			ctx.IndentedJSON(http.StatusBadRequest, err.StackTrace())
		})

		_db := ctx.Param("db")
		_tb := ctx.Param("tb")
		_groupBy := ctx.Query("$groupBy")
		_orderBy := ctx.Query("$order")
		_limit := ctx.Query("$limit")
		_offset := ctx.Query("$offset")
		_fields := ctx.Query("$fields")
		//_filter := ctx.Query("$filter")

		var _where []string
		var _whereData []interface{}
		_whereData = append(_whereData, "")
		for k, v := range ctx.Request.URL.Query() {
			if strings.Contains(k, "$") {
				continue
			}
			if len(v) == 0 {
				continue
			}

			_where = append(_where, fmt.Sprintf("%s=? ", k))
			_whereData = append(_whereData, v[0])
		}
		_whereData[0] = strings.Join(_where, " and ")

		rtm := restorm.Default()
		dts, err := rtm.ResGetMany(_db, _tb, _fields, _groupBy, _orderBy, _limit, _offset, _whereData...)
		errors.PanicM(err, "数据查询失败")
		ctx.JSON(http.StatusOK, gin.H{
			"data": dts,
		})
	})

	// 分页查询加上简单统计
	r.GET("/rds/:db/:tb/count", func(ctx *gin.Context) {
		defer errors.Resp(func(err *errors.Err) {
			ctx.IndentedJSON(http.StatusBadRequest, err.StackTrace())
		})

		_db := ctx.Param("db")
		_tb := ctx.Param("tb")

		var _where []string
		var _whereData []interface{}
		_whereData = append(_whereData, "")
		for k, v := range ctx.Request.URL.Query() {
			if len(v) == 0 {
				continue
			}

			_where = append(_where, fmt.Sprintf("%s=? ", k))
			_whereData = append(_whereData, v[0])
		}
		_whereData[0] = strings.Join(_where, " and ")

		rtm := restorm.Default()
		c, err := rtm.ResCount(_db, _tb, _whereData...)
		errors.PanicM(err, "数据统计失败")
		ctx.JSON(http.StatusOK, gin.H{
			"data": c,
		})
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
