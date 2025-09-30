package decode

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var ErrEmptyBody = errors.New("request body is empty")

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return ErrEmptyBody
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dst)
	if err != nil {
		if err == io.EOF {
			return ErrEmptyBody
		}
		return err
	}
	return nil
}
