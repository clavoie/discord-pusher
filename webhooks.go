package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	uuid "github.com/nu7hatch/gouuid"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

const webhookKind = "Webhook"

type hookHandlerFn func(*http.Request) *discordWebhook

var webhooks = map[string]hookHandlerFn{
	"bitbucket": bitbucketHandler,
	"unity":     unityHandler,
}

type hookDal struct {
	DiscordHook string
	Hook        string
	Type        string
}

func hookTypes() []string {
	types := make([]string, 0, len(webhooks))

	for hookName := range webhooks {
		types = append(types, hookName)
	}

	sort.Strings(types)
	return types
}

func addHook(typeName, webhookUrl string, r *http.Request) string {
	appContext := appengine.NewContext(r)
	userError := ""

	_, hasWebhookType := webhooks[typeName]

	if hasWebhookType == false {
		return fmt.Sprintf("Unknown webhook type: %v", typeName)
	}

	_, err := url.ParseRequestURI(webhookUrl)

	if err != nil {
		log.Errorf(appContext, "%v", err)
		return fmt.Sprintf("Invalid webhook url: %v", webhookUrl)
	}

	q := datastore.NewQuery(webhookKind).Filter("DiscordHook = ", webhookUrl).Filter("Type =", typeName).Limit(1)
	dals := make([]*hookDal, 0, 1)
	_, err = q.GetAll(appContext, &dals)

	if err != nil {
		log.Errorf(appContext, "%v", err)
		return "An error occured while communicating with the datastore"
	}

	if len(dals) > 0 {
		return "A webhook of that type and url already exist"
	}

	key := datastore.NewIncompleteKey(appContext, webhookKind, nil)
	id, err := uuid.NewV4()

	if err != nil {
		log.Errorf(appContext, "%v", err)
		return "An error occured while generating a new webhook"
	}

	dal := &hookDal{webhookUrl, strings.Replace(id.String(), "-", "", -1), typeName}
	_, err = datastore.Put(appContext, key, dal)

	if err != nil {
		log.Errorf(appContext, "%v", err)
		return "An error occured while communicating with the datastore"
	}

	time.Sleep(time.Second)
	return userError
}

func allHooks(r *http.Request) ([]*datastore.Key, []*hookDal) {
	appContext := appengine.NewContext(r)
	q := datastore.NewQuery(webhookKind)
	dals := make([]*hookDal, 0, 10)
	keys, err := q.GetAll(appContext, &dals)

	if err != nil {
		log.Errorf(appContext, "%v", err)
	}

	return keys, dals
}

func deleteHook(encodedKey string, r *http.Request) string {
	appContext := appengine.NewContext(r)
	key, err := datastore.DecodeKey(encodedKey)

	if err != nil {
		log.Errorf(appContext, "%v", err)
		return "Invalid webhook key"
	}

	err = datastore.Delete(appContext, key)

	if err != nil {
		log.Errorf(appContext, "%v", err)
		return "Error communicating with the db"
	}

	time.Sleep(time.Second)
	return ""
}

func dispatchHook(hookId string, r *http.Request) {
	appContext := appengine.NewContext(r)
	q := datastore.NewQuery(webhookKind).Filter("Hook =", hookId).Limit(1)
	dals := make([]*hookDal, 0, 1)
	_, err := q.GetAll(appContext, &dals)

	if err != nil {
		log.Errorf(appContext, "%v", err)
		return
	}

	if len(dals) < 1 {
		return
	}

	dal := dals[0]
	handler, hasHandler := webhooks[dal.Type]

	if hasHandler == false {
		log.Errorf(appContext, "unknown handler type: %v", dal.Type)
		return
	}

	discordPayload := handler(r)

	if discordPayload == nil {
		return
	}

	reader, err := tryEncodeJson(discordPayload)

	if err != nil {
		log.Errorf(appContext, "could not encode payload: %v", err)
		return
	}

	client := urlfetch.Client(appContext)
	resp, err := client.Post(dal.DiscordHook, "application/json", reader)

	if err != nil {
		log.Errorf(appContext, "post to discord failed: %v", err)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		dispatchErrHandler(resp, r)
	}
}

func dispatchErrHandler(resp *http.Response, r *http.Request) {
	appContext := appengine.NewContext(r)
	log.Errorf(appContext, "discord response status: %v", resp.StatusCode)

	defer func(r *http.Request) {
		err := resp.Body.Close()

		if err != nil {
			log.Errorf(appengine.NewContext(r), "could not close discord body: %v", err)
		}
	}(r)

	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(resp.Body)

	if err != nil {
		log.Errorf(appContext, "could not read from discord response body: %v", err)
		return
	}

	log.Errorf(appContext, "discord err body: %v", buffer.String())
}
