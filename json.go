package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func tryDecodeJson(dst interface{}, r *http.Request) bool {
	hc := contextFn(r)

	defer func() {
		err := r.Body.Close()

		if err != nil {
			hc.Errorf("could not close request body: %v", err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dst)

	if err != nil {
		hc.Errorf("could not decode json: %v", err)
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
