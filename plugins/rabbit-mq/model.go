package rabbit_mq

import (
	"github.com/streadway/amqp"
	"sync"
)



func GetPubChan() chan string {
	return make(chan string)
}

type RabbitMQ struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
	VHost         string `json:"v_host"`
	QName         string `json:"q_name"`
	Key           string `json:"key"`
	ExName        string `json:"ex_name"`
	ExType        string `json:"ex_type"`
	PrefetchCount int    `json:"prefetch_count"`
	connection    *amqp.Connection
	channel       *amqp.Channel
	mu            sync.RWMutex
}
