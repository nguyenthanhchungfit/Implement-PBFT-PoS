package core_rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type GClientConfig struct {
	Host string
	Port int16
}

type GClient struct {
	ClientConfig GClientConfig
	connection grpc.ClientConn
	innerClient ConsensusProtocolClient
}

func (client *GClient) SendProposeMessage(message *GProposeMessage) int32{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var opts []grpc.CallOption
	result, error := client.innerClient.OnProposeMessage(ctx, message, opts...)
	if error != nil {
		log.Println("sendProposeMessage failed %d", error)
		return -1
	}
	return result.Error
}

func NewGClient(cfg GClientConfig) *GClient {
	client := GClient{ClientConfig: cfg}
	serverAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	//var opts []grpc.DialOption
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	//defer conn.Close()
	client.innerClient = NewConsensusProtocolClient(conn)
	return &client
}