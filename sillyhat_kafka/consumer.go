package sillyhat_kafka

import (
	log "sillyhat-golang-tool/sillyhat_log/logrus"
	"time"
	"github.com/bsm/sarama-cluster"
	"github.com/Shopify/sarama"
	"strings"
)

type ConsumerDTO struct {
	Key, Value     []byte

	Topic string `json:"topic"`

	Partition int32 `json:"partition"`

	Offset int64 `json:"offset"`

}

type consumerConfig struct {

	maxWaitTime time.Duration

	maxEvent int
}

type kafkaClient struct {

	kafkaGroupId string

	kafkaNodeList string

	kafkaTopicList string

	consumerConfig *consumerConfig
}


func NewKafkaConsumerConfig(maxWaitTime time.Duration,maxEvent int) *consumerConfig{
	return &consumerConfig{maxWaitTime:maxWaitTime,maxEvent:maxEvent}
}


func NewKafkaClient(kafkaGroupId string,kafkaNodeList string,kafkaTopicList string,consumerConfig *consumerConfig) (*kafkaClient,error) {
	log.Info("NewKafkaClient [ kafkaGroupId :",kafkaGroupId,"; kafkaNodeList : ",kafkaNodeList,"; kafkaTopicList : ",kafkaTopicList,"]")
	return &kafkaClient{
		kafkaGroupId:kafkaGroupId,
		kafkaNodeList:kafkaNodeList,
		kafkaTopicList:kafkaTopicList,
		consumerConfig:consumerConfig,
	},nil
}

func (client kafkaClient) getKafkaClient() (*cluster.Consumer,error) {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	//config.Consumer.Offsets.Initial = sarama.OffsetOldest
	c, err := cluster.NewConsumer(strings.Split(client.kafkaNodeList, ","), client.kafkaGroupId, strings.Split(client.kafkaTopicList, ","), config)
	if err != nil {
		log.Printf("Failed open consumer: %v", err)
		return nil,err
	}

	go func() {
		for err := range c.Errors() {
			log.Info("Error: %s\n", err.Error())
		}
		log.Info("client.consumerClient.Errors end")
	}()

	go func() {
		for note := range c.Notifications() {
			log.Info("Rebalanced-------------: %v \n", note)
		}
		log.Info("client.consumerClient.Notificationsrors end")
	}()
	return c,nil
}


func (client kafkaClient) BatchConsumer() ([]ConsumerDTO,error){
	c,err := client.getKafkaClient()
	if err != nil{
		return nil,err
	}
	defer c.Close()

	var consumerArray []ConsumerDTO

	var endTime <- chan time.Time
	endTime = time.After(client.consumerConfig.maxWaitTime)

	for {
		select {
		case msg, ok := <- c.Messages():
			if ok {
				log.Debugf("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				//log.Infof("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				consumerArray = append(consumerArray,*&ConsumerDTO{Key:msg.Key,Value:msg.Value,Partition:msg.Partition,Offset:msg.Offset,Topic:msg.Topic})
				if len(consumerArray) >= client.consumerConfig.maxEvent{
					log.Println("Consumer max event:",len(consumerArray))
					goto Loop
				}
			}
		case <-endTime:
			log.Println("Consumer timeout,return length ",len(consumerArray))
			goto Loop
		}
	}
	Loop:
	log.Info("BatchConsumer end")
	return consumerArray,nil
}

func (client kafkaClient) BatchCommit(consumerArray []ConsumerDTO) (error){
	c,err := client.getKafkaClient()
	if err != nil{
		return err
	}
	defer c.Close()
	time.Sleep(10*time.Second)
	for _,consumerDTO := range consumerArray{
		c.MarkPartitionOffset(consumerDTO.Topic,consumerDTO.Partition,consumerDTO.Offset,"")
	}
	c.CommitOffsets()
	return nil
}