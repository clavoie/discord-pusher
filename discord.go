package main

type discordWebhook struct {
	Content   string          `json:"content"`
	Username  string          `json:"username"`
	AvatarUrl string          `json:"avatar_url"`
	Tts       bool            `json:"tts"`
	Embeds    []*discordEmbed `json:"embeds"`
}

type discordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}
