package restorm

import (
	"database/sql"
	"github.com/pubgo/g/errors"
	"math/big"
	"strconv"
	"time"
)

func _ToInt(p string) int {
	r, err := strconv.Atoi(p)
	errors.PanicM(err, "can not convert %s to int", p)
	return r
}

func _ToFloat(p string) float64 {
	f, err := strconv.ParseFloat(p, 0)
	errors.PanicM(err, "parse float type error,input(%s)", p)
	return f
}

func Converter(sqlType string) func(interface{}) interface{} {
	return func(dt interface{}) interface{} {
		_isNil := _isNone(dt)

		switch sqlType {
		case "tinyint", "int", "smallint", "mediumint", "bigint":
			if _isNil {
				return 0
			}

			switch _v := dt.(type) {
			case int64, int32, int16, int8, int, sql.NullInt64:
				return _v
			case string:
				return _ToInt(_v)
			case []byte:
				return big.NewInt(0).SetBytes(_v).Int64()
			}
		case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext":
			if _isNil {
				return ""
			}

			switch _v := dt.(type) {
			case string, sql.NullString:
				return _v
			case []byte:
				return string(_v)
			case int:
				return strconv.Itoa(_v)
			}
		case "date", "datetime", "time", "timestamp":
			if _isNil {
				return time.Now()
			}

			switch _v := dt.(type) {
			case string:
				t, err := time.Parse(time.RFC3339Nano, _v)
				errors.PanicM(err, "time parse error")
				return t
			case time.Time:
				return _v
			case []byte:
				t, err := time.Parse(time.RFC3339Nano, string(_v))
				errors.PanicM(err, "time parse error")
				return t
			}
		case "decimal", "double", "float":
			if _isNil {
				return 0.0
			}

			switch _v := dt.(type) {
			case float64, float32:
				return _v
			case string:
				return _ToFloat(_v)
			case []byte:
				return _ToFloat(string(_v))
			}

		case "binary", "blob", "longblob", "mediumblob", "varbinary":
			if _isNil {
				return ""
			}

			switch _v := dt.(type) {
			case string:
				return _v
			case []byte:
				return string(_v)
			}
		default:
			errors.PanicT(true, "unknown type")
		}
		return nil
	}
}
