// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: api/grpc/proto/users.proto

package users_proto

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

// AuthServicesUsersClient is the client API for AuthServicesUsers service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServicesUsersClient interface {
	CreateUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*ResponseCreateUser, error)
	ActivateUser(ctx context.Context, in *ActivateUserRequest, opts ...grpc.CallOption) (*ValidateResponse, error)
	ValidateEmail(ctx context.Context, in *ValidateEmailRequest, opts ...grpc.CallOption) (*ValidateResponse, error)
	ValidateNickname(ctx context.Context, in *ValidateNicknameRequest, opts ...grpc.CallOption) (*ValidateResponse, error)
	GetUserById(ctx context.Context, in *GetUserByIDRequest, opts ...grpc.CallOption) (*UserResponse, error)
	ValidIdentityNumber(ctx context.Context, in *RequestGetByIdentityNumber, opts ...grpc.CallOption) (*ResponseGetByIdentityNumber, error)
	UpdateUserPhoto(ctx context.Context, in *RequestUpdateUserPhoto, opts ...grpc.CallOption) (*ResponseUpdateUserPhoto, error)
	GetUserPhoto(ctx context.Context, in *RequestGetUserPhoto, opts ...grpc.CallOption) (*ResponseGetUserPhoto, error)
	ChangePassword(ctx context.Context, in *RequestChangePwd, opts ...grpc.CallOption) (*ResponseChangePwd, error)
	CreateUserBySystem(ctx context.Context, in *RequestCreateUserBySystem, opts ...grpc.CallOption) (*ResponseCreateUserBySystem, error)
	CreateUserWallet(ctx context.Context, in *RqCreateUserWallet, opts ...grpc.CallOption) (*ResponseCreateUserWallet, error)
	GetUserWalletByIdentityNumber(ctx context.Context, in *RqGetUserWalletByIdentityNumber, opts ...grpc.CallOption) (*ResGetUserWalletByIdentityNumber, error)
	GetUserByIdentityNumber(ctx context.Context, in *RqGetUserByIdentityNumber, opts ...grpc.CallOption) (*ResGetUserByIdentityNumber, error)
	UpdateUser(ctx context.Context, in *RqUpdateUser, opts ...grpc.CallOption) (*ResUpdateUser, error)
	RequestChangePassword(ctx context.Context, in *RqChangePwd, opts ...grpc.CallOption) (*ResAnny, error)
	UpdateUserIdentity(ctx context.Context, in *RqUpdateUserIdentity, opts ...grpc.CallOption) (*ResUpdateUser, error)
}

type authServicesUsersClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServicesUsersClient(cc grpc.ClientConnInterface) AuthServicesUsersClient {
	return &authServicesUsersClient{cc}
}

