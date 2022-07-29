package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/model"
	"github.com/Shopify/sarama"
)

type Sender interface {
	Send(context.Context, model.Report) error
}

type handler struct {
	sender Sender
}

func NewConsumerHandler(sender Sender) *handler {
	return &handler{sender: sender}
}

func (h *handler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var report model.Report

	if err := json.Unmarshal(msg.Value, &report); err != nil {
		//TODO: replace to logger from context
		fmt.Println("error", err)
		return nil
	}

	if err := h.sender.Send(ctx, report); err != nil {
		//TODO: replace to logger from context
		fmt.Println("error", err)
		return nil
	}
	return nil
}
