package restorm

import (
	"github.com/jmoiron/sqlx"
)

//Config is database connection configuration
type Config struct {
	Enable       bool   `toml:"enable" json:"enable"`
	Driver       string `toml:"driver" json:"driver"`
	Dsn          string `toml:"dsn" json:"dsn"`
	MaxOpenConns int    `toml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns int    `toml:"max_idle_conns" json:"max_idle_conns"`
	MaxLifetime  int    `toml:"max_left_time" json:"max_left_time"`
	ShowSQL      bool   `toml:"show_sql" json:"show_sql"`

	db *sqlx.DB
}
