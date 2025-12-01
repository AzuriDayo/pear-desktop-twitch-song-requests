package data

import (
	"embed"
)

//go:embed iofs/migrations/*
var MigrationsFS embed.FS
