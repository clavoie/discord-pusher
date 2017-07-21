package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/clavoie/discord-pusher-deps/types"
	uuid "github.com/nu7hatch/gouuid"
)

type hookHandlerFn func(*http.Request) *discordWebhook

var webhooks = map[string]hookHandlerFn{
	"bitbucket": bitbucketHandler,
	"unity":     unityHandler,
}

func hookTypes() []string {
	types := make([]string, 0, len(webhooks))

	for hookName := range webhooks {
		types = append(types, hookName)
	}

	sort.Strings(types)
	return types
}

func addHook(typeName, webhookUrl string, hc types.HookContext) string {
	userError := ""

	_, hasWebhookType := webhooks[typeName]

	if hasWebhookType == false {
		return fmt.Sprintf("Unknown webhook type: %v", typeName)
	}

	_, err := url.ParseRequestURI(webhookUrl)

	if err != nil {
		hc.Errorf("%v", err)
		return fmt.Sprintf("Invalid webhook url: %v", webhookUrl)
	}

	dal, err := hc.GetByTypeUrl(typeName, webhookUrl)

	if err != nil {
		hc.Errorf("%v", err)
		return "An error occured while communicating with the datastore"
	}

	if dal != nil {
		return "A webhook of that type and url already exist"
	}

	id, err := uuid.NewV4()

	if err != nil {
		hc.Errorf("%v", err)
		return "An error occured while generating a new webhook"
	}

	dal = &types.HookDal{webhookUrl, strings.Replace(id.String(), "-", "", -1), typeName}
	err = hc.Put(dal)

	if err != nil {
		hc.Errorf("%v", err)
		return "An error occured while communicating with the datastore"
	}

	return userError
}

func allHooks(hc types.HookContext) ([]string, []*types.HookDal) {
	keys, dals, err := hc.GetAll()

	if err != nil {
		hc.Errorf("%v", err)
		return []string{}, []*types.HookDal{}
	}

	return keys, dals
}

func deleteHook(encodedKey string, hc types.HookContext) string {
	err := hc.Delete(encodedKey)

	if err != nil {
		hc.Errorf("%v", err)
		return "Could not delete webhook"
	}

	return ""
}

func dispatchHook(hookId string, r *http.Request) {
	hc := contextFn(r)
	dal, err := hc.GetByHook(hookId)

	if err != nil {
		hc.Errorf("%v", err)
		return
	}

	if dal == nil {
		return
	}

	handler, hasHandler := webhooks[dal.Type]

	if hasHandler == false {
		hc.Errorf("unknown handler type: %v", dal.Type)
		return
	}

	discordPayload := handler(r)

	if discordPayload == nil {
		return
	}

	reader, err := tryEncodeJson(discordPayload)

	if err != nil {
		hc.Errorf("could not encode payload: %v", err)
		return
	}

	resp, err := hc.UrlPost(dal, reader)

	if err != nil {
		hc.Errorf("post to discord failed: %v", err)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		dispatchErrHandler(resp, hc)
	}
}

func dispatchErrHandler(resp *http.Response, hc types.HookContext) {
	hc.Errorf("discord response status: %v", resp.StatusCode)

	defer func(hc types.HookContext) {
		err := resp.Body.Close()

		if err != nil {
			hc.Errorf("could not close discord body: %v", err)
		}
	}(hc)

	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(resp.Body)

	if err != nil {
		hc.Errorf("could not read from discord response body: %v", err)
		return
	}

	hc.Errorf("discord err body: %v", buffer.String())
}
