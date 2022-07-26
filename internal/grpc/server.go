package grpc

import (
	"context"
	"errors"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ServerConfiguration struct {
	Addr            string
	ShutdownTimeout time.Duration
}

type grpcServerWrapper struct {
	name   string
	cfg    *ServerConfiguration
	server *grpc.Server
	logger *zap.Logger
}

func NewServer(config interface{}, log *zap.Logger) (*grpcServerWrapper, error) {
	cfg, ok := config.(*ServerConfiguration)
	if !ok || cfg == nil {
		return nil, errors.New("invalid server config")
	}
	srv := grpc.NewServer()
	return &grpcServerWrapper{
		cfg:    cfg,
		server: srv,
		logger: log,
	}, nil
}

func (w *grpcServerWrapper) Start(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	w.logger.Info("starting grpc listener", zap.String("addr", w.cfg.Addr))
	lis, err := net.Listen("tcp", w.cfg.Addr)
	if err != nil {
		return err
	}

	go func() {
		w.logger.Info("starting grpc server")
		if err := w.server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			w.logger.Error("grpc server failed", zap.Error(err))
		}
	}()
	return nil
}

func (w *grpcServerWrapper) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, w.cfg.ShutdownTimeout)
	defer cancel()
	grpcStopCh := make(chan struct{})
	go func() {
		w.logger.Info("gracefully shutting down grpc server")
		w.server.GracefulStop()
		close(grpcStopCh)
	}()
	select {
	case <-grpcStopCh:
		w.logger.Info("gracefully shut down grpc server")
	case <-ctx.Done():
		w.logger.Warn("stopping grpc server")
		w.server.Stop()
		<-grpcStopCh
		grpcStopCh = nil
	}
	return nil
}

func (w *grpcServerWrapper) Unwrap() interface{} {
	return w.server
}

func (w *grpcServerWrapper) Name() string {
	return w.name
}
