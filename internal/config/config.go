package config

import (
	"github.com/ITA-Dnipro/Dp-230-Report-Service/pkg/grpc"
	"github.com/ITA-Dnipro/Dp-230-Report-Service/pkg/kafka"
	"github.com/ITA-Dnipro/Dp-230-Report-Service/pkg/mail"
)

type Config struct {
	Server     *grpc.ServerConfiguration
	Client     *grpc.ClientConfiguration
	Producer   *kafka.ProducerConfiguration
	Consumer   *kafka.ConsumeGroupConfiguration
	MailSender *mail.MailConfiguration

	FrontEndServiceBaseURL string `envconfig:"FRONT_END_BASE_URL" default:"http://localhost:8080"`
}
