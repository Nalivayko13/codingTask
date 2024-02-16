package main

import (
	"context"
	"github.com/Nalivayko13/codingTask/gateway/gateway/store_service"
	"github.com/Nalivayko13/codingTask/gateway/logging"
	"github.com/Nalivayko13/codingTask/gateway/service"
	"github.com/Nalivayko13/codingTask/gateway/transport"
	"github.com/Nalivayko13/codingTask/gateway/transport/rest"
	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title           Swagger for Gateway Service
// @description     This is a simple service to create some operations with stores.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("loading env: ", err)
	}
	defaultLogLevel, err := zapcore.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Println("set up default log level")
		defaultLogLevel = zapcore.DebugLevel
	}
	logger := logging.NewLogger(defaultLogLevel)

	service.UrlGenerateToken = os.Getenv("URL_GENERATION_TOKEN")
	rest.UrlAuthServiceValidate = os.Getenv("URL_GENERATION_TOKEN_VALIDATE")

	rabbitConnection, err := service.NewRabbitMQConnection(&service.RabbitMQConfig{
		Host:     os.Getenv("RABBITMQ_HOST"),
		Port:     os.Getenv("RABBITMQ_PORT"),
		Username: os.Getenv("RABBITMQ_USERNAME"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
	})
	if err != nil {
		log.Fatal("failed connection with rabbitmq: ", err.Error())
	}
	rabbitChannel, err := service.NewRabbitChannel(rabbitConnection)
	if err != nil {
		log.Fatal("failed to create rabbitmq channel: ", err.Error())
	}
	rabbitQueue, err := service.NewRabbitExchangeAndQueue(rabbitChannel,
		os.Getenv("RABBITMQ_EXCHANGE_NAME"), os.Getenv("RABBITMQ_QUEUE_NAME"))
	if err != nil {
		log.Fatal("failed to declare rabbitmq queue: ", err.Error())
	}
	storeCli, err := store_service.NewStoreClient(os.Getenv("GRPC_HOST"))
	if err != nil {
		log.Fatal("failed to create grpc store client: ", err.Error())
	}

	services := service.NewService(rabbitChannel, rabbitConnection, rabbitQueue, storeCli, os.Getenv("RABBITMQ_EXCHANGE_NAME"))
	handlers := rest.NewHandler(services, *logger)

	srv := new(transport.Server)
	go func() {
		if err := srv.Run(os.Getenv("API_SERVER_HOST"), handlers.InitRoutes()); err != nil {
			log.Printf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Printf("Gateway Service Started on port %v", os.Getenv("API_SERVER_HOST"))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("shutting down: %w", err)
	}

	log.Println("Server exiting")
}
