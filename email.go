package main

import "time"

type Email struct {
	ID      string    `json:"id"`
	To      string    `json:"to"`
	Subject string    `json:"subject"`
	Body    string    `json:"body"`
	SendAt  time.Time `json:"send_at"`
}
