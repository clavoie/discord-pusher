package main

import (
	"html/template"
	"net/http"

	"github.com/clavoie/discord-pusher-deps/types"
)

var rootTemplate = template.Must(template.ParseFiles("index.html"))

type rootParams struct {
	Error     string
	Hooks     []*hookParam
	HookTypes []string
}

type hookParam struct {
	HookUrl string
	Key     string
	Type    string
	Url     string
}

func newParams(keys []string, dals []*types.HookDal) []*hookParam {
	hooks := make([]*hookParam, len(keys))

	for index, key := range keys {
		dal := dals[index]
		hookUrl := hookHandlerUrl + dal.Hook
		hook := &hookParam{hookUrl, key, dal.Type, dal.DiscordHook}
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
	params.Hooks = newParams(allHooks(contextFn(r)))
	params.HookTypes = hookTypes()

	rootTemplate.Execute(w, params)
}

func redirectToRoot(err string, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?e="+err, http.StatusFound)
}
