// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: consensus_protocol.proto

package core_rpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ConsensusProtocolClient is the client API for ConsensusProtocol service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConsensusProtocolClient interface {
	OnProposeMessage(ctx context.Context, in *GProposeMessage, opts ...grpc.CallOption) (*GResult, error)
	OnVoteMessage(ctx context.Context, in *GPreVoteMessage, opts ...grpc.CallOption) (*GResult, error)
	StartConsensus(ctx context.Context, in *GRequest, opts ...grpc.CallOption) (*GResult, error)
}

type consensusProtocolClient struct {
	cc grpc.ClientConnInterface
}

func NewConsensusProtocolClient(cc grpc.ClientConnInterface) ConsensusProtocolClient {
	return &consensusProtocolClient{cc}
}

func (c *consensusProtocolClient) OnProposeMessage(ctx context.Context, in *GProposeMessage, opts ...grpc.CallOption) (*GResult, error) {
	out := new(GResult)
	err := c.cc.Invoke(ctx, "/corerpc.ConsensusProtocol/OnProposeMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consensusProtocolClient) OnVoteMessage(ctx context.Context, in *GPreVoteMessage, opts ...grpc.CallOption) (*GResult, error) {
	out := new(GResult)
	err := c.cc.Invoke(ctx, "/corerpc.ConsensusProtocol/OnVoteMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consensusProtocolClient) StartConsensus(ctx context.Context, in *GRequest, opts ...grpc.CallOption) (*GResult, error) {
	out := new(GResult)
	err := c.cc.Invoke(ctx, "/corerpc.ConsensusProtocol/StartConsensus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConsensusProtocolServer is the server API for ConsensusProtocol service.
// All implementations must embed UnimplementedConsensusProtocolServer
// for forward compatibility
type ConsensusProtocolServer interface {
	OnProposeMessage(context.Context, *GProposeMessage) (*GResult, error)
	OnVoteMessage(context.Context, *GPreVoteMessage) (*GResult, error)
	StartConsensus(context.Context, *GRequest) (*GResult, error)
	mustEmbedUnimplementedConsensusProtocolServer()
}

// UnimplementedConsensusProtocolServer must be embedded to have forward compatible implementations.
type UnimplementedConsensusProtocolServer struct {
}

func (UnimplementedConsensusProtocolServer) OnProposeMessage(context.Context, *GProposeMessage) (*GResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnProposeMessage not implemented")
}
func (UnimplementedConsensusProtocolServer) OnVoteMessage(context.Context, *GPreVoteMessage) (*GResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnVoteMessage not implemented")
}
func (UnimplementedConsensusProtocolServer) StartConsensus(context.Context, *GRequest) (*GResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartConsensus not implemented")
}
func (UnimplementedConsensusProtocolServer) mustEmbedUnimplementedConsensusProtocolServer() {}

// UnsafeConsensusProtocolServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConsensusProtocolServer will
// result in compilation errors.
type UnsafeConsensusProtocolServer interface {
	mustEmbedUnimplementedConsensusProtocolServer()
}

func RegisterConsensusProtocolServer(s grpc.ServiceRegistrar, srv ConsensusProtocolServer) {
	s.RegisterService(&ConsensusProtocol_ServiceDesc, srv)
}

func _ConsensusProtocol_OnProposeMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GProposeMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusProtocolServer).OnProposeMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/corerpc.ConsensusProtocol/OnProposeMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusProtocolServer).OnProposeMessage(ctx, req.(*GProposeMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConsensusProtocol_OnVoteMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GPreVoteMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusProtocolServer).OnVoteMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/corerpc.ConsensusProtocol/OnVoteMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusProtocolServer).OnVoteMessage(ctx, req.(*GPreVoteMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConsensusProtocol_StartConsensus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusProtocolServer).StartConsensus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/corerpc.ConsensusProtocol/StartConsensus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusProtocolServer).StartConsensus(ctx, req.(*GRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ConsensusProtocol_ServiceDesc is the grpc.ServiceDesc for ConsensusProtocol service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConsensusProtocol_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "corerpc.ConsensusProtocol",
	HandlerType: (*ConsensusProtocolServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OnProposeMessage",
			Handler:    _ConsensusProtocol_OnProposeMessage_Handler,
		},
		{
			MethodName: "OnVoteMessage",
			Handler:    _ConsensusProtocol_OnVoteMessage_Handler,
		},
		{
			MethodName: "StartConsensus",
			Handler:    _ConsensusProtocol_StartConsensus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "consensus_protocol.proto",
}