package restorm

import "github.com/pubgo/g/errors"

var ErrTag = struct {
	DbCreateError    string `toml:"db_create_error"`
	DbDeleteError    string `toml:"db_delete_error"`
	DbUpdateError    string `toml:"db_update_error"`
	DbCountError     string `toml:"db_count_error"`
	DbGetError       string `toml:"db_get_error"`
	DbDeleteOneError string `toml:"db_delete_one_error"`
	DbUpdateOneError string `toml:"db_update_one_error"`
	DbGetOneError    string `toml:"db_get_one_error"`
}{
	DbCreateError:    "db_create_error",
	DbDeleteError:    "db_delete_error",
	DbUpdateError:    "db_update_error",
	DbCountError:     "db_count_error",
	DbGetError:       "db_get_error",
	DbDeleteOneError: "db_delete_one_error",
	DbUpdateOneError: "db_update_one_error",
	DbGetOneError:    "db_get_one_error",
}

func init() {
	errors.ErrTagRegistry(ErrTag)
}
