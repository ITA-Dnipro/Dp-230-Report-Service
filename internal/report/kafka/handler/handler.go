package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/report/model"
	"github.com/Shopify/sarama"
)

type Sender interface {
	Send(to string, message []byte) error
}

type handler struct {
	sender Sender
}

func New(sender Sender) *handler {
	return &handler{sender: sender}
}

func (h *handler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var report model.Report
	if err := json.Unmarshal(msg.Value, &report); err != nil {
		fmt.Println("error", err)
		return nil
	}
	fmt.Println(report)
	return nil
}
