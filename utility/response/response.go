package response

import (
	"encoding/json"
	"errors"
	"net/http"

	customErr "github.com/toggle-feature/utility/error"
)

type Meta struct {
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message,omitempty"`
	HttpStatus int         `json:"http_status"`
}

func Respond(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteSuccess(w http.ResponseWriter, data interface{}, message string) error {
	meta := Meta{
		Message:    message,
		Data:       data,
		HttpStatus: http.StatusOK,
	}

	Respond(w, meta, http.StatusOK)
	return nil
}

func WriteError(w http.ResponseWriter, err error) error {
	httpStatusError := http.StatusInternalServerError
	if (errors.As(err, &customErr.Error{})) {
		errCustom := err.(customErr.Error)
		if errCustom.HttpStatus != 0 {
			httpStatusError = errCustom.HttpStatus
		}
	}

	meta := Meta{
		Message:    err.Error(),
		HttpStatus: httpStatusError,
	}

	Respond(w, meta, httpStatusError)
	return err
}
