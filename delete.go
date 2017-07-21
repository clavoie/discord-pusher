package main

import "net/http"

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	encodedKey := r.URL.Query().Get("k")
	userError := deleteHook(encodedKey, r)
	redirectToRoot(userError, w, r)
}
