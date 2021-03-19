package api

import (
	"encoding/json"
	"fmt"
	"go-app/model"
	"io"
	"net/http"
	"strconv"
	"strings"

	errors "github.com/vasupal1996/goerror"
	"go.mongodb.org/mongo-driver/bson/primitive"

	goErr "errors"
)

// DecodeJSONBody decode json data from request to an interface
func (a *API) DecodeJSONBody(r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		if r.Header.Get("Content-Type") != "application/json" {
			err := errors.New("unsupported content-type request: Content-Type header is not application/json", &errors.BadRequest)
			// ctx.SetErr(err)
			return err
		}
	}

	if r.ContentLength == 0 {
		return errors.New("Request body must not be empty", &errors.BadRequest)
	}

	// r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case goErr.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			// ctx.SetErr(errors.New(msg, &errors.BadRequest))
			return errors.New(msg, &errors.BadRequest)
			// return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case goErr.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			// ctx.SetErr(errors.New(msg, &errors.BadRequest))
			return errors.New(msg, &errors.BadRequest)
			// return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case goErr.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			// ctx.SetErr(errors.New(msg, &errors.BadRequest))
			return errors.New(msg, &errors.BadRequest)
			// return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			// ctx.SetErr(errors.New(msg, &errors.BadRequest))
			return errors.New(msg, &errors.BadRequest)
			// return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			// ctx.SetErr(errors.New(msg, &errors.BadRequest))
			return errors.New(msg, &errors.BadRequest)
			// return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			// ctx.SetErr(errors.New(msg, &errors.BadRequest))
			return errors.New(msg, &errors.BadRequest)
			// return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			// ctx.SetErr(errors.New(err.Error(), &errors.BadRequest))
			return errors.New(err.Error(), &errors.BadRequest)
			// return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		// ctx.SetErr(errors.New(msg, &errors.BadRequest))
		return errors.New(msg, &errors.BadRequest)
		// return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}
	return nil
}

//Extracts page value from the request
func GetPageValue(r *http.Request) int {
	var val int
	var err error
	if r.URL.Query().Get("page") != "" {
		val, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			val = 0
		}
	}
	return val
}

//Extracts status value from the request
func GetStatusValue(r *http.Request) string {
	var val string
	if r.URL.Query().Get("status") != "" {
		val = r.URL.Query().Get("status")
		if val == "" {
			val = model.Publish
		}
	}
	return val
}

//Extracts ObjectID from the request
func GetObjectID(r *http.Request) (primitive.ObjectID, error) {
	var val primitive.ObjectID
	var err error
	if r.URL.Query().Get("id") != "" {
		val, err = primitive.ObjectIDFromHex(r.URL.Query().Get("id"))
		if err != nil {
			return primitive.NilObjectID, err
		}
	}
	if r.URL.Query().Get("id") == "" {
		return primitive.NilObjectID, nil
	}
	return val, nil
}
