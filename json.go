package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func tryDecodeJson(dst interface{}, r *http.Request) bool {
	defer func(r *http.Request) {
		err := r.Body.Close()

		if err != nil {
			log.Errorf(appengine.NewContext(r), "could not close request body: %v", err)
		}
	}(r)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dst)

	if err != nil {
		log.Errorf(appengine.NewContext(r), "could not decode json: %v", err)
		return false
	}

	return true
}

func tryEncodeJson(payload interface{}) (io.Reader, error) {
	var buffer bytes.Buffer

	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(payload)

	if err != nil {
		return nil, err
	}

	return &buffer, nil
}
