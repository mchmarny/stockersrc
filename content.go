package main

import (
	"time"
)

// TextContent represents generic text event
type TextContent struct {
	Symbol    string    `json:"symbol"`
	ID        string    `json:"cid"`
	CreatedAt time.Time `json:"created"`
	Author    string    `json:"author"`
	Lang      string    `json:"lang"`
	Source    string    `json:"source"`
	Content   string    `json:"content"`
	Magnitude float32   `json:"magnitude"`
	Score     float32   `json:"score"`
	IsRetweet bool      `json:"retweet"`
}
