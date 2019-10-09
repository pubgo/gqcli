package restorm

import (
	"database/sql"
	"fmt"
	"github.com/pubgo/mycli/pkg/logs"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pubgo/g/errors"
)

var once sync.Once
var _db *restOrm

// Default is default RestOrm instance
func Default() IRestOrm {
	once.Do(func() {
		_db = &restOrm{
			cfg: make(map[string]*Config),
		}
	})
	return _db
}

type restOrm struct {
	cfg map[string]*Config
}

func (t *restOrm) DbStats() map[string]sql.DBStats {
	stats := make(map[string]sql.DBStats)
	for k, v := range t.cfg {
		stats[k] = v.db.Stats()
	}
	return stats
}

func (t *restOrm) dbConnect(dbName string, cfg *Config) (err error) {
	defer errors.RespErr(&err)

	if !cfg.Enable {
		if _l := logger.Warn(); _l.Enabled() {
			_l.Msgf("db %s is not enable", dbName)
		}
		return
	}

	db, err := sqlx.Connect(cfg.Driver, cfg.Dsn)
	errors.PanicM(err, "%s connect error", cfg.Dsn)
	errors.PanicM(db.Ping(), "%s ping error", dbName)

	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}

	if cfg.MaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)
	}
	cfg.db = db
	return
}

func (t *restOrm) DbAdd(dbName string, cfg *Config) {
	_retryAt(time.Second, func(dur time.Duration) {
		errors.PanicM(t.dbConnect(dbName, cfg), "db(%s) connect error", dbName)
		t.cfg[dbName] = cfg
	})
}

func (t *restOrm) DbDelete(dbName string) {
	if db, ok := t.cfg[dbName]; ok {
		_retryAt(time.Second, func(dur time.Duration) {
			errors.PanicM(db.db.Close, "db(%s) closed error", dbName)
			delete(t.cfg, dbName)
		})
	}
}

// DbUpdate 更新字段和类型
func (t *restOrm) DbUpdate(dbName string, cfg *Config) {
	t.DbDelete(dbName)
	t.DbAdd(dbName, cfg)
}

// 检查数据库名称和数据表名
func (t *restOrm) checkDbAndTb(dbName, tbName string) error {
	if t.cfg[dbName] == nil {
		return fmt.Errorf("db name or table name does not exist,(%s, %s)", dbName, tbName)
	}
	return nil
}

// 批量创建记录
func (t *restOrm) ResCreateMany(dbName, tbName string, dts ...map[string]interface{}) (err error) {
	defer errors.RespErr(&err)

	// 检查db和tb
	errors.Panic(t.checkDbAndTb(dbName, tbName))

	_cfg := t.cfg[dbName]
	_tx, err := _cfg.db.Beginx()
	errors.PanicM(err, "db(%s) create tx error", dbName)

	_sql := &sqlBuilder{table: tbName}
	for _, dt := range dts {
		_sqlData := _sql.insertString(dt)
		if _cfg.ShowSQL {
			logs.P("_sqlData", _sql.args)
		}

		_, err := _tx.Exec(_sqlData, _sql.args...)
		if err != nil {
			errors.PanicM(_tx.Rollback(), "tx rollback error: %s", err)
		}

		errors.PanicMM(err, func(err *errors.Err) {
			err.Msg("db create error")
			err.M("input", dt)
			err.M("sql", _sqlData)
			err.M("db", dbName)
			err.M("tb", tbName)
			err.SetTag(ErrTag.DbCreateError)
		})
	}

	errors.PanicM(_tx.Commit(), "tx commit error")
	return
}

// 批量删除记录
func (t *restOrm) ResDeleteMany(dbName, tbName string, where ...interface{}) (err error) {
	defer errors.RespErr(&err)

	// 检查db和tb
	errors.Panic(t.checkDbAndTb(dbName, tbName))

	_cfg := t.cfg[dbName]
	_sql := &sqlBuilder{table: tbName}
	_sql.Where(where...)

	_sqlData := _sql.deleteString()
	if _cfg.ShowSQL {
		logs.P("_sqlData", _sql.args)
	}

	_, err = _cfg.db.Exec(_sqlData, _sql.args...)
	errors.PanicMM(err, func(err *errors.Err) {
		err.Msg("db delete error")
		err.M("sql", _sqlData)
		err.M("db", dbName)
		err.M("tb", tbName)
		err.SetTag(ErrTag.DbDeleteError)
	})

	return
}

