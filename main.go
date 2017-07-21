package main

import (
	"net/http"
)

func init() {
	http.HandleFunc(hookHandlerUrl, hookHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/", rootHandler)
}
