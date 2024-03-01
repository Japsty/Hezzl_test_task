package migrations_clickhouse

import "embed"

//go:embed *.sql
var Embed embed.FS
