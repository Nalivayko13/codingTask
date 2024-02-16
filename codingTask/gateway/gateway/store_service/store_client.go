package store_service

import (
	"context"
	"fmt"
	"github.com/Nalivayko13/codingTask/gateway/model"
	storeProto "github.com/Nalivayko13/codingTask/gateway/pkg/store_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StoreGRPCClient struct {
	cli storeProto.StoreServiceClient
}

func NewStoreClient(host string) (*StoreGRPCClient, error) {
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC store client: %w", err)
	}
	if conn == nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}
	cli := storeProto.NewStoreServiceClient(conn)
	return &StoreGRPCClient{cli: cli}, nil
}

func (s *StoreGRPCClient) GetStoreByID(ctx context.Context, storeID int) (*model.Store, error) {
	store, err := s.cli.GetStore(ctx, &storeProto.StoreID{StoreID: int64(storeID)})
	if err != nil {
		return nil, err
	}
	return convertProtoStoreToModel(store), nil
}

func (s *StoreGRPCClient) GetHistoryByStoreID(ctx context.Context, storeID int) ([]model.Version, *model.Store, error) {
	history, err := s.cli.GetHistory(ctx, &storeProto.StoreID{StoreID: int64(storeID)})
	if err != nil {
		return nil, nil, err
	}
	var versions []model.Version
	for _, v := range history.History {
		ver := convertProtoVersionToModel(v, int(history.Info.ID))
		versions = append(versions, *ver)
	}
	return versions, convertProtoStoreToModel(history.Info), nil
}

func (s *StoreGRPCClient) GetVersionByID(ctx context.Context, storeID, versionID int) (*model.Version, *model.Store, error) {
	version, err := s.cli.GetVersion(ctx,
		&storeProto.VersionID{
			StoreID:   int64(storeID),
			VersionID: int64(versionID),
		})
	if err != nil {
		return nil, nil, err
	}
	return convertProtoVersionToModel(version, storeID), convertProtoStoreToModel(version.Info), nil
}
