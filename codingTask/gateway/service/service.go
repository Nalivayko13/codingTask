package service

import (
	"context"
	"github.com/Nalivayko13/codingTask/gateway/model"
	"github.com/streadway/amqp"
)

type Service struct {
	Channel  *amqp.Channel
	Conn     *amqp.Connection
	Queue    string
	Exchange string
	store    Store
}

func NewService(channel *amqp.Channel, rabbitMQConn *amqp.Connection, rabbitMQQueue string,
	store Store, exchangeName string) *Service {
	return &Service{
		Channel:  channel,
		Conn:     rabbitMQConn,
		Queue:    rabbitMQQueue,
		store:    store,
		Exchange: exchangeName,
	}
}

type GatewayService interface {
	AuthUser(ctx context.Context, user *model.User) (string, error)
	CreateStore(ctx context.Context, store *model.Store) error
	CreateVersion(ctx context.Context, version *model.Version) error
	DeleteStore(ctx context.Context, storeID string) error
	DeleteVersion(ctx context.Context, creator string, storeID, versionID string) error
	GetStore(ctx context.Context, storeID int) (*model.Store, error)
	GetHistory(ctx context.Context, storeID int) ([]model.Version, *model.Store, error)
	GetVersion(ctx context.Context, storeID, versionID int) (*model.Version, *model.Store, error)
}

type Store interface {
	GetStoreByID(ctx context.Context, storeID int) (*model.Store, error)
	GetHistoryByStoreID(ctx context.Context, storeID int) ([]model.Version, *model.Store, error)
	GetVersionByID(ctx context.Context, storeID, versionID int) (*model.Version, *model.Store, error)
}
