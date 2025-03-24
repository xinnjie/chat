package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/tinode/chat/server/auth"
	"github.com/tinode/chat/server/config"
	adapter "github.com/tinode/chat/server/db"
	"github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbx "github.com/tinode/chat/pbx"
)

// AdminServer is a gRPC server implementation of the admin service.
type AdminServer struct {
	pbx.UnimplementedAdminServiceServer
	logger *slog.Logger

	config  config.Config
	db      adapter.Adapter
	cluster *Cluster
}

// NewAdminServer creates a new instance of AdminServer.
func NewAdminServer(logger *slog.Logger, config config.Config, db adapter.Adapter, cluster *Cluster) *AdminServer {
	return &AdminServer{
		config:  config,
		db:      db,
		logger:  logger,
		cluster: cluster,
	}
}

// MakeUserRoot changes a user's authentication level to ROOT.
func (a *AdminServer) MakeUserRoot(ctx context.Context, req *pbx.MakeUserRootRequest) (*pbx.MakeUserRootResponse, error) {
	logger := a.logger.With(slog.String("user", req.UserId))

	// Parse user ID
	userId := types.ParseUserId(req.UserId)
	if userId.IsZero() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID '%s'", req.UserId)
	}

	// Update user auth level to ROOT
	if err := a.db.AuthUpdRecord(userId, "basic", "", auth.LevelRoot, nil, time.Time{}); err != nil {
		logger.Warn("Failed to promote user to ROOT", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to promote user: %v", err)
	}

	logger.Info("User successfully promoted to ROOT")
	return &pbx.MakeUserRootResponse{}, nil
}

// GetUser retrieves detailed information about a user.
func (a *AdminServer) GetUser(ctx context.Context, req *pbx.GetUserRequest) (*pbx.User, error) {
	// Parse user ID
	userId := types.ParseUserId(req.UserId)
	if userId.IsZero() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID '%s'", req.UserId)
	}

	// Get user data from database
	user, err := a.db.UserGet(userId)
	if err != nil {
		a.logger.Warn("Failed to get user", "error", err, "userId", userId.String())
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %s", req.UserId)
	}

	// Create response
	resp := &pbx.User{
		UserId: userId.UserId(),
	}

	// Handle Public field - convert to string if it's not nil
	if user.Public != nil {
		// Check if it's already a string
		if publicStr, ok := user.Public.(string); ok {
			resp.Public = publicStr
		} else {
			// Try to marshal it to JSON
			publicBytes, err := json.Marshal(user.Public)
			if err == nil {
				resp.Public = string(publicBytes)
			} else {
				a.logger.Warn("Failed to marshal user public data", "error", err)
			}
		}
	}

	// Add lastSeen timestamp if available
	if user.LastSeen != nil && !user.LastSeen.IsZero() {
		resp.LastSeen = timestamppb.New(*user.LastSeen)
	}

	return resp, nil
}

// UpdateUserState changes a user's account state (normal, suspended).
func (a *AdminServer) UpdateUserState(ctx context.Context, req *pbx.UpdateUserStateRequest) (*pbx.UpdateUserStateResponse, error) {
	// Parse user ID
	userId := types.ParseUserId(req.UserId)
	if userId.IsZero() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID '%s'", req.UserId)
	}

	// TODO: Implement user state update
	// Check the state value
	switch req.State {
	case pbx.AccountState_ACCOUNT_STATE_NORMAL, pbx.AccountState_ACCOUNT_STATE_SUSPENDED:
		// Valid states, continue
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid account state")
	}

	// TODO: Update user state in database

	return &pbx.UpdateUserStateResponse{}, nil
}

// DeleteUser permanently removes a user account.
func (a *AdminServer) DeleteUser(ctx context.Context, req *pbx.DeleteUserRequest) (*pbx.DeleteUserResponse, error) {
	// Parse user ID
	userId := types.ParseUserId(req.UserId)
	if userId.IsZero() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID '%s'", req.UserId)
	}

	// TODO: Implement user deletion
	// Check if user exists
	// Perform either soft or hard delete based on req.Hard flag

	return &pbx.DeleteUserResponse{}, nil
}

// CreateGroupChat creates a new group chat.
func (a *AdminServer) CreateGroupChat(ctx context.Context, req *pbx.CreateGroupChatRequest) (*pbx.CreateGroupChatResponse, error) {
	// Parse user ID
	userId := types.ParseUserId(req.UserId)
	if userId.IsZero() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID '%s'", req.UserId)
	}
	logger := a.logger.With(slog.String("user", req.UserId))

	// Check if the user exists
	user, err := a.db.UserGet(userId)
	if err != nil {
		logger.Warn("Failed to get user", "error", err, "userId", userId.String())
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %s", req.UserId)
	}

	// Generate a unique topic ID for the group chat
	// Group topics in Tinode use "grp" prefix followed by a random string
	topicId := a.cluster.GenerateLocalTopicName()
	now := types.TimeNow()

	// Create the topic
	topic := &types.Topic{
		ObjHeader: types.ObjHeader{Id: topicId, CreatedAt: now},
		Access: types.DefaultAccess{
			Auth: types.ModeCPublic,
			Anon: types.ModeNone,
		},
		Public: map[string]interface{}{
			"fn": req.ChatName,
		},
	}

	// Create the topic with the user as owner
	if err := store.Topics.Create(topic, userId, nil); err != nil {
		logger.Warn("Failed to create group chat", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create group chat: %v", err)
	}

	logger.Info("Group chat created", "topicId", topicId)

	return &pbx.CreateGroupChatResponse{
		TopicId: topicId,
	}, nil
}
