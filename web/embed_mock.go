//go:build test
// +build test

package web

import (
	"embed"
)

//go:embed src/*
var FrontendFS embed.FS
