package main

import (
	"github.com/Nalivayko13/codingTask/store/logging"
	"github.com/Nalivayko13/codingTask/store/repository"
	"github.com/Nalivayko13/codingTask/store/service"
	"github.com/Nalivayko13/codingTask/store/transport"
	"github.com/Nalivayko13/codingTask/store/transport/rabbitmq"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

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

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_DBNAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatal("failed to create db: ", err)
	}
	defer db.Close()

	repos := repository.NewPostgresRepo(db)

	if err := repos.Migrate(db); err != nil {
		log.Fatal("failed to start migrations: ", err.Error())
	}

	rabbitConn, err := rabbitmq.NewRabbitMQConnection(&rabbitmq.RabbitConfig{
		User:   os.Getenv("RABBITMQ_USERNAME"),
		Passwd: os.Getenv("RABBITMQ_PASSWORD"),
		Host:   os.Getenv("RABBITMQ_HOST"),
		Port:   os.Getenv("RABBITMQ_PORT"),
	})
	if err != nil {
		log.Fatal("failed to create rabbitmq conn: ", err.Error())
	}
	rabbitChan, err := rabbitmq.NewRabbitChannel(rabbitConn)
	if err != nil {
		log.Fatal("failed to create rabbitmq chann: ", err.Error())
	}
	rabbitQueue, err := rabbitmq.NewRabbitQueue(rabbitChan, os.Getenv("RABBITMQ_QUEUE_NAME"),
		os.Getenv("RABBITMQ_EXCHANGE_NAME"))
	if err != nil {
		log.Fatal("failed to create rabbitmq queue: ", err.Error())
	}

	services := service.NewService(repos)

	eventHandler := rabbitmq.NewEventHandler(services, logger)

	s, lis, err := transport.NewGRPCServer(services, *logger, os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatal("failed to create grpc server")
	}

	logger.Log.Info("grpc server is started")
	go func() {
		err = transport.RunGRPCServer(s, lis)
		if err != nil {
			log.Printf("error occured while running grpc server: %s", err.Error())
		}
	}()

	eventHandler.Consume(rabbitChan, rabbitQueue)
}
