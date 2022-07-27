package main

import (
	"context"
	"log"

	"github.com/kelseyhightower/envconfig"

	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/app"
	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/config"
	reportHandler "github.com/ITA-Dnipro/Dp-230-Report-Service/internal/report/kafka/handler"
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

	app.Consumer.Handler = reportHandler.New(app.MailSender)

	return app.Run(context.Background())
}
