package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// custom error used if unable to parse JSON string in a request body
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

// custom Runtime type
type Runtime int32

// adds a custom MarshalJSON() method on the Runtime type
// to return "<runtime> mins"
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	// jsonValue needs double quotes to be a valid JSON string
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

// adds a custom UnmarshalJSON() method on the Runtime type
// to parse out "mins" from runtime field in request body
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// expected input is "<runtime> mins"
	// so remove quotes
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// split input string for sanity checking values
	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// parse string into an int32 value
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// convert i (int32) into a Runtime type and assign it
	*r = Runtime(i)

	return nil
}
