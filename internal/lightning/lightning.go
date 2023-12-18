// internal/lightning/lightning.go
package lightning

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type macaroonCred struct {
	macaroon []byte
}

func (m *macaroonCred) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"macaroon": string(m.macaroon),
	}, nil
}

func (m *macaroonCred) RequireTransportSecurity() bool {
	return true
}

// NewMacaroonCredentials creates a new gRPC per-RPC credential with the provided macaroon.
func NewMacaroonCredentials(macaroon []byte) credentials.PerRPCCredentials {
	return &macaroonCred{
		macaroon: macaroon,
	}
}

// Config represents the configuration for Lightning node connection.
type Config struct {
	RPCServerAddress string `json:"lnd_node_endpoint"`
	TLSCert          string `json:"lnd_node_tls_cert_path"`
	Macaroon         string `json:"lnd_node_macaroon_hex_path"`
	SlackToken       string `json:"slack_token"`
	SlackChannelID   string `json:"slack_channel_id"`
}

// LoadConfig reads the configuration from the voltbot_config.json file.
func LoadConfig() (*Config, error) {
	configFile := "config/voltbot_config.json"
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Run establishes a connection to the Lightning node and prints node information.
func Run() (*Config, *lnrpc.GetInfoResponse, error) {
	// Load the configuration
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Use config values
	rpcServerAddress := config.RPCServerAddress
	TLSCert := config.TLSCert
	Macaroon := config.Macaroon

	// Load the TLS certificate for a secure connection
	creds, err := credentials.NewClientTLSFromFile(TLSCert, "")
	if err != nil {
		log.Fatalf("Error loading TLS certificate: %v", err)
	}

	// Create gRPC dial options with TLS credentials
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	// Create a gRPC connection to the Lightning node
	conn, err := grpc.Dial(rpcServerAddress, opts...)
	if err != nil {
		log.Fatalf("Error connecting to Lightning node: %v", err)
	}
	defer conn.Close()

	// Load the macaroon for authentication
	macaroonBytes, err := os.ReadFile(Macaroon)
	if err != nil {
		log.Fatalf("Error reading macaroon file: %v", err)
	}

	// Create a new gRPC context with the macaroon
	ctx := context.Background()
	perRPCCredentials := NewMacaroonCredentials(macaroonBytes)
	opts = append(opts, grpc.WithPerRPCCredentials(perRPCCredentials))

	conn, err = grpc.Dial(rpcServerAddress, opts...)
	if err != nil {
		log.Fatalf("Error connecting to Lightning node with macaroon: %v", err)
	} else {
		log.Printf("Connection established")
	}
	defer conn.Close()

	// Create an lnrpc Lightning client
	client := lnrpc.NewLightningClient(conn)

	// Get node info
	nodeInfo, err := client.GetInfo(ctx, &lnrpc.GetInfoRequest{})
	if err != nil {
		log.Fatalf("Error retrieving node info: %v", err)
	}

	log.Printf("Node info: %v", nodeInfo)
	return config, nodeInfo, nil
}
