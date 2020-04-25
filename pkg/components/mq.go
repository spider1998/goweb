package components

import (
	"github.com/nsqio/go-nsq"
	"time"
)

type NsqMq struct {
	sender *Sender
}

type Consumer struct {
	Logger
	addr    string
	queue   string
	channel string
	config  *nsq.Config
}

func NewMq(addr string, logger Logger) *NsqMq {
	return &NsqMq{
		sender: NewSender(addr, logger),
	}
}

func (mq *NsqMq) Consume(addr, queue string, channel string, logger Logger, handle func(msq *nsq.Message) error) {
	go NewConsumer(addr, queue, channel, logger, nsq.NewConfig()).Consumer(handle)
}

func (mq *NsqMq) Send(queue string, payload []byte) {
	mq.sender.Send(queue, payload)
}

func NewConsumer(addr, queue, channel string, logger Logger, config *nsq.Config) *Consumer {
	c := new(Consumer)
	c.Logger = logger
	c.addr = addr
	c.queue = queue
	c.channel = channel
	if config == nil {
		config = nsq.NewConfig()
	}
	c.config = config
	return c
}

func (c *Consumer) Consumer(handler func(msg *nsq.Message) error) {
	if handler == nil {
		panic("consumer handler is nil")
	}
	for {
		c.consume(handler)
		<-time.After(time.Second)
	}
}

func (c *Consumer) consume(handler func(msg *nsq.Message) error) {
	c.Infof("start nsq runner.")
	consumer, err := nsq.NewConsumer(c.queue, c.channel, c.config)
	if err != nil {
		c.Errorf("fail to create nsq consumer.", "error", err)
		return
	}
	consumer.SetLogger(NewNSQLogger(c), nsq.LogLevelWarning)
	consumer.AddHandler(nsq.HandlerFunc(handler))
	err = consumer.ConnectToNSQD(c.addr)
	if err != nil {
		c.Errorf("fail to connect to nsqd.", "error", err)
		return
	}
	<-consumer.StopChan
	c.Infof("nsq runner stopped.")
}

type Sender struct {
	Logger
	addr  string
	msgCH chan message
}

type message struct {
	queue   string
	payload []byte
}

func NewSender(addr string, logger Logger) *Sender {
	s := new(Sender)
	s.Logger = logger
	s.addr = addr
	s.msgCH = make(chan message, 1024)
	go s.runLogSender()
	return s
}

func (s *Sender) Send(queue string, payload []byte) {
	s.Infof("enqueue message.", "queue", queue)
	select {
	case s.msgCH <- message{queue: queue, payload: payload}:
	default:
		s.Warnf("mq channel is full.")
	}
}

func (s *Sender) runLogSender() {
	for {
		func() {
			s.Infof("run mq sender.")
			producer, err := nsq.NewProducer(s.addr, nsq.NewConfig())
			if err != nil {
				s.Errorf("fail to create new nsq producer.", "error", err)
				return
			}
			defer producer.Stop()
			producer.SetLogger(NewNSQLogger(s), nsq.LogLevelWarning)

			for msg := range s.msgCH {
				s.Debugf("send msg.", "queue", msg.queue, "payload", string(msg.payload))
				s.Infof("send msg.", "queue", msg.queue)
				err = producer.Publish(msg.queue, msg.payload)
				if err != nil {
					s.Errorf("fail to publish log message.", "error", err)
					return
				}
			}
		}()
		time.Sleep(time.Second)
	}
}

func NewNSQLogger(logger Logger) NSQLogger {
	return NSQLogger{logger}
}

type NSQLogger struct {
	logger Logger
}

func (logger NSQLogger) Output(calldepth int, s string) error {
	logger.logger.Infof(s, "service", "nsq_logger")
	return nil
}
