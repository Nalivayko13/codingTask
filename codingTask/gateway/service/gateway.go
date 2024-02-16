package service

import (
	"context"
	"encoding/json"
	"github.com/Nalivayko13/codingTask/gateway/model"
	"github.com/streadway/amqp"
	"time"
)

func (s *Service) CreateStore(ctx context.Context, store *model.Store) error {
	store.CreatedAt = time.Now().String()
	err := s.send(newEvent(store, "create_store", "", "", ""))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateVersion(ctx context.Context, version *model.Version) error {
	version.CreatedAt = time.Now().String()
	err := s.send(newEvent(version, "create_version", version.Creator, "", ""))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteStore(ctx context.Context, storeID string) error {
	err := s.send(newEvent(nil, "delete_store", "", storeID, ""))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteVersion(ctx context.Context, creator string, storeID, versionID string) error {
	err := s.send(newEvent(nil, "delete_version", creator, string(storeID), string(versionID)))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetStore(ctx context.Context, storeID int) (*model.Store, error) {
	return s.store.GetStoreByID(ctx, storeID)
}

func (s *Service) GetHistory(ctx context.Context, storeID int) ([]model.Version, *model.Store, error) {
	return s.store.GetHistoryByStoreID(ctx, storeID)
}

func (s *Service) GetVersion(ctx context.Context, storeID, versionID int) (*model.Version, *model.Store, error) {
	return s.store.GetVersionByID(ctx, storeID, versionID)
}

func (s *Service) send(message []byte) error {
	err := s.Channel.Publish(
		s.Exchange,
		s.Queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func newEvent(data interface{}, action, creator, storeID, versionID string) []byte {
	message := map[string]interface{}{
		"store_id":   storeID,
		"version_id": versionID,
		"data":       data,
		"action":     action,
		"creator":    creator,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return nil
	}

	return body
}
