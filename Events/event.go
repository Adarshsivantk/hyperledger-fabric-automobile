
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// blockEventListener listens for new blocks on a channel.
func blockEventListener(organization string, channelName string) {
	orgProfile := profile[organization]

	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

	// Create a gRPC client connection.
	clientConnection := newGrpcConnection(
		tlsCertPath,
		gatewayPeer,
		peerEndpoint,
	)
	defer clientConnection.Close()

	id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

	// Connect to the Fabric Gateway.
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)

	// Create a context for event listening.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("\n*** Start Block event listening")

	events, err := network.BlockEvents(
		ctx,
		client.WithStartBlock(1),
	)
	if err != nil {
		panic(fmt.Errorf(
			"failed to start Block event listening: %w",
			err,
		))
	}

	for event := range events {
		fmt.Printf(
			"Received block number %d\n",
			event.GetHeader().GetNumber(),
		)
	}
}

// chaincodeEventListener listens for chaincode events on a channel.
func chaincodeEventListener(
	organization string,
	channelName string,
	chaincodeName string,
) {
	orgProfile := profile[organization]

	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

	// Create a gRPC client connection.
	clientConnection := newGrpcConnection(
		tlsCertPath,
		gatewayPeer,
		peerEndpoint,
	)
	defer clientConnection.Close()

	id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

	// Connect to the Fabric Gateway.
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)

	// Create a context for event listening.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("*** Start Chaincode event listening")

	events, err := network.ChaincodeEvents(ctx, chaincodeName)
	if err != nil {
		panic(fmt.Errorf(
			"failed to start Chaincode event listening: %w",
			err,
		))
	}

	for event := range events {
		fmt.Printf(
			"Received Chaincode Event: %s\nData: %s\n",
			event.EventName,
			event.Payload,
		)
	}
}

// pvtBlockEventListener listens for blocks containing private data.
func pvtBlockEventListener(organization string, channelName string) {
	orgProfile := profile[organization]

	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

	// Create a gRPC client connection.
	clientConnection := newGrpcConnection(
		tlsCertPath,
		gatewayPeer,
		peerEndpoint,
	)
	defer clientConnection.Close()

	id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

	// Connect to the Fabric Gateway.
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)

	// Create a context for event listening.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("*** Start Block and private data event listening")

	events, err := network.BlockAndPrivateDataEvents(
		ctx,
		client.WithStartBlock(1),
	)
	if err != nil {
		panic(fmt.Errorf(
			"failed to start Block event listening: %w",
			err,
		))
	}

	for event := range events {
		if event.GetPrivateDataMap() != nil {
			fmt.Printf(
				"Received block: %d with private data\n",
				event.GetBlock().GetHeader().GetNumber(),
			)
		}
	}
}