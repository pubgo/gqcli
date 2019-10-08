package restorm

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pubgo/g/errors"
	"github.com/pubgo/schema"
)

var once sync.Once
var _db *RestOrm

func Default() *RestOrm {
	once.Do(func() {
		_db = &RestOrm{
			cfg: make(map[string]*Config),
		}
	})
	return _db
}

type RestOrm struct {
	cfg map[string]*Config
}

// colsTransfer get table field name and type
func (t *RestOrm) colsTransfer(name string, cfg *Config) (err error) {
	defer errors.RespErr(&err)

	tbs, err := schema.Tables(cfg.db.DB)
	errors.PanicM(err, "get table schema error")

	if cfg.tbs == nil {
		cfg.tbs = make(map[string]*tb)
	}

	for name, tps := range tbs {
		cfg.tbs[name] = &tb{ColT: make(map[string]*converter)}
		for _, f := range tps {
			fieldTp := strings.ToLower(f.DatabaseTypeName())
			cfg.tbs[name].ColT[f.Name()] = &converter{TpName: fieldTp, converter: Converter(fieldTp)}
		}
	}
	return
}

func (t *RestOrm) DbStats() map[string]sql.DBStats {
	stats := make(map[string]sql.DBStats)
	for k, v := range t.cfg {
		stats[k] = v.db.Stats()
	}
	return stats
}

func (t *RestOrm) ColTs() map[string]map[string]map[string]string {
	var dt = make(map[string]map[string]map[string]string)
	for k, v := range t.cfg {
		dt[k] = make(map[string]map[string]string)
		for k1, v1 := range v.tbs {
			dt[k][k1] = make(map[string]string)
			for k2, v := range v1.ColT {
				dt[k][k1][k2] = v.TpName
			}
		}
	}
	return dt
}

func (t *RestOrm) dbConnect(key string, conf *Config) (err error) {
	defer errors.RespErr(&err)

	if !conf.Enable {
		return
	}

	db, err := sqlx.Connect(conf.Driver, conf.Dsn)
	errors.PanicM(err, "%s connect error", conf.Dsn)
	errors.PanicM(db.Ping(), "ping error")

	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	if conf.MaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)
	}
	conf.db = db

	errors.Panic(t.colsTransfer(key, conf))
	return
}

func (t *RestOrm) DbConfigAdd(name string, cfg *Config) {
	_retryAt(time.Second, func(dur time.Duration) {
		errors.PanicM(t.dbConnect(name, cfg), "db connect error")
		t.cfg[name] = cfg
	})
}

func (t *RestOrm) DbConfigDelete(name string) {
	if db, ok := t.cfg[name]; ok {
		_retryAt(time.Second, func(dur time.Duration) {
			errors.PanicM(db.db.Close, "db(%s) closed error", name)
			delete(t.cfg, name)
		})
	}
}

// DbUpdate 更新字段和类型
func (t *RestOrm) DbUpdate(name string) (err error) {
	defer errors.RespErr(&err)

	cfg, ok := t.cfg[name]
	errors.PanicT(!ok, "db(%s) does not exist", name)
	errors.Panic(t.colsTransfer(name, cfg))
	t.cfg[name] = cfg
	return
}

func (t *RestOrm) DbConfigUpdate(name string, cfg *Config) {
	t.DbConfigDelete(name)
	t.DbConfigAdd(name, cfg)
}

// 检查数据库名称和数据表名
func (t *RestOrm) checkDbAndTb(dbName, tbName string) error {
	db := t.cfg[dbName]
	if db == nil || db.tbs[tbName] == nil {
		return fmt.Errorf("db name or table name does not exist,(%s, %s)", dbName, tbName)
	}
	return nil
}

// 批量创建记录
func (t *RestOrm) ResCreateMany(dbName, tbName string, dts ...map[string]interface{}) (err error) {
	defer errors.RespErr(&err)

	// 检查db和tb
	errors.Panic(t.checkDbAndTb(dbName, tbName))

	_db := t.cfg[dbName].db
	_tx, err := _db.Beginx()
	errors.PanicM(err, "db(%s) create tx error", dbName)

	_sql := &sqlBuilder{table: tbName}
	for _, dt := range dts {
		_, err := _tx.Exec(_sql.insertString(dt), _sql.args...)
		if err != nil {
			errors.PanicM(_tx.Rollback, "tx rollback error: %s", err)
		}

		errors.PanicMM(err, func(err *errors.Err) {
			err.Msg("db create error")
			err.M("input", dt)
			err.M("db", dbName)
			err.M("tb", tbName)
			err.SetTag(ErrTag.DbCreateError)
		})
	}

	errors.PanicM(_tx.Commit, "tx commit error")
	return
}

