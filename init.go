package main

import (
	"net/http"
	"os"

	"github.com/clavoie/discord-pusher-deps/ae"
	"github.com/clavoie/discord-pusher-deps/types"
)

type newContextFn func(*http.Request) types.HookContext

var contextFn newContextFn

var pusherTypes = map[string]newContextFn{
	"app_engine": ae.NewHookContext,
}

func init() {
	setContextFn()

	http.HandleFunc(hookHandlerUrl, hookHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/", rootHandler)
}

func setContextFn() {
	applicationType := os.Getenv("DISCORD_PUSHER_TYPE")

	fn, hasFn := pusherTypes[applicationType]

	if hasFn {
		contextFn = fn
		return
	}

	if len(os.Args) > 0 {
		fn, hasFn = pusherTypes[os.Args[0]]

		if hasFn {
			contextFn = fn
			return
		}
	}

	contextFn = ae.NewHookContext
}
