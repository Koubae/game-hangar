package model

import "time"

type Provider struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Category    string    `json:"category"`
	Disabled    bool      `json:"disabled"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}
