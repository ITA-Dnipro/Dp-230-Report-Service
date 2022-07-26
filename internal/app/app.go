package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/config"
	grpcWrapper "github.com/ITA-Dnipro/Dp-230-Report-Service/internal/grpc"
	kafkaWrapper "github.com/ITA-Dnipro/Dp-230-Report-Service/internal/kafka"
	"github.com/Shopify/sarama"
)

type Component interface {
	Start(context.Context) error
	Stop(context.Context) error
	Unwrap() interface{}
	Name() string
}

type App struct {
	config         config.Config
	logger         *zap.Logger
	initComponents map[string]Component
}

func NewApp(cfg config.Config) (*App, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	components, err := initComponents(cfg, logger)
	if err != nil {
		return nil, err
	}

	app := &App{
		logger:         logger,
		initComponents: components,
	}
	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer cancel()
		signaled := make(chan os.Signal, 1)
		signal.Notify(signaled, os.Interrupt, syscall.SIGTERM)
		select {
		case s := <-signaled:
			a.logger.Info("exiting", zap.Stringer("signal", s))
			signal.Stop(signaled)

		case <-ctx.Done():
			signal.Stop(signaled)
		}
	}()

	for name, comp := range a.initComponents {
		if err := comp.Start(ctx); err != nil {
			a.logger.Warn("failed to stop", zap.String("component", name), zap.Error(err))
		}
	}

	<-ctx.Done()
	if err := a.stop(); err != nil {
		a.logger.Warn("failed to stop app", zap.Error(err))
	}
	return ctx.Err()
}

func (a *App) stop() error {
	shutdownTimeout := time.Second * 30
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	_ = ctx
	// for _, c := range a.components {
	// 	err := c.Stop(ctx)
	// 	if err != nil {
	// 		a.logger.Warn("failed to stop", zap.String("component", c.Name()), zap.Error(err))
	// 	}
	// }
	return nil
}
func (a *App) Unwrap(dep interface{}) error {
	t := reflect.TypeOf(dep)
	v := reflect.ValueOf(dep)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		fn := tf.Name
		ft := tf.Type.Elem().String()
		comp, ok := a.initComponents[fn]
		if !ok {
			return errors.New(fmt.Sprintf("could not find config for component %s", fn))
		}

		switch ft {
		case "grpc.Server":
			uc, ok := comp.Unwrap().(*grpc.Server)
			if !ok {
				return errors.New("wrong type")
			}
			v.Field(i).Set(reflect.ValueOf(uc))
		case "grpc.ClientConn":
			uc, ok := comp.Unwrap().(*grpc.ClientConn)
			if !ok {
				return errors.New("wrong type")
			}
			v.Field(i).Set(reflect.ValueOf(uc))
		case "sarama.SyncProducer":
			uc, ok := comp.Unwrap().(sarama.SyncProducer)
			if !ok {
				return errors.New("wrong type")
			}
			v.Field(i).Set(reflect.ValueOf(uc))
		case "kafka.ConsumerHandler":
			uc, ok := comp.Unwrap().(*grpc.ClientConn)
			if !ok {
				return errors.New("wrong type")
			}
			v.Field(i).Set(reflect.ValueOf(uc))

		default:
			return errors.New(fmt.Sprintf("unsuported type %s", ft))
		}
	}
	return nil
}

func verifyConfigStructOrPtr(cfg interface{}) error {
	cfgVal := reflect.ValueOf(cfg)
	if cfgVal.Kind() == reflect.Ptr {
		if cfgVal.IsNil() {
			return errors.New("application configuration should be non-nil")
		}
		cfgVal = cfgVal.Elem()
	}
	if cfgVal.Kind() != reflect.Struct {
		return errors.New("application configuration should be a struct")
	}
	return nil
}

func initComponent(fieldType string, val interface{}, log *zap.Logger) (Component, error) {
	var comp interface{}
	var err error
	switch fieldType {
	case "grpc.ServerConfiguration":
		cfg := val.(*grpcWrapper.ServerConfiguration)
		comp, err = grpcWrapper.NewServer(cfg, log)
	case "grpc.ClientConfiguration":
		cfg := val.(*grpcWrapper.ClientConfiguration)
		comp, err = grpcWrapper.NewClient(cfg, log)
	case "kafka.ProducerConfiguration":
		cfg := val.(*kafkaWrapper.ProducerConfiguration)
		comp, err = kafkaWrapper.NewSyncProducer(cfg, log)
	case "kafka.ConsumeGroupConfiguration":
		cfg := val.(*kafkaWrapper.ConsumeGroupConfiguration)
		comp, err = kafkaWrapper.NewConsumerGroup(cfg, log)
	default:
		return nil, errors.New("unknown config type")
	}
	if err != nil {
		return nil, err
	}
	c, ok := comp.(Component)
	if !ok {
		return nil, errors.New("cant convert component")
	}
	return c, nil
}

func initComponents(cfg config.Config, log *zap.Logger) (map[string]Component, error) {
	res := make(map[string]Component)
	t := reflect.TypeOf(cfg)
	v := reflect.ValueOf(cfg)
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i).Type.Elem().String()
		fv := v.Field(i).Interface()
		cp, err := initComponent(ft, fv, log)
		if err != nil {
			return nil, err
		}
		res[t.Field(i).Name] = cp
	}
	return res, nil
}
