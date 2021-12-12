package data

import (
	"fmt"
	"strconv"
)

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
