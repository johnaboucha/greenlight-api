package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// envelope to wrap JSON output in a parent, "movie", key
type envelope map[string]interface{}

func (app *application) readIDParam(r *http.Request) (int64, error) {

	// retrieves a slice of URL parameters
	params := httprouter.ParamsFromContext(r.Context())

	// checks and returns valid ID parameter
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	//js, err := json.Marshal(data)
	js, err := json.MarshalIndent(data, "", "\t") // add indents at a slight performance cost
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// reads incoming JSON request and handles bad input better than default Go error handling
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	// limit the request body to 1MB
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// restrict unknown fields in request body
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// decode JSON into target
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// syntax error
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

			// End of file error
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badlly-formed JSON")

		// Incorrect field type
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// empty body
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// unknown field names, cannot be mapped
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// request body was too large
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// invalid argument
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// otherwise, just pass along the error
		default:
			return err
		}
	}

	// run Decode again to check if there is
	// additional data after the valid input
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
