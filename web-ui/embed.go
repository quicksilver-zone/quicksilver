package webui

import (
	"embed"
	"net/http"
)

//go:embed out/*
var webFS embed.FS

func ui() {
	http.Handle("/", http.FileServer(http.FS(webFS)))
	http.ListenAndServe(":8080", nil)
}
