package webui

import (
	"embed"
	"fmt"
	"net/http"
	"time"
)

//go:embed out/*
var webFS embed.FS

func UI() {
	http.Handle("/", http.FileServer(http.FS(webFS)))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      nil, // Use default mux
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting embedded web app server: ", err)
	}
}
