package controller

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"orderservice/internal/model"
	"orderservice/internal/usecase"
	"orderservice/pkg/consumer"
)

type KafkaController struct {
	uc       usecase.OrderUsecase
	consumer *consumer.Consumer
}

func NewKafkaController(uc usecase.OrderUsecase, cons *consumer.Consumer) *KafkaController {
	return &KafkaController{
		uc:       uc,
		consumer: cons,
	}
}

func (kc *KafkaController) Start(ctx context.Context) error {
	log.Println("kafka controller: starting consumer loop")

	err := kc.consumer.Consume(ctx, func(ctx context.Context, key, value []byte) error {
		return kc.handleMessage(ctx, key, value)
	})

	if err != nil {
		log.Printf("kafka controller: consumer stopped with error: %v", err)
		return err
	}
	return nil
}

func (kc *KafkaController) Stop(ctx context.Context) error {
	log.Println("kafka controller: stopping...")
	if err := kc.consumer.Close(); err != nil {
		log.Printf("kafka controller: close error: %v", err)
		return err
	}
	return nil
}

func (kc *KafkaController) handleMessage(ctx context.Context, key, value []byte) error {
	var ord model.Order
	if err := json.Unmarshal(value, &ord); err != nil {
		return err
	}

	processCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := kc.uc.CreateOrder(processCtx, &ord); err != nil {
		return err
	}

	log.Printf("kafka controller: order %s processed successfully", ord.OrderUID)
	return nil
}
