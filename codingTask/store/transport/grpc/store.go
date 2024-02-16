package grpc

import (
	"context"
	"github.com/Nalivayko13/codingTask/store/logging"
	"github.com/Nalivayko13/codingTask/store/model"
	storeProto "github.com/Nalivayko13/codingTask/store/pkg/store_service"
	"github.com/Nalivayko13/codingTask/store/service"
	"go.uber.org/zap"
)

type GRPCServer struct {
	Service *service.Service
	Logger  logging.Logger
	storeProto.UnimplementedStoreServiceServer
}

func (s *GRPCServer) GetStore(ctx context.Context, storeID *storeProto.StoreID) (*storeProto.Store, error) {
	store, err := s.Service.GetStore(ctx, int(storeID.StoreID))
	if err != nil {
		s.Logger.Log.Error("failed to get store via grpc", zap.Error(err))
		return nil, err
	}
	return convertStoreToProto(store), nil
}

func convertStoreToProto(store *model.Store) *storeProto.Store {
	return &storeProto.Store{
		ID:        int64(store.ID),
		Name:      store.Name,
		Address:   store.Address,
		Creator:   store.Creator,
		CreatedAt: store.CreatedAt,
	}
}

func (s *GRPCServer) GetHistory(ctx context.Context, storeID *storeProto.StoreID) (*storeProto.History, error) {
	history, err := s.Service.GetVersions(ctx, int(storeID.StoreID))
	if err != nil {
		s.Logger.Log.Error("failed to get history via grpc", zap.Error(err))
		return nil, err
	}
	store, err := s.Service.GetStore(ctx, int(storeID.StoreID))
	if err != nil {
		s.Logger.Log.Error("failed to get store via grpc", zap.Error(err))
		return nil, err
	}
	var versions []*storeProto.Version
	for _, v := range history {
		versions = append(versions, convertToProtoVersion(&v))
	}
	return &storeProto.History{
		History: versions,
		Info:    convertStoreToProto(store),
	}, nil
}
func (s *GRPCServer) GetVersion(ctx context.Context, versionID *storeProto.VersionID) (*storeProto.Version, error) {
	version, err := s.Service.GetVersion(ctx, int(versionID.StoreID), int(versionID.VersionID))
	if err != nil {
		s.Logger.Log.Error("failed to get version via grpc", zap.Error(err))
		return nil, err
	}
	store, err := s.Service.GetStore(ctx, int(versionID.StoreID))
	if err != nil {
		s.Logger.Log.Error("failed to get store via grpc", zap.Error(err))
		return nil, err
	}

	return &storeProto.Version{
		ID:            int64(version.ID),
		Info:          convertStoreToProto(store),
		VersionNumber: int64(version.VersionNumber),
		Creator:       version.Creator,
		Owner:         version.Owner,
		OpenAt:        version.OpenAt,
		CloseAt:       version.CloseAt,
		CreatedAt:     version.CreatedAt,
	}, nil
}

func convertToProtoVersion(version *model.Version) *storeProto.Version {
	return &storeProto.Version{
		ID:            int64(version.ID),
		VersionNumber: int64(version.VersionNumber),
		Creator:       version.Creator,
		CreatedAt:     version.CreatedAt,
		Owner:         version.Owner,
		OpenAt:        version.OpenAt,
		CloseAt:       version.CloseAt,
	}
}
