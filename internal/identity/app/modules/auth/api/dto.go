package api

import "time"

type DTOProvider struct {
	ID          int64     `json:"id"`
	Source      string    `json:"source"`
	Type        string    `json:"type"`
	DisplayName string    `json:"display_name"`
	Category    string    `json:"category"`
	Disabled    bool      `json:"disabled"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}
