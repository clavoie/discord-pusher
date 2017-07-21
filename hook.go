package main

import "net/http"

const hookHandlerUrl = "/hook/"

func hookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	path := r.URL.Path
	sliceLen := len(hookHandlerUrl)

	if len(path) <= sliceLen {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	hookId := path[sliceLen:]
	dispatchHook(hookId, r)
}
