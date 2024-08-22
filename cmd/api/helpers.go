package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"harry2an.com/expenses/internal/validator"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dest any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(2000000))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&dest)
	if err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		var destinationErr *json.InvalidUnmarshalError

		switch {
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("malformed payload")
		case errors.Is(err, io.EOF):
			return errors.New("payload not found")
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("syntax error at char %d", syntaxErr.Offset)
		case err.Error() == "http: request body too large":
			return errors.New("payload too large")
		case errors.As(err, &typeErr):
			return errors.New("unmatched type in payload")
		case errors.As(err, &destinationErr):
			return errors.New("destination is not a pointer")
		case strings.HasPrefix(err.Error(), "json: unknown fields "):
			return fmt.Errorf("unknown field %s", strings.TrimPrefix(err.Error(), "json: unknown fields "))
		default:
			return err
		}
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("multiple payloads detected")
	}

	return nil
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id param")
	}

	return id, nil
}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.Println(fmt.Errorf("%s", err))
			}
		}()

		fn()
	}()
}

func (app *application) readInt(q url.Values, key string, defaultVal int, v *validator.Validator) int {
	val := q.Get(key)
	if val == "" {
		return defaultVal
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		v.AddError(key, "must be an integer")
		return defaultVal
	}

	return i
}
