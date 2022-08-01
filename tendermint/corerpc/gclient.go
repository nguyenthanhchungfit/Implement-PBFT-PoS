package core_rpc

import (
	"context"
	"fmt"
	"github.com/implement-pbft-pos/tendermint/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type GClientConfig struct {
	Host string
	Port uint16
	Timeout uint64
}

func (clientCfg GClientConfig) String() string {
	return fmt.Sprintf("GClientConf(Host: %s, Port: %d)", clientCfg.Host, clientCfg.Port)
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
		utils.ErrorStdOutLogger.Printf("sendProposeMessage failed %d", error)
		return -1
	}
	return result.Error
}

func (client *GClient) SendPreVoteMessage(message *GPreVoteMessage) int32{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var opts []grpc.CallOption
	result, error := client.innerClient.OnPreVoteMessage(ctx, message, opts...)
	if error != nil {
		utils.ErrorStdOutLogger.Printf("SendPreVoteMessage failed %d", error)
		return -1
	}
	return result.Error
}

func (client *GClient) SendPreCommitMessage(message *GPreCommitMessage) int32{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var opts []grpc.CallOption
	result, error := client.innerClient.OnPreCommitMessage(ctx, message, opts...)
	if error != nil {
		utils.ErrorStdOutLogger.Printf("SendPreCommitMessage failed %d", error)
		return -1
	}
	return result.Error
}

func (client *GClient) ConnectToRemote()  {
	host := client.ClientConfig.Host
	port := client.ClientConfig.Port
	serverAddr := fmt.Sprintf("%s:%d", host, port)
	//var opts []grpc.DialOption
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	//defer conn.Close()
	client.innerClient = NewConsensusProtocolClient(conn)
}