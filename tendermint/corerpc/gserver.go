package core_rpc

import (
	"context"
	"fmt"
	"github.com/implement-pbft-pos/tendermint/utils"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GServerConfig struct {
	Id   int32
	Port uint16
}

type IProcessor interface {
	ReceiveProposeMessage(height, round, validRound int32, data string)
	ReceivePreVoteMessage(height, round int32, hashValue []byte)
	ReceivePreCommitMessage(height, round int32, hashValue []byte)
}

type GServer struct {
	ServerCfg GServerConfig
	Processor IProcessor
}

func (s *GServer) StartServer() {
	port := s.ServerCfg.Port
	id := s.ServerCfg.Id
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		utils.ErrorStdOutLogger.Printf("Server %d failed to listen at port %d with error: %v", id, port, err)
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("Server %d start listen at port: %d", id, port)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	handler := &consensusProtocolServer{Processor: s.Processor}
	RegisterConsensusProtocolServer(grpcServer, handler)
	grpcServer.Serve(lis)
}

type consensusProtocolServer struct {
	UnimplementedConsensusProtocolServer
	Processor IProcessor
}

func (s *consensusProtocolServer) OnProposeMessage(ctx context.Context, message *GProposeMessage) (*GResult, error) {
	s.Processor.ReceiveProposeMessage(message.Height, message.Round, message.ValidRound, message.Data.Data)
	return &GResult{Error: 0, Data: "onProposeMessage"}, nil
}

func (s *consensusProtocolServer) OnPreVoteMessage(ctx context.Context, message *GPreVoteMessage) (*GResult, error) {
	s.Processor.ReceivePreVoteMessage(message.Height, message.Round, message.HashValue)
	return &GResult{Error: 0, Data: "OnPreVoteMessage"}, nil
}

func (s *consensusProtocolServer) OnPreCommitMessage(ctx context.Context, message *GPreCommitMessage) (*GResult, error) {
	s.Processor.ReceivePreCommitMessage(message.Height, message.Round, message.HashValue)
	return &GResult{Error: 0, Data: "OnPreCommitMessage"}, nil
}

