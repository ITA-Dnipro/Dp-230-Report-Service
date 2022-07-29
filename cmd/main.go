package main

import (
	"context"
	"log"

	"github.com/kelseyhightower/envconfig"

	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/app"
	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/config"
	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/report/kafka/handler"
	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/report/mail"
)

func main() {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err.Error())
	}
	log.Fatal(run(cfg))
}

func run(cfg config.Config) error {
	app, err := app.NewApp(cfg)
	if err != nil {
		return err
	}
	sender := mail.NewReportMailSender(app.MailSender, cfg.FrontEndServiceBaseURL)
	app.Consumer.Handler = handler.NewConsumerHandler(sender)

	return app.Run(context.Background())
}
