package main

import (
	"net/http"

	"github.com/gobuffalo/packr/v2"
)

var (
	frontendBox = packr.New("frontend", "dist/frontend/")
)

func frontendHandler() http.Handler {
	return http.FileServer(frontendBox)
}
