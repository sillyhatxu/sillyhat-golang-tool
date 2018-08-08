package sillyhat_kafka

import (
	log "sillyhat-golang-tool/sillyhat_log/logrus"
	"github.com/Shopify/sarama"
	"strings"
	"time"
	"sillyhat-golang-tool/sillyhat_kafka/command"
	"encoding/json"
)

type Producer struct {

	Topics string

	Address string

	config *sarama.Config
}

func NewKafkaProducerClient(address,topics string) *Producer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true //必须有这个选项
	config.Producer.Timeout = 5 * time.Second
	return &Producer{Address:address,Topics:topics,config:config}
}

func (p Producer) Send(commandName,commandBody string) error {
	commandDTO := command.NewCommandDTO(commandName,commandBody)
	commandByte, err := json.Marshal(commandDTO)
	if err != nil {
		log.Error("CommandDTO to json error",err.Error())
		return err
	}
	kafkaProducer,err := sarama.NewAsyncProducer(strings.Split(p.Address, ","), p.config)
	defer kafkaProducer.Close()
	if err != nil {
		return err
	}
	//必须有这个匿名函数内容
	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					log.Error(err)
				}
			case <-success:
			}
		}
	}(kafkaProducer)
	msg := &sarama.ProducerMessage{
		Topic: p.Topics,
		Value: sarama.ByteEncoder(commandByte),
	}
	kafkaProducer.Input() <- msg
	return nil
}

//func (p Producer) Send(message string) error {
//	kafkaProducer,err := sarama.NewAsyncProducer(strings.Split(p.Address, ","), p.config)
//	defer kafkaProducer.Close()
//	if err != nil {
//		return err
//	}
//	//必须有这个匿名函数内容
//	go func(p sarama.AsyncProducer) {
//		errors := p.Errors()
//		success := p.Successes()
//		for {
//			select {
//				case err := <-errors:
//					if err != nil {
//						log.Error(err)
//					}
//				case <-success:
//			}
//		}
//	}(kafkaProducer)
//	msg := &sarama.ProducerMessage{
//		Topic: p.Topics,
//		Value: sarama.ByteEncoder(message),
//	}
//	kafkaProducer.Input() <- msg
//	return nil
//}