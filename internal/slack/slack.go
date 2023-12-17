// internal/slack/slack.go
package slack

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
	"github.com/stackzoo/voltbot/internal/lightning"
)

// Run starts the Slack integration.
func Run() {
	// Replace with your Slack token
	slackToken := "toke-here"
	channelID := "C06AF4FQJN9"

	api := slack.New(slackToken)

	// Get Lightning node info
	nodeInfo, err := lightning.GetNodeInfo()
	if err != nil {
		log.Fatalf("Error retrieving node info: %v", err)
	}

	// Format message
	message := fmt.Sprintf("Node Alias: %s\nNode Public Key: %s\nNode Active channels: %v\n",
		nodeInfo.IdentityAddress, nodeInfo.LightningId, nodeInfo.NumActiveChannels)

	// Post message to Slack channel
	_, _, err = api.PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		log.Fatalf("Error posting message to Slack: %v", err)
	}
}
