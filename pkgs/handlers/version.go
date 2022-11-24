package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Version struct {
	Version string
}

func GetVersionHandler(version string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		v := Version{
			Version: version,
		}

		jsonOut, err := json.Marshal(v)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
		fmt.Fprint(w, string(jsonOut))
	}
}
