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
	// Get Lightning node info
	config, nodeInfo, err := lightning.Run()
	if err != nil {
		log.Fatalf("Error retrieving node info: %v", err)
	}
	// Format message
	message := fmt.Sprintf("LND Version: %s\nNode Alias: %s\nNode Public Key: %s\nNode Active channels: %v\nNode Peers: %v\nBlock Height: %v\n",
		nodeInfo.Version, nodeInfo.Alias, nodeInfo.IdentityPubkey, nodeInfo.NumActiveChannels, nodeInfo.NumPeers, nodeInfo.BlockHeight)

	// Use config values
	slackToken := config.SlackToken
	channelID := config.SlackChannelID

	api := slack.New(slackToken)

	// Post message to Slack channel
	_, _, err = api.PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		log.Fatalf("Error posting message to Slack: %v", err)
	}

}