// 批量修改记录
func (t *restOrm) ResUpdateMany(dbName, tbName string, data map[string]interface{}, where ...interface{}) (err error) {
	defer errors.RespErr(&err)

	errors.PanicT(_isNone(data), "data is nil")

	// 检查db和tb
	errors.Panic(t.checkDbAndTb(dbName, tbName))

	_cfg := t.cfg[dbName]
	_sql := &sqlBuilder{table: tbName}
	_sql.Where(where...)

	_sqlData := _sql.updateString(data)
	if _cfg.ShowSQL {
		logs.P("_sqlData", _sql.args)
	}

	_, err = _cfg.db.Exec(_sqlData, _sql.args...)
	errors.PanicMM(err, func(err *errors.Err) {
		err.Msg("db update error")
		err.M("sql", _sqlData)
		err.M("db", dbName)
		err.M("tb", tbName)
		err.M("data", data)
		err.SetTag(ErrTag.DbUpdateError)
	})

	return
}

// 过滤条件统计
func (t *restOrm) ResCount(dbName, tbName string, where ...interface{}) (c int64, err error) {
	defer errors.RespErr(&err)

	// 检查db和tb
	errors.Panic(t.checkDbAndTb(dbName, tbName))

	_cfg := t.cfg[dbName]
	_sql := &sqlBuilder{table: tbName}
	_sql.Where(where...)

	_sqlData := _sql.countString()
	if _cfg.ShowSQL {
		logs.P(_sqlData, _sql.args)
	}

	errors.PanicMM(_cfg.db.Select(&c, _sqlData, _sql.args...), func(err *errors.Err) {
		err.Msg("db count error")
		err.M("sql", _sqlData)
		err.M("dbName", dbName)
		err.M("tbName", tbName)
		err.SetTag(ErrTag.DbCountError)
	})
	return
}

func (t *restOrm) rows2Map(dbName, tbName string, rows *sqlx.Rows) (dts []map[string]interface{}, err error) {
	defer errors.RespErr(&err)

	for rows.Next() {
		dest := make(map[string]interface{})

		cons, err := rows.ColumnTypes()
		errors.PanicM(err, "get rows column types error")

		values := make([]interface{}, len(cons))
		for i, con := range cons {
			switch strings.ToLower(con.DatabaseTypeName()) {
			case "tinyint", "int", "smallint", "mediumint", "bigint", "integer":
				values[i] = new(sql.NullInt64)
			case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext":
				values[i] = new(sql.NullString)
			case "date", "datetime", "time", "timestamp":
				values[i] = new(time.Time)
			case "decimal", "double", "float":
				values[i] = new(sql.NullFloat64)
			case "binary", "blob", "longblob", "mediumblob", "varbinary":
				values[i] = new(sql.NullString)
			case "bool":
				values[i] = new(sql.NullBool)
			case "":
				switch con.ScanType().Name() {
				case "int64", "int32", "int":
					values[i] = new(sql.NullInt64)
				case "float64", "float32":
					values[i] = new(sql.NullFloat64)
				case "string":
					values[i] = new(sql.NullString)
				}
			default:
				errors.PanicT(true, "unknown type")
			}
		}

		errors.PanicM(rows.Scan(values...), "rows scan error")
		for i, column := range cons {
			colN := column.Name()
			switch _v := values[i].(type) {
			case *sql.NullFloat64:
				dest[colN] = _v.Float64
			case *sql.NullString:
				dest[colN] = _v.String
			case *sql.NullInt64:
				dest[colN] = _v.Int64
			case *time.Time:
				dest[colN] = _v.Unix()
			case *sql.NullBool:
				dest[colN] = _v.Bool
			default:
				errors.PanicT(true, "unknown type")
			}
		}
		dts = append(dts, dest)
	}
	return
}

// 过滤查询
func (t *restOrm) ResGetMany(dbName, tbName string, fields string, groupBy string, order string, limit, offset string, where ...interface{}) (dts []map[string]interface{}, err error) {
	defer errors.RespErr(&err)

	// 检查db和tb
	errors.Panic(t.checkDbAndTb(dbName, tbName))

	_cfg := t.cfg[dbName]
	_sql := &sqlBuilder{table: tbName, fields: fields}
	_sql.groupBy = groupBy
	_sql.orderBy = order
	_sql.limit = limit
	_sql.offset = offset
	_sql.Where(where...)

	_sqlData := _sql.queryString()
	if t.cfg[dbName].ShowSQL {
		logs.P(_sqlData, _sql.args)
	}

	rows, err := _cfg.db.Queryx(_sqlData, _sql.args...)
	errors.PanicMM(err, func(err *errors.Err) {
		err.Msg("db get error")
		err.M("sql", _sqlData)
		err.M("db", dbName)
		err.M("tb", tbName)
		err.SetTag(ErrTag.DbGetError)
	})

	return t.rows2Map(dbName, tbName, rows)
}
