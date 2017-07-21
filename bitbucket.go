package main

import (
	"fmt"
	"net/http"
)

type bitbucketBuildStatus struct {
	// author
	Repo struct {
		Name string `json:"name"`
	} `json:"repository"`
	CommitStatus struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		State       string `json:"state"`
		Key         string `json:"key"`
		Url         string `json:"url"`
		Refname     string `json:"refname"`
		Type        string `json:"type"`
		CreatedOn   string `json:"created_on"`
		UpdatedOn   string `json:"updated_on"`
		Commit      struct {
			Message string `json:"message"`
		} `json:"commit"`
	} `json:"commit_status"`
}

func bitbucketHandler(r *http.Request) *discordWebhook {
	buildStatus := new(bitbucketBuildStatus)

	if tryDecodeJson(buildStatus, r) == false {
		return nil
	}

	discordPayload := new(discordWebhook)
	discordPayload.Content = fmt.Sprintf("%v build on %v: %v", buildStatus.Repo.Name, buildStatus.CommitStatus.Refname, buildStatus.CommitStatus.State)

	embed := new(discordEmbed)
	embed.Title = discordPayload.Content
	embed.Description = buildStatus.CommitStatus.Commit.Message
	embed.Url = buildStatus.CommitStatus.Url

	discordPayload.Embeds = []*discordEmbed{embed}

	return discordPayload
}
