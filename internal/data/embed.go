package data

import (
	"embed"
)

//go:embed iofs/migrations/*.sql
var FS embed.FS
