package auth

import (
	"context"
	"log/slog"
	ssov1 "proto/gen/go/sso"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	Validate = validator.New()
	// errors
	InvalidArgument = "Invalid argument"
	InternalError   = "Internal error"
	InvalidAppID    = "Invalid app_id"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterUser(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	validationErr := Validate.Struct(req)
	if validationErr != nil {
		slog.Error("Failed to validate", validationErr)
		return nil, status.Error(codes.InvalidArgument, "Invalid request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, InvalidArgument)
	}

	token, err := s.auth.Login(ctx, req.Email, req.Password, int(req.GetAppId()))
	if err != nil {
		// TODO:

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	validationErr := Validate.Struct(req)
	if validationErr != nil {
		slog.Error("Failed to validate", validationErr)
		return nil, status.Error(codes.InvalidArgument, "Invalid request")
	}

	userID, err := s.auth.RegisterUser(ctx, req.Email, req.Password)
	if err != nil {
		// TODO:

		return nil, status.Error(codes.InvalidArgument, InvalidArgument)
	}

	return &ssov1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, InvalidArgument)
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.UserId)
	if err != nil {
		// TODO:

		return nil, status.Error(codes.InvalidArgument, InvalidArgument)
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
