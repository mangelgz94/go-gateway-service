// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: find-number-position-api.proto

package find_number_position_api

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

// FindNumberPositionAPIServiceClient is the client API for FindNumberPositionAPIService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FindNumberPositionAPIServiceClient interface {
	GetNumberPosition(ctx context.Context, in *GetNumberPositionRequest, opts ...grpc.CallOption) (*GetNumberPositionResponse, error)
}

type findNumberPositionAPIServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFindNumberPositionAPIServiceClient(cc grpc.ClientConnInterface) FindNumberPositionAPIServiceClient {
	return &findNumberPositionAPIServiceClient{cc}
}

func (c *findNumberPositionAPIServiceClient) GetNumberPosition(ctx context.Context, in *GetNumberPositionRequest, opts ...grpc.CallOption) (*GetNumberPositionResponse, error) {
	out := new(GetNumberPositionResponse)
	err := c.cc.Invoke(ctx, "/find_number_position_api.FindNumberPositionAPIService/GetNumberPosition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FindNumberPositionAPIServiceServer is the server API for FindNumberPositionAPIService service.
// All implementations must embed UnimplementedFindNumberPositionAPIServiceServer
// for forward compatibility
type FindNumberPositionAPIServiceServer interface {
	GetNumberPosition(context.Context, *GetNumberPositionRequest) (*GetNumberPositionResponse, error)
	mustEmbedUnimplementedFindNumberPositionAPIServiceServer()
}

// UnimplementedFindNumberPositionAPIServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFindNumberPositionAPIServiceServer struct {
}

func (UnimplementedFindNumberPositionAPIServiceServer) GetNumberPosition(context.Context, *GetNumberPositionRequest) (*GetNumberPositionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNumberPosition not implemented")
}
func (UnimplementedFindNumberPositionAPIServiceServer) mustEmbedUnimplementedFindNumberPositionAPIServiceServer() {
}

// UnsafeFindNumberPositionAPIServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FindNumberPositionAPIServiceServer will
// result in compilation errors.
type UnsafeFindNumberPositionAPIServiceServer interface {
	mustEmbedUnimplementedFindNumberPositionAPIServiceServer()
}

func RegisterFindNumberPositionAPIServiceServer(s grpc.ServiceRegistrar, srv FindNumberPositionAPIServiceServer) {
	s.RegisterService(&FindNumberPositionAPIService_ServiceDesc, srv)
}

func _FindNumberPositionAPIService_GetNumberPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNumberPositionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FindNumberPositionAPIServiceServer).GetNumberPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/find_number_position_api.FindNumberPositionAPIService/GetNumberPosition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FindNumberPositionAPIServiceServer).GetNumberPosition(ctx, req.(*GetNumberPositionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FindNumberPositionAPIService_ServiceDesc is the grpc.ServiceDesc for FindNumberPositionAPIService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FindNumberPositionAPIService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "find_number_position_api.FindNumberPositionAPIService",
	HandlerType: (*FindNumberPositionAPIServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNumberPosition",
			Handler:    _FindNumberPositionAPIService_GetNumberPosition_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "find-number-position-api.proto",
}