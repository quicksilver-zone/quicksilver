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
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}
		w.Write(version)
	}
}
