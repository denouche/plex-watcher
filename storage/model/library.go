package model

import "time"

type Library struct {
	ID                       string     `json:"id"`
	OwnerID                  string     `json:"ownerId"`
	TeamID                   string     `json:"teamId"`
	ChannelID                string     `json:"channelId"`
	BaseURL                  string     `json:"baseUrl"`
	Token                    string     `json:"token"`
	LastNotifiedMediaAddedAt *time.Time `json:"lastNotifiedMediaAddedAt"`
}
