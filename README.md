# codingTask

How to start?

0. Create ```.env``` file for gateway and store services
 #### example of ```.env``` for gateway:

```
#Server
HOST=localhost
API_SERVER_HOST=:8081
LOG_LEVEL=debug
GRPC_PORT=:8090
GRPC_HOST=store:8090

#RabbitMQ
RABBITMQ_PASSWORD=guest
RABBITMQ_USERNAME=guest
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_EXCHANGE_NAME=store
RABBITMQ_QUEUE_NAME=gateway

#Auth
URL_GENERATION_TOKEN=http://auth:8080/generate
URL_GENERATION_TOKEN_VALIDATE=http://auth:8080/validate

```

 #### example of ```.env``` for gateway:

```
#Postgres
DB_PASSWORD=mysecretpassword
DB_USERNAME=postgres
DB_SSLMODE=disable
DB_PORT=5432
DB_DBNAME=postgres
DB_HOST=postgres
DB_TESTNAME=test

#Server
HOST=localhost
gRPC_HOST=localhost:8090
API_SERVER_HOST=:8082
LOG_LEVEL=debug
GRPC_PORT=:8090

#Rabbit
RABBITMQ_PASSWORD=guest
RABBITMQ_USERNAME=guest
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_VHOST=/
RABBITMQ_EXCHANGE_NAME=store
RABBITMQ_QUEUE_NAME=gateway
```

1. docker compose up --build
2. Test the work and see documentation using swagger by endpoint http://localhost:8081/docs/index.html