func (c *authServicesUsersClient) CreateUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*ResponseCreateUser, error) {
	out := new(ResponseCreateUser)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) ActivateUser(ctx context.Context, in *ActivateUserRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	out := new(ValidateResponse)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/ActivateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) ValidateEmail(ctx context.Context, in *ValidateEmailRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	out := new(ValidateResponse)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/ValidateEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) ValidateNickname(ctx context.Context, in *ValidateNicknameRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	out := new(ValidateResponse)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/ValidateNickname", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) GetUserById(ctx context.Context, in *GetUserByIDRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/GetUserById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) ValidIdentityNumber(ctx context.Context, in *RequestGetByIdentityNumber, opts ...grpc.CallOption) (*ResponseGetByIdentityNumber, error) {
	out := new(ResponseGetByIdentityNumber)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/ValidIdentityNumber", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) UpdateUserPhoto(ctx context.Context, in *RequestUpdateUserPhoto, opts ...grpc.CallOption) (*ResponseUpdateUserPhoto, error) {
	out := new(ResponseUpdateUserPhoto)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/UpdateUserPhoto", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) GetUserPhoto(ctx context.Context, in *RequestGetUserPhoto, opts ...grpc.CallOption) (*ResponseGetUserPhoto, error) {
	out := new(ResponseGetUserPhoto)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/GetUserPhoto", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) ChangePassword(ctx context.Context, in *RequestChangePwd, opts ...grpc.CallOption) (*ResponseChangePwd, error) {
	out := new(ResponseChangePwd)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/ChangePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) CreateUserBySystem(ctx context.Context, in *RequestCreateUserBySystem, opts ...grpc.CallOption) (*ResponseCreateUserBySystem, error) {
	out := new(ResponseCreateUserBySystem)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/CreateUserBySystem", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) CreateUserWallet(ctx context.Context, in *RqCreateUserWallet, opts ...grpc.CallOption) (*ResponseCreateUserWallet, error) {
	out := new(ResponseCreateUserWallet)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/CreateUserWallet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) GetUserWalletByIdentityNumber(ctx context.Context, in *RqGetUserWalletByIdentityNumber, opts ...grpc.CallOption) (*ResGetUserWalletByIdentityNumber, error) {
	out := new(ResGetUserWalletByIdentityNumber)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/GetUserWalletByIdentityNumber", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) GetUserByIdentityNumber(ctx context.Context, in *RqGetUserByIdentityNumber, opts ...grpc.CallOption) (*ResGetUserByIdentityNumber, error) {
	out := new(ResGetUserByIdentityNumber)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/GetUserByIdentityNumber", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) UpdateUser(ctx context.Context, in *RqUpdateUser, opts ...grpc.CallOption) (*ResUpdateUser, error) {
	out := new(ResUpdateUser)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/UpdateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) RequestChangePassword(ctx context.Context, in *RqChangePwd, opts ...grpc.CallOption) (*ResAnny, error) {
	out := new(ResAnny)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/RequestChangePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServicesUsersClient) UpdateUserIdentity(ctx context.Context, in *RqUpdateUserIdentity, opts ...grpc.CallOption) (*ResUpdateUser, error) {
	out := new(ResUpdateUser)
	err := c.cc.Invoke(ctx, "/users_proto.authServicesUsers/UpdateUserIdentity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServicesUsersServer is the server API for AuthServicesUsers service.
// All implementations should embed UnimplementedAuthServicesUsersServer
// for forward compatibility
type AuthServicesUsersServer interface {
	CreateUser(context.Context, *UserRequest) (*ResponseCreateUser, error)
	ActivateUser(context.Context, *ActivateUserRequest) (*ValidateResponse, error)
	ValidateEmail(context.Context, *ValidateEmailRequest) (*ValidateResponse, error)
	ValidateNickname(context.Context, *ValidateNicknameRequest) (*ValidateResponse, error)
	GetUserById(context.Context, *GetUserByIDRequest) (*UserResponse, error)
	ValidIdentityNumber(context.Context, *RequestGetByIdentityNumber) (*ResponseGetByIdentityNumber, error)
	UpdateUserPhoto(context.Context, *RequestUpdateUserPhoto) (*ResponseUpdateUserPhoto, error)
	GetUserPhoto(context.Context, *RequestGetUserPhoto) (*ResponseGetUserPhoto, error)
	ChangePassword(context.Context, *RequestChangePwd) (*ResponseChangePwd, error)
	CreateUserBySystem(context.Context, *RequestCreateUserBySystem) (*ResponseCreateUserBySystem, error)
	CreateUserWallet(context.Context, *RqCreateUserWallet) (*ResponseCreateUserWallet, error)
	GetUserWalletByIdentityNumber(context.Context, *RqGetUserWalletByIdentityNumber) (*ResGetUserWalletByIdentityNumber, error)
	GetUserByIdentityNumber(context.Context, *RqGetUserByIdentityNumber) (*ResGetUserByIdentityNumber, error)
	UpdateUser(context.Context, *RqUpdateUser) (*ResUpdateUser, error)
	RequestChangePassword(context.Context, *RqChangePwd) (*ResAnny, error)
	UpdateUserIdentity(context.Context, *RqUpdateUserIdentity) (*ResUpdateUser, error)
}

// UnimplementedAuthServicesUsersServer should be embedded to have forward compatible implementations.
type UnimplementedAuthServicesUsersServer struct {
}

func (UnimplementedAuthServicesUsersServer) CreateUser(context.Context, *UserRequest) (*ResponseCreateUser, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedAuthServicesUsersServer) ActivateUser(context.Context, *ActivateUserRequest) (*ValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivateUser not implemented")
}
func (UnimplementedAuthServicesUsersServer) ValidateEmail(context.Context, *ValidateEmailRequest) (*ValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateEmail not implemented")
}
func (UnimplementedAuthServicesUsersServer) ValidateNickname(context.Context, *ValidateNicknameRequest) (*ValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateNickname not implemented")
}
func (UnimplementedAuthServicesUsersServer) GetUserById(context.Context, *GetUserByIDRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserById not implemented")
}
func (UnimplementedAuthServicesUsersServer) ValidIdentityNumber(context.Context, *RequestGetByIdentityNumber) (*ResponseGetByIdentityNumber, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidIdentityNumber not implemented")
}
func (UnimplementedAuthServicesUsersServer) UpdateUserPhoto(context.Context, *RequestUpdateUserPhoto) (*ResponseUpdateUserPhoto, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserPhoto not implemented")
}
func (UnimplementedAuthServicesUsersServer) GetUserPhoto(context.Context, *RequestGetUserPhoto) (*ResponseGetUserPhoto, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserPhoto not implemented")
}
func (UnimplementedAuthServicesUsersServer) ChangePassword(context.Context, *RequestChangePwd) (*ResponseChangePwd, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (UnimplementedAuthServicesUsersServer) CreateUserBySystem(context.Context, *RequestCreateUserBySystem) (*ResponseCreateUserBySystem, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserBySystem not implemented")
}
func (UnimplementedAuthServicesUsersServer) CreateUserWallet(context.Context, *RqCreateUserWallet) (*ResponseCreateUserWallet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserWallet not implemented")
}
func (UnimplementedAuthServicesUsersServer) GetUserWalletByIdentityNumber(context.Context, *RqGetUserWalletByIdentityNumber) (*ResGetUserWalletByIdentityNumber, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserWalletByIdentityNumber not implemented")
}
func (UnimplementedAuthServicesUsersServer) GetUserByIdentityNumber(context.Context, *RqGetUserByIdentityNumber) (*ResGetUserByIdentityNumber, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByIdentityNumber not implemented")
}
func (UnimplementedAuthServicesUsersServer) UpdateUser(context.Context, *RqUpdateUser) (*ResUpdateUser, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedAuthServicesUsersServer) RequestChangePassword(context.Context, *RqChangePwd) (*ResAnny, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestChangePassword not implemented")
}
func (UnimplementedAuthServicesUsersServer) UpdateUserIdentity(context.Context, *RqUpdateUserIdentity) (*ResUpdateUser, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserIdentity not implemented")
}

// UnsafeAuthServicesUsersServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServicesUsersServer will
// result in compilation errors.
type UnsafeAuthServicesUsersServer interface {
	mustEmbedUnimplementedAuthServicesUsersServer()
}

func RegisterAuthServicesUsersServer(s grpc.ServiceRegistrar, srv AuthServicesUsersServer) {
	s.RegisterService(&AuthServicesUsers_ServiceDesc, srv)
}

func _AuthServicesUsers_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).CreateUser(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_ActivateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).ActivateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/ActivateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).ActivateUser(ctx, req.(*ActivateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_ValidateEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).ValidateEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/ValidateEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).ValidateEmail(ctx, req.(*ValidateEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_ValidateNickname_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateNicknameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).ValidateNickname(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/ValidateNickname",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).ValidateNickname(ctx, req.(*ValidateNicknameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_GetUserById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).GetUserById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/GetUserById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).GetUserById(ctx, req.(*GetUserByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_ValidIdentityNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestGetByIdentityNumber)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).ValidIdentityNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/ValidIdentityNumber",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).ValidIdentityNumber(ctx, req.(*RequestGetByIdentityNumber))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_UpdateUserPhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestUpdateUserPhoto)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).UpdateUserPhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/UpdateUserPhoto",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).UpdateUserPhoto(ctx, req.(*RequestUpdateUserPhoto))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_GetUserPhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestGetUserPhoto)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).GetUserPhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/GetUserPhoto",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).GetUserPhoto(ctx, req.(*RequestGetUserPhoto))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_ChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestChangePwd)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).ChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/ChangePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).ChangePassword(ctx, req.(*RequestChangePwd))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_CreateUserBySystem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestCreateUserBySystem)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).CreateUserBySystem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/CreateUserBySystem",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).CreateUserBySystem(ctx, req.(*RequestCreateUserBySystem))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_CreateUserWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RqCreateUserWallet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).CreateUserWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/CreateUserWallet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).CreateUserWallet(ctx, req.(*RqCreateUserWallet))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_GetUserWalletByIdentityNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RqGetUserWalletByIdentityNumber)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).GetUserWalletByIdentityNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/GetUserWalletByIdentityNumber",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).GetUserWalletByIdentityNumber(ctx, req.(*RqGetUserWalletByIdentityNumber))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_GetUserByIdentityNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RqGetUserByIdentityNumber)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).GetUserByIdentityNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/GetUserByIdentityNumber",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).GetUserByIdentityNumber(ctx, req.(*RqGetUserByIdentityNumber))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RqUpdateUser)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/UpdateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).UpdateUser(ctx, req.(*RqUpdateUser))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_RequestChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RqChangePwd)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).RequestChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/RequestChangePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).RequestChangePassword(ctx, req.(*RqChangePwd))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthServicesUsers_UpdateUserIdentity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RqUpdateUserIdentity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServicesUsersServer).UpdateUserIdentity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users_proto.authServicesUsers/UpdateUserIdentity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServicesUsersServer).UpdateUserIdentity(ctx, req.(*RqUpdateUserIdentity))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthServicesUsers_ServiceDesc is the grpc.ServiceDesc for AuthServicesUsers service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthServicesUsers_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "users_proto.authServicesUsers",
	HandlerType: (*AuthServicesUsersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _AuthServicesUsers_CreateUser_Handler,
		},
		{
			MethodName: "ActivateUser",
			Handler:    _AuthServicesUsers_ActivateUser_Handler,
		},
		{
			MethodName: "ValidateEmail",
			Handler:    _AuthServicesUsers_ValidateEmail_Handler,
		},
		{
			MethodName: "ValidateNickname",
			Handler:    _AuthServicesUsers_ValidateNickname_Handler,
		},
		{
			MethodName: "GetUserById",
			Handler:    _AuthServicesUsers_GetUserById_Handler,
		},
		{
			MethodName: "ValidIdentityNumber",
			Handler:    _AuthServicesUsers_ValidIdentityNumber_Handler,
		},
		{
			MethodName: "UpdateUserPhoto",
			Handler:    _AuthServicesUsers_UpdateUserPhoto_Handler,
		},
		{
			MethodName: "GetUserPhoto",
			Handler:    _AuthServicesUsers_GetUserPhoto_Handler,
		},
		{
			MethodName: "ChangePassword",
			Handler:    _AuthServicesUsers_ChangePassword_Handler,
		},
		{
			MethodName: "CreateUserBySystem",
			Handler:    _AuthServicesUsers_CreateUserBySystem_Handler,
		},
		{
			MethodName: "CreateUserWallet",
			Handler:    _AuthServicesUsers_CreateUserWallet_Handler,
		},
		{
			MethodName: "GetUserWalletByIdentityNumber",
			Handler:    _AuthServicesUsers_GetUserWalletByIdentityNumber_Handler,
		},
		{
			MethodName: "GetUserByIdentityNumber",
			Handler:    _AuthServicesUsers_GetUserByIdentityNumber_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _AuthServicesUsers_UpdateUser_Handler,
		},
		{
			MethodName: "RequestChangePassword",
			Handler:    _AuthServicesUsers_RequestChangePassword_Handler,
		},
		{
			MethodName: "UpdateUserIdentity",
			Handler:    _AuthServicesUsers_UpdateUserIdentity_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/grpc/proto/users.proto",
}
