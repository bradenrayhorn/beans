//go:build !test
// +build !test

package web

import (
	"embed"
)

//go:embed dist/*
var FrontendFS embed.FS
