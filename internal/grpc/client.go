package grpc

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ClientConfiguration struct {
	Target string
}

type grpcClientWrapper struct {
	name       string
	cfg        *ClientConfiguration
	clientConn *grpc.ClientConn
	logger     *zap.Logger
}

func NewClient(config interface{}, log *zap.Logger) (*grpcClientWrapper, error) {
	cfg, ok := config.(*ClientConfiguration)
	if !ok || cfg == nil {
		return nil, errors.New("invalid client config")
	}
	conn, err := grpc.Dial(cfg.Target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &grpcClientWrapper{
		cfg:        cfg,
		clientConn: conn,
		logger:     log,
	}, nil
}

func (w *grpcClientWrapper) Start(context.Context) error {
	return nil
}

func (w *grpcClientWrapper) Stop(context.Context) error {
	w.logger.Info("closing grpc client", zap.String("target", w.cfg.Target))
	return w.clientConn.Close()
}

func (w *grpcClientWrapper) Unwrap() interface{} {
	return w.clientConn
}

func (w *grpcClientWrapper) Name() string {
	return w.name
}
