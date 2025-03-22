package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/tinode/chat/server/auth"
	"github.com/tinode/chat/server/config"
	adapter "github.com/tinode/chat/server/db"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbx "github.com/tinode/chat/pbx"
)

// AdminServer is a gRPC server implementation of the admin service.
type AdminServer struct {
	pbx.UnimplementedAdminServiceServer
	logger *slog.Logger

	config config.Config
	db adapter.Adapter
}

// NewAdminServer creates a new instance of AdminServer.
func NewAdminServer(logger *slog.Logger, config config.Config, db adapter.Adapter) *AdminServer {
	return &AdminServer{
		config: config,
		db: db,
		logger: logger,
	}
}

// MakeUserRoot changes a user's authentication level to ROOT.
func (a *AdminServer) MakeUserRoot(ctx context.Context, req *pbx.MakeUserRootReq) (*pbx.MakeUserRootResp, error) {
	logger := a.logger.With(slog.String("user", req.UserId))
	
	// Parse user ID
	userId := types.ParseUserId(req.UserId)
	if userId.IsZero() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID '%s'", req.UserId)
	}

	// Update user auth level to ROOT
	if err := a.db.AuthUpdRecord(userId, "basic", "", auth.LevelRoot, nil, time.Time{}); err != nil {
		logger.Warn("Failed to promote user to ROOT: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to promote user: %v", err)
	}
	
	logger.Info("User '%s' successfully promoted to ROOT", req.UserId)
	return &pbx.MakeUserRootResp{}, nil
}

// GetUser retrieves detailed information about a user.
func (a *AdminServer) GetUser(ctx context.Context, req *pbx.GetUserReq) (*pbx.GetUserResp, error) {
	logs.Info.Printf("GetUser request for user %s", req.UserId)

	// TODO: Implement getting user information

	return &pbx.GetUserResp{}, nil
}

// UpdateUserState changes a user's account state (normal, suspended).
func (a *AdminServer) UpdateUserState(ctx context.Context, req *pbx.UpdateUserStateReq) (*pbx.UpdateUserStateResp, error) {
	logs.Info.Printf("UpdateUserState request for user %s to state %v", req.UserId, req.State)

	// TODO: Implement user state update
	// Check the state value
	switch req.State {
	case pbx.AccountState_ACCOUNT_STATE_NORMAL, pbx.AccountState_ACCOUNT_STATE_SUSPENDED:
		// Valid states, continue
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid account state")
	}

	// TODO: Update user state in database

	return &pbx.UpdateUserStateResp{}, nil
}

// DeleteUser permanently removes a user account.
func (a *AdminServer) DeleteUser(ctx context.Context, req *pbx.DeleteUserReq) (*pbx.DeleteUserResp, error) {
	logs.Info.Printf("DeleteUser request for user %s (hard=%v)", req.UserId, req.Hard)

	// TODO: Implement user deletion
	// Check if user exists
	// Perform either soft or hard delete based on req.Hard flag

	return &pbx.DeleteUserResp{}, nil
}
