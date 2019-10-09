package restorm

import "database/sql"

// IRestOrm is Restorm interface
type IRestOrm interface {
	// DbAdd 添加数据库配置，并获取表字段和类型结构
	// dbName 数据库别名，最好是数据库的名字
	// cfg 配置文件
	DbAdd(dbName string, cfg *Config)

	// DbDelete 删除数据库配置，暂停数据库连接以及删除表字段和类型
	DbDelete(dbName string)

	// DbUpdate 更新数据库配置，更新表字段以及类型或者更新数据库连接
	DbUpdate(dbName string, cfg *Config)

	// ResCreateMany 批量添加记录
	// dbName 数据库名称或者添加配置文件时的别名
	// tbName 数据库表名称
	ResCreateMany(dbName, tbName string, dts ...map[string]interface{}) (err error)

	// ResDeleteMany 批量删除记录
	ResDeleteMany(dbName, tbName string, where ...interface{}) (err error)

	// ResUpdateMany 根据过滤修改记录
	ResUpdateMany(dbName, tbName string, data map[string]interface{}, where ...interface{}) (err error)

	// ResCount 根据过滤统计查询信息
	// where 数据库过滤条件
	ResCount(dbName, tbName string, where ...interface{}) (c int64, err error)

	// ResGetMany 根据过滤统计查询信息
	// fields 查询的字段信息
	// groupBy 分组字段
	// order 排序字段
	// limit 查询后的限制数量
	// offset 偏移量
	// where 数据库过滤条件
	ResGetMany(dbName, tbName string, fields string, groupBy string, order string, limit, offset string, where ...interface{}) (dts []map[string]interface{}, err error)

	// DbStats 获取数据库连接状态
	DbStats() map[string]sql.DBStats
}
