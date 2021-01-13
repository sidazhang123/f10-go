package rabbit_mq

import (
	"fmt"
	"github.com/sidazhang123/f10-go/basic"
	"go.uber.org/zap"

	log2 "github.com/micro/go-micro/v2/util/log"
	"github.com/sidazhang123/f10-go/basic/config"
	"github.com/streadway/amqp"
	"sync"

	z "github.com/sidazhang123/f10-go/plugins/zap"
)

var (
	m      sync.RWMutex
	log    = z.GetLogger()
	inited bool
	mq     *RabbitMQ
)

func init() {
	basic.Register(initRMQ)
}

type Subscriber interface {
	Consumer(msg amqp.Delivery, errList []error, mux *sync.Mutex)
}

func initRMQ() {
	m.Lock()
	defer m.Unlock()
	if inited {

		log.Error(fmt.Errorf("[Init] RMQ was initialized").Error())
		return
	}

	mq = &RabbitMQ{}
	err := config.C().App("rabbitmq", mq)
	if err != nil {
		panic(err)
	}
	mq.mqConnect()
	inited = true
}

var mqConn *amqp.Connection
var mqChan *amqp.Channel

func (r *RabbitMQ) mqConnect() {
	var err error
	RabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", r.Username, r.Password, r.Host, r.Port, r.VHost)
	mqConn, err = amqp.Dial(RabbitUrl)
	r.connection = mqConn
	if err != nil {
		log.Error("Failed to open MQ conn", zap.Any("err", err))
		return
	}
	mqChan, err = mqConn.Channel()
	if err != nil {
		log.Error("Failed to open MQ channel", zap.Any("err", err))
		return
	}
	r.channel = mqChan

	err = r.channel.ExchangeDeclare(r.ExName, r.ExType, false, false, false, false, nil)
	if err != nil {
		log.Error("Failed to declare exchange", zap.Any("err", err))
		return
	}

	_, err = r.channel.QueueDeclare(r.QName, false, false, false, false, nil)
	if err != nil {
		log.Error("Failed to declare queue", zap.Any("err", err))
		return
	}

	err = r.channel.QueueBind(r.QName, r.Key, r.ExName, false, nil)
	if err != nil {
		log.Error("Failed to bind queue with exchange", zap.Any("err", err))
		return
	}

	err = r.channel.Qos(r.PrefetchCount, 0, true)
	if err != nil {
		log.Error("Failed to set QoS", zap.Any("err", err))
		return
	}

}

func (r *RabbitMQ) mqClose() {
	if !r.connection.IsClosed() {
		err := r.channel.Close()
		if err != nil {
			log.Error("Failed to close MQ channel", zap.Any("err", err))
		}
		err = r.connection.Close()
		if err != nil {
			log.Error("Failed to close MQ conn", zap.Any("err", err))
		}
	}
}

func GetRMQ() *RabbitMQ {
	if inited == true && mq.connection != nil && !mq.connection.IsClosed() {
		return mq
	} else {
		mq.mqConnect()
		return mq
	}
}
func (r *RabbitMQ) PurgeQueue() (int, error) {
	return r.channel.QueuePurge(r.QName, false)
}
func (r *RabbitMQ) RegisterPublisher(pub chan string) {

	if r.connection == nil || r.connection.IsClosed() {
		r.mqConnect()
	}

	go func() {
		for i := range pub {
			err := r.channel.Publish(r.ExName, r.Key, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(i),
			})
			if err != nil {
				log.Error("Failed to publish msg", zap.Any("err", err.Error()))
			}
		}
	}()
}

func (r *RabbitMQ) RegisterSubscriber(receiver Subscriber, errList []error, sigChan chan struct{}) (e error) {
	// close conn when all messages are processed
	if r.connection == nil || r.connection.IsClosed() {
		r.mqConnect()
	}
	log2.Info(fmt.Sprintf("rabbitmq:142 r.connection.IsClosed() %t", r.connection.IsClosed()))
	msgList, e := r.channel.Consume(r.QName, "", false, false, false, false, nil)
	log2.Info(fmt.Sprintf("rabbitmq:144 len(msgList) %d", len(msgList)))
	if e != nil {
		log.Error("Failed to consume channel", zap.Any("err", e.Error()))
		return
	}
	mux := sync.Mutex{}
loop:
	for {
		select {
		case msg := <-msgList:
			{
				log2.Info(fmt.Sprintf("rabbitmq:154 msg.body %s", string(msg.Body)))
				receiver.Consumer(msg, errList, &mux)
			}
		case <-sigChan:
			break loop
		}
	}
	e = r.connection.Close()

	return
}

func Ack(msg amqp.Delivery) {

	err := msg.Ack(false)
	if err != nil {
		log.Error("ACK Exception", zap.Any("err", err.Error()))

		return
	}

}
func Nack(msg amqp.Delivery) {
	err := msg.Nack(false, true)
	if err != nil {
		log.Error("NACK Exception", zap.Any("err", err.Error()))

		return
	}
}
