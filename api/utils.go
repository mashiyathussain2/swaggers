package api

import (
	"encoding/json"

	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// DecodeJSONBody decode json data from request to an interface
func (a *API) DecodeJSONBody(r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		if r.Header.Get("Content-Type") != "application/json" {
			err := errors.New("unsupported content-type request: Content-Type header is not application/json")
			// ctx.SetErr(err)
			return err
		}
	}

	if r.ContentLength == 0 {
		return errors.New("Request body must not be empty")
	}

	// r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			// ctx.SetErr(errors.New(msg))
			return errors.New(msg)
			// return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			// ctx.SetErr(errors.New(msg))
			return errors.New(msg)
			// return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			// ctx.SetErr(errors.New(msg))
			return errors.New(msg)
			// return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			// ctx.SetErr(errors.New(msg))
			return errors.New(msg)
			// return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			// ctx.SetErr(errors.New(msg))
			return errors.New(msg)
			// return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			// ctx.SetErr(errors.New(msg))
			return errors.New(msg)
			// return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			// ctx.SetErr(errors.New(err.Error()))
			return errors.New(err.Error())
			// return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		// ctx.SetErr(errors.New(msg))
		return errors.New(msg)
		// return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}
	return nil
}
