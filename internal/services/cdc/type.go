package cdc

import (
	"github.com/SyaibanAhmadRamadhan/event-bus/debezium"
	"github.com/mini-e-commerce-microservice/auth-service/internal/model"
)

type DebeziumPayload[T any] struct {
	Payload T `json:"payload"`
}

type UserData struct {
	model.User
	Op debezium.Operation `json:"op"`
}
