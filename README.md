# VOLTBOT
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/stackzoo/voltbot)](https://goreportcard.com/report/github.com/stackzoo/voltbot)  
<img src="images/voltbot-logo-nobg.png" alt="logo" width="160"/>

Lightning node bot ⚡🤖

### Supported Products

- [Slack](https://slack.com/)

### Supported Lightning Nodes Implementations

- [LND](https://github.com/lightningnetwork/lnd)

## Abstract
Voltbot is a lightweight bot that retrieves lightning node stats and send them via different channels.  
This can be beneficial for executing and monitoring nodes.  
Voltbot communicates with the *LND* instance through [gRPC](https://grpc.io/).  




## Instructions

To operate, the bot requires reading configuration data from a JSON file.  
The configuration file must be located within the `config` folder at the root directory and must be named `voltbot_config.json`.  
You can take a look at the example file inside the `config` folder:  
```json
{
    "lnd_node_endpoint": "<node-endpoint:port>",
    "lnd_node_tls_cert_path": "config/tls.cert",
    "lnd_node_macaroon_hex_path": "config/macaroon.hex",
    "slack_channel_id": "<slack-channel-id>",
    "slack_token": "<slack-token>"
}
```  
At present, the retrieval of statistics is hard-coded to occur every *1440 minutes* (24 hours).  

## Example

<img src="images/slack-message.png" alt="logo" width="900"/>


