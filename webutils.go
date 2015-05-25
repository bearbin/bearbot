package main

import (
	"net/http"

	"github.com/PuerkitoBio/ghost/handlers"
)

func gzipHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(handlers.GZIPHandler(h, nil))
}
