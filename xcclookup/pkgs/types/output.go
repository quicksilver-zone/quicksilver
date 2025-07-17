package types

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ingenuity-build/multierror"
)

func outputResponse(w http.ResponseWriter, response *Response, errors map[string]error, clearMessages bool) {
	if len(errors) > 0 {
		response.Errors = multierror.New(errors)
	}

	if clearMessages {
		response.Messages = nil
	}

	jsonOut, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}
	fmt.Fprint(w, string(jsonOut))
}

var OutputEpoch = func(w http.ResponseWriter, response *Response, errors map[string]error) {
	outputResponse(w, response, errors, false)
}

var OutputCurrent = func(w http.ResponseWriter, response *Response, errors map[string]error) {
	outputResponse(w, response, errors, true)
}
