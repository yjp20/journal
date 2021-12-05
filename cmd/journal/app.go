package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type App struct {
	Addr     string
	DB       *sql.DB
	Password string
	Models   Models
	Origins  []string
	errorLog *log.Logger
}

func (app *App) writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	w.Write(b)
	return nil
}

func (app *App) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	return app.parseJSON(http.MaxBytesReader(w, r.Body, int64(maxBytes)), dst)
}

func (app *App) parseJSON(reader io.Reader, dst interface{}) error {
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError) && unmarshalTypeError.Field != "":
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		case errors.As(err, &unmarshalTypeError) && unmarshalTypeError.Field == "":
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *App) readOptionalTime(qs url.Values, key string) (*time.Time, error) {
	s := qs.Get(key)

	if s == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, s)
	return &t, err
}
