package routes

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ParseResponse(resp *http.Response, schema interface{}) error {
	var unmarshalErr *json.UnmarshalTypeError

	// headerContentTtype := resp.Header.Get("Content-Type")
	// fmt.Println(headerContentTtype)
	// if headerContentTtype != "application/json" {
	// 	return errors.New("content type is not application/json")
	// }

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields() //throws error if uneeded JSON is added
	err := decoder.Decode(schema)   //decodes the incoming JSON into the struct
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			return errors.New("Bad Request. Wrong Type provided for field " + unmarshalErr.Field)
		} else {
			return errors.New("Bad Request " + err.Error())
		}
	}
	return nil
}
