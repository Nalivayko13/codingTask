package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/Nalivayko13/codingTask/store/logging"
	"github.com/Nalivayko13/codingTask/store/model"
	"github.com/Nalivayko13/codingTask/store/service"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
	"strconv"
)

type EventAction struct {
	Action    string `json:"action"`
	StoreID   string `json:"store_id"`
	VersionID string `json:"version_id"`
	Creator   string `json:"creator"`
}

type EventStore struct {
	Action    string      `json:"action"`
	Data      model.Store `json:"data"`
	StoreID   string      `json:"store_id"`
	Creator   string      `json:"creator"`
	VersionID string      `json:"version_id"`
}

type EventVersion struct {
	Action    string        `json:"action"`
	Data      model.Version `json:"data"`
	StoreID   string        `json:"store_id"`
	Creator   string        `json:"creator"`
	VersionID string        `json:"version_id"`
}

type EventHandler struct {
	service *service.Service
	log     *logging.Logger
}

func NewEventHandler(serv *service.Service, log *logging.Logger) *EventHandler {
	return &EventHandler{
		service: serv,
		log:     log,
	}
}

func (h *EventHandler) HandleEvent(data amqp.Delivery) {
	eventAction, err := extractEventAction(data)
	if err != nil {
		h.log.Log.Error("failed to extract event", zap.Error(err))
		return
	}
	ctx := context.Background()

	switch eventAction.Action {
	case "create_store":
		h.CreateStore(ctx, data)
	case "create_version":
		h.CreateVersion(ctx, data)
	case "delete_store":
		h.DeleteStore(ctx, eventAction.StoreID)
	case "delete_version":
		h.DeleteVersion(ctx, eventAction.StoreID, eventAction.VersionID, eventAction.Creator)
	default:
		h.log.Log.Warn("action not found", zap.String("action", eventAction.Action))
	}
}

func (h *EventHandler) Consume(rabbitChan *amqp.Channel, rabbitQueue amqp.Queue) {
	msg, err := rabbitChan.Consume(
		rabbitQueue.Name,
		"store",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		h.log.Log.Error("failed to consume mess: ", zap.Error(err))
	}

	var ch chan struct{}

	go func() {
		for d := range msg {
			log.Printf("Received a message: %s", d.Body)
			h.HandleEvent(d)
		}
	}()

	<-ch
}

func (h *EventHandler) CreateStore(ctx context.Context, data amqp.Delivery) {
	store, err := extractStore(data)
	if err != nil {
		h.log.Log.Error("failed to extract data", zap.Error(err))
		return
	}
	err = h.service.CreateStore(ctx, store)
	if err != nil {
		h.log.Log.Error("failed to create store", zap.Error(err))
	}
	h.log.Log.Debug("created store")
}

func (h *EventHandler) CreateVersion(ctx context.Context, data amqp.Delivery) {
	version, err := extractVersion(data)
	if err != nil {
		h.log.Log.Error("Failed to extract data", zap.Error(err))
		return
	}

	err = h.service.CreateVersion(ctx, version)
	if err != nil {
		h.log.Log.Error("Failed to create version", zap.Error(err))
	}
	h.log.Log.Info("version created")
}

func (h *EventHandler) DeleteStore(ctx context.Context, storeIDstr string) {
	storeID, err := strconv.Atoi(storeIDstr)
	if err != nil {
		h.log.Log.Error("failed to convert storeID to int", zap.Error(err))
	}
	err = h.service.DeleteStore(ctx, storeID)

	if err != nil {
		h.log.Log.Error("failed to delete store", zap.Error(err))
	}
	h.log.Log.Info("store deleted")
}

func (h *EventHandler) DeleteVersion(ctx context.Context, storeIDstr, versionIDstr, creator string) {
	storeID, err := strconv.Atoi(storeIDstr)
	if err != nil {
		h.log.Log.Error("failed to convert storeID to int", zap.Error(err))
		return
	}
	versionID, err := strconv.Atoi(versionIDstr)
	if err != nil {
		h.log.Log.Error("failed to convert VersionID to int", zap.Error(err))
		return
	}

	err = h.service.DeleteVersion(ctx, versionID, storeID)
	if err != nil {
		h.log.Log.Error("failed to delete store version", zap.Error(err))
		return
	}
	h.log.Log.Info("Store version deleted successfully")
}

func extractEventAction(data amqp.Delivery) (EventAction, error) {
	var ev EventAction
	err := json.Unmarshal(data.Body, &ev)
	if err != nil {
		return EventAction{}, err
	}
	return ev, nil
}

func extractStore(data amqp.Delivery) (*model.Store, error) {
	var storeData EventStore
	err := json.Unmarshal(data.Body, &storeData)
	if err != nil {
		return nil, err
	}
	return &storeData.Data, nil
}

func extractVersion(data amqp.Delivery) (*model.Version, error) {
	var version EventVersion
	err := json.Unmarshal(data.Body, &version)
	if err != nil {
		return nil, err
	}
	return &version.Data, nil
}