// 批量删除记录
func (t *RestOrm) ResDeleteMany(dbName, tbName string, filter ...interface{}) (err error) {
	defer errors.RespErr(&err)

	// 检查db和tb
	errors.Panic(t.checkDbAndTb(dbName, tbName))

	_db := t.cfg[dbName].db
	_sql := &sqlBuilder{table: tbName}
	_sql.Where(filter...)

	_, err = _db.Exec(_sql.deleteString(), _sql.args...)
	errors.PanicMM(err, func(err *errors.Err) {
		err.Msg("db delete error")
		err.M("input", filter)
		err.M("db", dbName)
		err.M("tb", tbName)
		err.SetTag(ErrTag.DbDeleteError)
	})

	return
}

// 批量修改记录
func (t *RestOrm) ResUpdateMany(dbName, tbName string, data map[string]interface{}, filter ...interface{}) (err error) {
	defer errors.RespErr(&err)

	errors.PanicT(_isNone(data), "data is nil")

	// 检查db和tb
	errors.Panic(t.checkDbAndTb(dbName, tbName))

	_db := t.cfg[dbName].db
	_sql := &sqlBuilder{table: tbName}
	_sql.Where(filter...)

	_, err = _db.Exec(_sql.updateString(data), _sql.args...)
	errors.PanicMM(err, func(err *errors.Err) {
		err.Msg("db update error")
		err.M("input", filter)
		err.M("db", dbName)
		err.M("tb", tbName)
		err.M("data", data)
		err.SetTag(ErrTag.DbUpdateError)
	})

	return
}

// 过滤条件统计
func (t *RestOrm) ResCount(dbName, tbName string, filter ...interface{}) (c int64, err error) {
	defer errors.RespErr(&err)

	// 检查db和tb
	errors.Panic(t.checkDbAndTb(dbName, tbName))

	_db := t.cfg[dbName].db
	_sql := &sqlBuilder{table: tbName}
	_sql.Where(filter...)

	errors.PanicMM(_db.Select(&c, _sql.countString(), _sql.args...), func(err *errors.Err) {
		err.Msg("db count error")
		err.M("input", filter)
		err.M("dbName", dbName)
		err.M("tbName", tbName)
		err.SetTag(ErrTag.DbCountError)
	})
	return
}

func (t *RestOrm) rows2Map(dbName, tbName string, rows *sqlx.Rows) (res []map[string]interface{}, err error) {
	defer errors.RespErr(&err)

	// 检查db和tb
	errors.Panic(t.checkDbAndTb(dbName, tbName))

	var dts []map[string]interface{}
	for rows.Next() {
		dest := make(map[string]interface{})

		cons, err := rows.ColumnTypes()
		errors.Panic(err)

		values := make([]interface{}, len(cons))
		for i := range values {
			values[i] = new(interface{})
		}

		_tb := t.cfg[dbName].tbs[tbName]
		errors.Panic(rows.Scan(values...))
		for i, column := range cons {
			k := column.Name()
			if _fn, ok := _tb.ColT[k]; ok {
				dest[k] = _fn.converter(values[i])
			} else {
				dest[k] = Converter(strings.ToLower(column.DatabaseTypeName()))(values[i])
			}
		}
		dts = append(dts, dest)
	}

	res = dts
	return
}

// 过滤查询
func (t *RestOrm) ResGetMany(dbName, tbName string, fields string, groupBy string, order string, limit, offset string, filter ...interface{}) (dts []map[string]interface{}, err error) {
	defer errors.RespErr(&err)

	// 检查db和tb
	errors.Panic(t.checkDbAndTb(dbName, tbName))

	_db := t.cfg[dbName].db
	_sql := &sqlBuilder{table: tbName, fields: fields}
	_sql.groupBy = groupBy
	_sql.orderBy = order
	_sql.limit = limit
	_sql.offset = offset
	_sql.Where(filter...)

	rows, err := _db.Queryx(_sql.queryString(), _sql.args...)
	errors.PanicMM(err, func(err *errors.Err) {
		err.Msg("db get error")
		err.M("input", filter)
		err.M("db", dbName)
		err.M("tb", tbName)
		err.M("sql", _sql)
		err.SetTag(ErrTag.DbGetError)
	})

	return t.rows2Map(dbName, tbName, rows)
}
