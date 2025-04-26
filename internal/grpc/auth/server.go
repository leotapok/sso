package auth

import (
	"context"
	ssov1 "github.com/leotapok/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appID int64,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (UserID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

const (
	emptyValue = 0
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {

	if err := validateLogin(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {

	if err := validateRegister(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())

	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.Email == "" {
		return status.Error(codes.InvalidArgument, "login required")
	}
	if req.Password == "" {
		return status.Error(codes.InvalidArgument, "password required")
	}
	if req.AppId == emptyValue {
		return status.Error(codes.InvalidArgument, "appId required")
	}
	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password required")
	}
	return nil

}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id required")
	}
	return nil
}
