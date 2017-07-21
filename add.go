package main

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := r.ParseForm()

	if err != nil {
		log.Errorf(appengine.NewContext(r), "%v", err)
		redirectToRoot("Could not parse the request", w, r)
		return
	}

	userError := addHook(r.PostForm.Get("type"), r.PostForm.Get("url"), r)
	redirectToRoot(userError, w, r)
}
