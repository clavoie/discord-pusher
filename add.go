package main

import (
	"net/http"
)

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	hc := contextFn(r)

	if err != nil {
		hc.Errorf("%v", err)
		redirectToRoot("Could not parse the request", w, r)
		return
	}

	userError := addHook(r.PostForm.Get("type"), r.PostForm.Get("url"), hc)
	redirectToRoot(userError, w, r)
}
