package transport

import (
	"github.com/Nalivayko13/codingTask/store/logging"
	storeProto "github.com/Nalivayko13/codingTask/store/pkg/store_service"
	"github.com/Nalivayko13/codingTask/store/service"
	g "github.com/Nalivayko13/codingTask/store/transport/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func NewGRPCServer(service *service.Service, logger logging.Logger, port string) (*grpc.Server, net.Listener, error) {
	s := grpc.NewServer()
	str := &g.GRPCServer{
		Service: service,
		Logger:  logger,
	}
	storeProto.RegisterStoreServiceServer(s, str)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Log.Error("GRPCServer error", zap.Error(err))
		return nil, nil, err
	}
	reflection.Register(s)
	return s, lis, nil
}

func RunGRPCServer(s *grpc.Server, lis net.Listener) error {
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}
