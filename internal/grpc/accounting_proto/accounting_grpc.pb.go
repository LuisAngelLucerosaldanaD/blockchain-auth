// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: api/grpc/proto/accounting.proto

package accounting_proto

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

// AccountingServicesAccountingClient is the client API for AccountingServicesAccounting service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountingServicesAccountingClient interface {
	GetAccountingByWalletById(ctx context.Context, in *RequestGetAccountingByWalletId, opts ...grpc.CallOption) (*ResponseGetAccountingByWalletId, error)
	CreateAccounting(ctx context.Context, in *RequestCreateAccounting, opts ...grpc.CallOption) (*ResponseCreateAccounting, error)
	SetAmountToAccounting(ctx context.Context, in *RequestSetAmountToAccounting, opts ...grpc.CallOption) (*ResponseSetAmountToAccounting, error)
}

type accountingServicesAccountingClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountingServicesAccountingClient(cc grpc.ClientConnInterface) AccountingServicesAccountingClient {
	return &accountingServicesAccountingClient{cc}
}

func (c *accountingServicesAccountingClient) GetAccountingByWalletById(ctx context.Context, in *RequestGetAccountingByWalletId, opts ...grpc.CallOption) (*ResponseGetAccountingByWalletId, error) {
	out := new(ResponseGetAccountingByWalletId)
	err := c.cc.Invoke(ctx, "/accounting_proto.accountingServicesAccounting/GetAccountingByWalletById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountingServicesAccountingClient) CreateAccounting(ctx context.Context, in *RequestCreateAccounting, opts ...grpc.CallOption) (*ResponseCreateAccounting, error) {
	out := new(ResponseCreateAccounting)
	err := c.cc.Invoke(ctx, "/accounting_proto.accountingServicesAccounting/CreateAccounting", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountingServicesAccountingClient) SetAmountToAccounting(ctx context.Context, in *RequestSetAmountToAccounting, opts ...grpc.CallOption) (*ResponseSetAmountToAccounting, error) {
	out := new(ResponseSetAmountToAccounting)
	err := c.cc.Invoke(ctx, "/accounting_proto.accountingServicesAccounting/SetAmountToAccounting", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountingServicesAccountingServer is the server API for AccountingServicesAccounting service.
// All implementations should embed UnimplementedAccountingServicesAccountingServer
// for forward compatibility
type AccountingServicesAccountingServer interface {
	GetAccountingByWalletById(context.Context, *RequestGetAccountingByWalletId) (*ResponseGetAccountingByWalletId, error)
	CreateAccounting(context.Context, *RequestCreateAccounting) (*ResponseCreateAccounting, error)
	SetAmountToAccounting(context.Context, *RequestSetAmountToAccounting) (*ResponseSetAmountToAccounting, error)
}

// UnimplementedAccountingServicesAccountingServer should be embedded to have forward compatible implementations.
type UnimplementedAccountingServicesAccountingServer struct {
}

func (UnimplementedAccountingServicesAccountingServer) GetAccountingByWalletById(context.Context, *RequestGetAccountingByWalletId) (*ResponseGetAccountingByWalletId, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccountingByWalletById not implemented")
}
func (UnimplementedAccountingServicesAccountingServer) CreateAccounting(context.Context, *RequestCreateAccounting) (*ResponseCreateAccounting, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccounting not implemented")
}
func (UnimplementedAccountingServicesAccountingServer) SetAmountToAccounting(context.Context, *RequestSetAmountToAccounting) (*ResponseSetAmountToAccounting, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAmountToAccounting not implemented")
}

// UnsafeAccountingServicesAccountingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountingServicesAccountingServer will
// result in compilation errors.
type UnsafeAccountingServicesAccountingServer interface {
	mustEmbedUnimplementedAccountingServicesAccountingServer()
}

func RegisterAccountingServicesAccountingServer(s grpc.ServiceRegistrar, srv AccountingServicesAccountingServer) {
	s.RegisterService(&AccountingServicesAccounting_ServiceDesc, srv)
}

func _AccountingServicesAccounting_GetAccountingByWalletById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestGetAccountingByWalletId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountingServicesAccountingServer).GetAccountingByWalletById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accounting_proto.accountingServicesAccounting/GetAccountingByWalletById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountingServicesAccountingServer).GetAccountingByWalletById(ctx, req.(*RequestGetAccountingByWalletId))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountingServicesAccounting_CreateAccounting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestCreateAccounting)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountingServicesAccountingServer).CreateAccounting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accounting_proto.accountingServicesAccounting/CreateAccounting",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountingServicesAccountingServer).CreateAccounting(ctx, req.(*RequestCreateAccounting))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountingServicesAccounting_SetAmountToAccounting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestSetAmountToAccounting)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountingServicesAccountingServer).SetAmountToAccounting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accounting_proto.accountingServicesAccounting/SetAmountToAccounting",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountingServicesAccountingServer).SetAmountToAccounting(ctx, req.(*RequestSetAmountToAccounting))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountingServicesAccounting_ServiceDesc is the grpc.ServiceDesc for AccountingServicesAccounting service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountingServicesAccounting_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "accounting_proto.accountingServicesAccounting",
	HandlerType: (*AccountingServicesAccountingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAccountingByWalletById",
			Handler:    _AccountingServicesAccounting_GetAccountingByWalletById_Handler,
		},
		{
			MethodName: "CreateAccounting",
			Handler:    _AccountingServicesAccounting_CreateAccounting_Handler,
		},
		{
			MethodName: "SetAmountToAccounting",
			Handler:    _AccountingServicesAccounting_SetAmountToAccounting_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/grpc/proto/accounting.proto",
}
