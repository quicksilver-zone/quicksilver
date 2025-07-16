package handlers

import (
	"fmt"
	"net/http"

	"github.com/quicksilver-zone/quicksilver/xcclookup/pkgs/types"
)

func GetVersionHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		version, err := types.GetVersion()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(version)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
	}
}
