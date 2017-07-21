package main

import (
	"html/template"
	"net/http"

	"google.golang.org/appengine/datastore"
)

var rootTemplate = template.Must(template.ParseFiles("index.html"))

type rootParams struct {
	Error     string
	Hooks     []*hookParam
	HookTypes []string
}

type hookParam struct {
	Key  string
	Type string
	Url  string
}

func newParams(keys []*datastore.Key, dals []*hookDal) []*hookParam {
	hooks := make([]*hookParam, len(keys))

	for index, key := range keys {
		dal := dals[index]
		hook := &hookParam{key.Encode(), dal.Type, dal.DiscordHook}
		hooks[index] = hook
	}

	return hooks
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := r.URL.Query().Get("e")

	params := new(rootParams)
	params.Error = err
	params.Hooks = newParams(allHooks(r))
	params.HookTypes = hookTypes()

	rootTemplate.Execute(w, params)
}

func redirectToRoot(err string, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?e="+err, http.StatusFound)
}
