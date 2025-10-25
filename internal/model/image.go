package model

import "time"

type Image struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	URL       string    `json:"url"`
	DrawURL   string    `json:"draw_url"`
	CreatedAt time.Time `json:"created_at"`
}
