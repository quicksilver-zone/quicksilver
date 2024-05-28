package types

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ingenuity-build/multierror"
)

var OutputEpoch = func(w http.ResponseWriter, response *Response, errors map[string]error) {
	fmt.Println("check for errors...")
	if len(errors) > 0 {
		fmt.Printf("found %d error(s)\n", len(errors))
		response.Errors = multierror.New(errors)
		fmt.Println(response.Errors)
	}

	jsonOut, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}
	fmt.Fprint(w, string(jsonOut))
}

var OutputCurrent = func(w http.ResponseWriter, response *Response, errors map[string]error) {
	fmt.Println("check for errors...")
	if len(errors) > 0 {
		fmt.Printf("found %d error(s)\n", len(errors))
		response.Errors = multierror.New(errors)
		fmt.Println(response.Errors)
	}

	response.Messages = nil
	jsonOut, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprint(w, string(jsonOut))
}
