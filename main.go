package main

import (
	"time"

	"github.com/stackzoo/voltbot/internal/slack"
)

func main() {
	// Run slack.Run() immediately
	slack.Run()

	// Create a ticker that triggers every 6 hours
	ticker := time.NewTicker(360 * time.Minute)

	// Run the slack.Run() function each time the ticker triggers
	for range ticker.C {
		slack.Run()
	}
}
