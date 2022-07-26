package main

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/app"
	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/config"
)

type AppDep struct {
	Server   *grpc.Server
	Client   *grpc.ClientConn
	Producer sarama.SyncProducer
	// Consumer *kafka.ConsumerHandler
}

func main() {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err.Error())
	}
	log.Fatal(run(cfg))
}

func run(cfg config.Config) error {
	a, err := app.NewApp(cfg)
	if err != nil {
		return err
	}
	var dep AppDep
	if err := a.Unwrap(&dep); err != nil {
		return err
	}
	return a.Run(context.Background())
}
