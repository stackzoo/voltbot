// internal/lightning/lightning.go
package lightning

import (
	"context"
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

// Run establishes a connection to the Lightning node and prints node information.
func Run() {
	// Replace these values with your Lightning node's details
	rpcServerAddress := "umbrel.local:10009"
	tlsCertPath := "tls.cert"
	macaroonPath := "macaroon.hex" // Update with your actual macaroon filename

	// Load the TLS certificate for a secure connection
	creds, err := credentials.NewClientTLSFromFile(tlsCertPath, "")
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
	macaroonBytes, err := os.ReadFile(macaroonPath)
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

	// Print the retrieved information
	log.Printf("Node Alias: %s\n", nodeInfo.IdentityAddress)
	log.Printf("Node Public Key: %s\n", nodeInfo.LightningId)
	log.Printf("Node Active channels: %v\n", nodeInfo.NumActiveChannels)
}

// GetNodeInfo retrieves information about the Lightning node.
func GetNodeInfo() (*lnrpc.GetInfoResponse, error) {
	// Replace these values with your Lightning node's details
	rpcServerAddress := "umbrel.local:10009"
	tlsCertPath := "tls.cert"
	macaroonPath := "macaroon.hex" // Update with your actual macaroon filename

	// Load the TLS certificate for a secure connection
	creds, err := credentials.NewClientTLSFromFile(tlsCertPath, "")
	if err != nil {
		return nil, err
	}

	// Create gRPC dial options with TLS credentials
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	// Create a gRPC connection to the Lightning node
	conn, err := grpc.Dial(rpcServerAddress, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Load the macaroon for authentication
	macaroonBytes, err := os.ReadFile(macaroonPath)
	if err != nil {
		return nil, err
	}

	// Create a new gRPC context with the macaroon
	ctx := context.Background()
	perRPCCredentials := NewMacaroonCredentials(macaroonBytes)
	opts = append(opts, grpc.WithPerRPCCredentials(perRPCCredentials))

	conn, err = grpc.Dial(rpcServerAddress, opts...)
	if err != nil {
		return nil, err
	}

	// Create an lnrpc Lightning client
	client := lnrpc.NewLightningClient(conn)

	// Get node info
	nodeInfo, err := client.GetInfo(ctx, &lnrpc.GetInfoRequest{})
	if err != nil {
		return nil, err
	}

	return nodeInfo, nil
}
