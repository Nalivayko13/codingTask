package service

import (
	"context"
	"github.com/Nalivayko13/codingTask/store/model"
)

func (s *Service) CreateStore(ctx context.Context, store *model.Store) error {
	err := s.Repo.CreateStore(ctx, store)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateVersion(ctx context.Context, version *model.Version) error {
	err := s.Repo.CreateVersion(ctx, version)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetStore(ctx context.Context, id int) (*model.Store, error) {
	storeDB, err := s.Repo.GetStore(ctx, id)
	if err != nil {
		return nil, err
	}
	return storeDB, err
}

func (s *Service) GetHistory(ctx context.Context, storeID int) ([]model.Version, error) {
	history, err := s.Repo.GetVersions(ctx, storeID)
	if err != nil {
		return nil, err
	}
	return history, err
}

func (s *Service) GetVersion(ctx context.Context, storeID, versionID int) (*model.Version, error) {
	version, err := s.Repo.GetVersion(ctx, storeID, versionID)
	if err != nil {
		return nil, err
	}
	return version, err
}
