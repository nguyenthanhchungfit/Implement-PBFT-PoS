package core_rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GServerConfig struct {
	Port int16
}

type GServer struct {
	ServerCfg GServerConfig
}

func (s *GServer) StartServer(){
	port := s.ServerCfg.Port;
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}else{
		log.Printf("Start listen at port: %d", port)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	handler := &consensusProtocolServer{}
	RegisterConsensusProtocolServer(grpcServer, handler)
	grpcServer.Serve(lis)
}

type consensusProtocolServer struct {
	UnimplementedConsensusProtocolServer
}

func (s *consensusProtocolServer) OnProposeMessage(ctx context.Context, message *GProposeMessage) (*GResult, error) {
	fmt.Printf("OnProposeMessage %d\n", message.Round)
	return &GResult{Error: 0, Data: "onProposeMessage"}, nil
}

func (s *consensusProtocolServer) OnPreVoteMessage(ctx context.Context, message *GPreVoteMessage) (*GResult, error) {
	fmt.Println("OnPreVoteMessage")
	return &GResult{Error: 0, Data: "OnPreVoteMessage"}, nil
}

func (s *consensusProtocolServer) OnPreCommitMessage(ctx context.Context, message *GPreCommitMessage) (*GResult, error) {
	fmt.Println("OnPreCommitMessage")
	return &GResult{Error: 0, Data: "OnPreCommitMessage"}, nil
}

//func main(){
//	flag.Parse()
//	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//	var opts []grpc.ServerOption
//	if *tls {
//		if *certFile == "" {
//			*certFile = data.Path("x509/server_cert.pem")
//		}
//		if *keyFile == "" {
//			*keyFile = data.Path("x509/server_key.pem")
//		}
//		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
//		if err != nil {
//			log.Fatalf("Failed to generate credentials %v", err)
//		}
//		opts = []grpc.ServerOption{grpc.Creds(creds)}
//	}
//	grpcServer := grpc.NewServer(opts...)
//	RegisterConsensusProtocolServer(grpcServer, newServer())
//	grpcServer.Serve(lis)
//}