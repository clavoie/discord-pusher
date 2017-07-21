package main

import (
	"fmt"
	"net/http"
)

type unityBuildStatus struct {
	ProjectName     string `json:"projectName"`
	BuildTargetName string `json:"buildTargetName"`
	BuildNumber     int    `json:"buildNumber"`
	BuildStatus     string `json:"buildStatus"`
	Links           struct {
		DashboardUrl struct {
			Href string `json:"href"`
		} `json:"dashboard_url"`
		DashboardSummary struct {
			Href string `json:"href"`
		} `json:"dashboard_summary"`
	} `json:"links"`
}

func unityHandler(r *http.Request) *discordWebhook {
	buildStatus := new(unityBuildStatus)

	if tryDecodeJson(buildStatus, r) == false {
		return nil
	}

	discordPayload := new(discordWebhook)
	discordPayload.Content = buildStatus.ProjectName

	embed := new(discordEmbed)
	embed.Title = fmt.Sprintf("#%v - %v", buildStatus.BuildNumber, buildStatus.BuildTargetName)
	embed.Description = buildStatus.BuildStatus
	embed.Url = buildStatus.Links.DashboardUrl.Href + buildStatus.Links.DashboardSummary.Href

	discordPayload.Embeds = []*discordEmbed{embed}

	return discordPayload
}
