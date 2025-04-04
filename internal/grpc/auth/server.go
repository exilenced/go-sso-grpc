﻿package auth

import (
	"context"
	ssov1 "github.com/exilenced/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
)

type Auth interface {
	Login(ctx context.Context, username string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, username string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}
type serverAPI struct {
	auth Auth
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})

}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}
	token, err := s.auth.Login(ctx, req.GetUsername(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		//TODO: wrap
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}
	userID, err := s.auth.RegisterNewUser(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		//TODO: wrap
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		//TODO: wrap
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetUsername() == "" {
		return status.Error(codes.InvalidArgument, "username required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app id required")
	}
	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetUsername() == "" {
		return status.Error(codes.InvalidArgument, "username required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password required")
	}
	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user id required")
	}
	return nil
}
