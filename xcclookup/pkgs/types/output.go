package types

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/multierr"
)

func outputResponse(w http.ResponseWriter, response *Response, errors map[string]error, clearMessages bool) {
	// Set appropriate headers for JSON API responses
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	if len(errors) > 0 {
		var err error
		for key, e := range errors {
			err = multierr.Append(err, fmt.Errorf("%s: %w", key, e))
		}
		response.Errors = &ErrorString{error: err}
	}

	if clearMessages {
		response.Messages = nil
	}

	jsonOut, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonOut))
}

var OutputEpoch = func(w http.ResponseWriter, response *Response, errors map[string]error) {
	outputResponse(w, response, errors, false)
}

var OutputCurrent = func(w http.ResponseWriter, response *Response, errors map[string]error) {
	outputResponse(w, response, errors, true)
}
