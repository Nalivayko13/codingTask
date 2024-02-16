package service

import (
	"context"
	"github.com/Nalivayko13/codingTask/store/model"
)

type Service struct {
	Repo
}

func NewService(repo Repo) *Service {
	return &Service{Repo: repo}
}

type Repo interface {
	CreateStore(ctx context.Context, store *model.Store) error
	CreateVersion(ctx context.Context, version *model.Version) error
	DeleteStore(ctx context.Context, storeID int) error
	DeleteVersion(ctx context.Context, versionID, storeID int) error
	GetStore(ctx context.Context, id int) (*model.Store, error)
	GetVersions(ctx context.Context, storeID int) ([]model.Version, error)
	GetVersion(ctx context.Context, storeID, versionID int) (*model.Version, error)
}
