package kafka

import (
	"context"
	"errors"
	"time"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type SyncProducer sarama.SyncProducer

type ProducerConfiguration struct {
	Addrs           []string `default:":9091"`
	ShutdownTimeout time.Duration
}

type syncProducerWrapper struct {
	name         string
	cfg          *ProducerConfiguration
	syncProducer sarama.SyncProducer
	logger       *zap.Logger
}

func NewSyncProducer(config interface{}, log *zap.Logger) (*syncProducerWrapper, error) {
	cfg, ok := config.(*ProducerConfiguration)
	if !ok || cfg == nil {
		return nil, errors.New("invalid sync producer config")
	}
	saramaConfig := sarama.NewConfig()
	// Required by SyncProducer
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	syncProducer, err := sarama.NewSyncProducer(cfg.Addrs, saramaConfig)
	if err != nil {
		return nil, err
	}
	return &syncProducerWrapper{
		cfg:          cfg,
		syncProducer: syncProducer,
		logger:       log,
	}, nil
}

func (w *syncProducerWrapper) Unwrap() interface{} {
	return w.syncProducer
}

func (w *syncProducerWrapper) Start(context.Context) error {
	return nil
}

func (w *syncProducerWrapper) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, w.cfg.ShutdownTimeout)
	defer cancel()

	stop := make(chan error, 1)
	go func() {
		w.logger.Info("closing kafka sync producer")
		stop <- w.syncProducer.Close()
	}()

	select {
	case e := <-stop:
		if e != nil {
			w.logger.Error("could not close kafka sync producer", zap.Error(e))
			return e
		}
		w.logger.Info("gracefully closed kafka sync producer")
	case <-ctx.Done():
		w.logger.Warn("ShutdownTimeout exceeded while closing kafka sync producer")
		return ctx.Err()
	}
	return nil
}

func (w *syncProducerWrapper) Name() string {
	return w.name
}
