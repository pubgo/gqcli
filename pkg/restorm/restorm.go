package restorm

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pubgo/g/errors"
	"github.com/pubgo/schema"
	"strings"
	"sync"
	"time"
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

// colsTransfer get field type and name
func (t *RestOrm) colsTransfer(name string, cfg *Config) (err error) {
	defer errors.RespErr(&err)

	tbs, err := schema.Tables(cfg.db.DB)
	errors.PanicM(err, "get tables error")

	if cfg.colT == nil {
		cfg.colT = make(map[string]map[string]*converter)
	}

	for name, tps := range tbs {
		cfg.colT[name] = make(map[string]*converter)
		for _, f := range tps {
			field := strings.ToLower(f.DatabaseTypeName())
			cfg.colT[name][f.Name()] = &converter{Name: field, converter: Converter(field)}
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
		for k1, v1 := range v.colT {
			dt[k][k1] = make(map[string]string)
			for k2, v := range v1 {
				dt[k][k1][k2] = v.Name
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

	// update table fields type
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
			errors.PanicM(db.db.Close, "db(%s) close error", name)
			delete(t.cfg, name)
		})
	}
}

func (t *RestOrm) DbUpdate(name string) (err error) {
	defer errors.RespErr(&err)

	cfg, ok := t.cfg[name]
	errors.PanicT(!ok, "db(%s) does not exist", name)

	_retryAt(time.Second, func(dur time.Duration) {
		errors.PanicM(t.dbConnect(name, cfg), "db(%s) update error", name)
		t.cfg[name] = cfg
	})

	return
}

func (t *RestOrm) DbConfigUpdate(name string, cfg *Config) (err error) {
	defer errors.RespErr(&err)

	_retryAt(time.Second, func(dur time.Duration) {
		errors.PanicM(t.dbConnect(name, cfg), "db(%s) config update error", name)
		t.cfg[name] = cfg
	})

	return
}

// 创建记录
func (t *RestOrm) ResCreateMany(dbName, tbName string, dts ...map[string]interface{}) (err error) {
	defer errors.RespErr(&err)

	_db := t.cfg[dbName].db
	_tx, err := _db.Beginx()
	errors.Panic(err)

	_sql := &sqlBuilder{table: tbName}
	for _, dt := range dts {
		_, err := _tx.Exec(_sql.insertString(dt), _sql.args...)
		if err != nil {
			errors.PanicM(_tx.Rollback, "tx rollback error: %s", err)
		}

		errors.PanicMM(err, func(err *errors.Err) {
			err.Msg("db create error")
			err.M("input", dt)
			err.M("dbName", dbName)
			err.M("tbName", tbName)
			err.SetTag(ErrTag.DbCreateError)
		})
	}
	errors.PanicM(_tx.Commit(), "tx commit error")

	return
}

// 删除记录
func (t *RestOrm) ResDeleteMany(dbName, tbName string, filter ...interface{}) (err error) {
	defer errors.RespErr(&err)

	_db := t.cfg[dbName].db
	_sql := &sqlBuilder{table: tbName}
	_sql.Where(filter...)

	_, err = _db.Exec(_sql.deleteString(), _sql.args...)
	errors.PanicMM(err, func(err *errors.Err) {
		err.Msg("db delete error")
		err.M("input", filter)
		err.M("dbName", dbName)
		err.M("tbName", tbName)
		err.SetTag(ErrTag.DbDeleteError)
	})

	return
}

// 修改记录
func (t *RestOrm) ResUpdateMany(dbName, tbName string, data map[string]interface{}, filter ...interface{}) (err error) {
	defer errors.RespErr(&err)

	errors.PanicT(_isNone(data), "update data is nil")

	_db := t.cfg[dbName].db
	_sql := &sqlBuilder{table: tbName}
	_sql.Where(filter...)

	_, err = _db.Exec(_sql.updateString(data), _sql.args...)
	errors.PanicMM(err, func(err *errors.Err) {
		err.Msg("db update error")
		err.M("input", filter)
		err.M("dbName", dbName)
		err.M("tbName", tbName)
		err.M("data", data)
		err.SetTag(ErrTag.DbUpdateError)
	})

	return
}

func (t *RestOrm) ResCount(dbName, tbName string, filter ...interface{}) (c int64, err error) {
	defer errors.RespErr(&err)

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

	var dts []map[string]interface{}
	for rows.Next() {
		dest := make(map[string]interface{})

		cons, err := rows.ColumnTypes()
		errors.Panic(err)

		values := make([]interface{}, len(cons))
		for i := range values {
			values[i] = new(interface{})
		}

		_tb := t.cfg[dbName].colT[tbName]
		errors.Panic(rows.Scan(values...))
		for i, column := range cons {
			k := column.Name()
			if _fn, ok := _tb[k]; ok {
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

// 查询
func (t *RestOrm) ResGetMany(dbName, tbName string, fields string, groupBy string, order string, limit, offset string, filter ...interface{}) (dts []map[string]interface{}, err error) {
	defer errors.RespErr(&err)

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
		err.M("dbName", dbName)
		err.M("tbName", tbName)
		err.M("sql", _sql)
		err.SetTag(ErrTag.DbGetError)
	})

	return t.rows2Map(dbName, tbName, rows)
}
