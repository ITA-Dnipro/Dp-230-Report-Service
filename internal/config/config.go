package config

import (
	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/grpc"
	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/kafka"
)

type Config struct {
	Server   *grpc.ServerConfiguration
	Client   *grpc.ClientConfiguration
	Producer *kafka.ProducerConfiguration
	Consumer *kafka.ConsumeGroupConfiguration
}
