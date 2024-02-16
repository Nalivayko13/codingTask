package store_service

import (
	"github.com/Nalivayko13/codingTask/gateway/model"
	storeProto "github.com/Nalivayko13/codingTask/gateway/pkg/store_service"
)

func convertProtoStoreToModel(store *storeProto.Store) *model.Store {
	return &model.Store{
		ID:        int(store.ID),
		Name:      store.Name,
		Address:   store.Address,
		CreatedAt: store.CreatedAt,
	}
}

func convertProtoVersionToModel(version *storeProto.Version, storeID int) *model.Version {
	return &model.Version{
		ID:            int(version.ID),
		StoreID:       storeID,
		VersionNumber: int(version.VersionNumber),
		Creator:       version.Creator,
		CreatedAt:     version.CreatedAt,
		OpenAt:        version.OpenAt,
		CloseAt:       version.CloseAt,
		Owner:         version.Owner,
	}
}
