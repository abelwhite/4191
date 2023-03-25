// Filename: ./cmd/api/helpers.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{} //envelope is a type that will contain school infos //container inside a container

func (app *application) readIDParams(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid Id Parameter")
	}
	return id, nil
}

func (app *application) WriteJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	//convert the data to json format
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n') //if no error
	//add any headers that were sent
	for key, value := range headers {
		w.Header()[key] = value
	}

	//add header info
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(js))
	return nil

}

// to help signify the error from the api
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, destination interface{}) error {
	//we will try to decode a JSON request
	err := json.NewDecoder(r.Body).Decode(destination)
	if err != nil {
		//something went wrong
		//lets find out the type of error
		var syntaxError *json.SyntaxError
		var unMarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		//lets check for the type of decode error
		switch {
		case errors.As(err, &syntaxError): //we use As if we looking for the type of error
			return fmt.Errorf("body contains badly-formed JSON(at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF): //we use Is if we looking for a specific type of error
			return fmt.Errorf("body contains badly-formed JSON")

		case errors.As(err, &unMarshalTypeError):
			if unMarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q)", unMarshalTypeError.Field)
			}
			return fmt.Errorf("body contains empty JSON field %q)", unMarshalTypeError.Field)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}

	}
	return nil

}
