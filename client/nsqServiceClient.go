package client

import (
	"github.com/Diode222/GomokuGameReferee/conf"
	"github.com/nsqio/go-nsq"
)

func NewNsqClient() (*nsq.Producer, error) {
	producer, err := nsq.NewProducer(conf.NSQ_ADDR, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	return producer, err
}
