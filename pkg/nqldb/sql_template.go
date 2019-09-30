package nqldb

import (
	"github.com/pubgo/g/errors"
	"text/template"
)

var _base_tmpl string

func init() {
	_, err := template.New("main").Parse(`
CREATE TABLE "{{TableName}}" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "create_at" integer(11) NOT NULL,
  "update_at" integer(11) NOT NULL
);

CREATE INDEX "create_at" ON "{{TableName}}" ( "create_at" DESC );
CREATE INDEX "update_at" ON "{{TableName}}" ( "update_at" DESC );
`)
	errors.Panic(err)
}
