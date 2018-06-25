package sillyhat_kafka

import (
	log "sillyhat-golang-tool/sillyhat_log/logrus"
	"time"
	"github.com/bsm/sarama-cluster"
	"github.com/Shopify/sarama"
	"strings"
	"encoding/json"
)

//type ConsumerBatchCallback func(eventArray []EventDTO) error

type EventDTO struct {

	EventName string `json:"eventName"`

	EventBody string `json:"eventBody"`

	EventTimestamp int64 `json:"eventTimestamp"`
}

type CommandDTO struct {

	CommandName string `json:"commandName"`

	CommandBody string `json:"commandBody"`

	CommandTimestamp int64 `json:"commandTimestamp"`
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

	consumerClient *cluster.Consumer
}


func NewKafkaConsumerConfig(maxWaitTime time.Duration,maxEvent int) *consumerConfig{
	return &consumerConfig{maxWaitTime:maxWaitTime,maxEvent:maxEvent}
}


func NewKafkaClient(kafkaGroupId string,kafkaNodeList string,kafkaTopicList string,consumerConfig *consumerConfig) (*kafkaClient,error) {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	c, err := cluster.NewConsumer(strings.Split(kafkaNodeList, ","), kafkaGroupId, strings.Split(kafkaTopicList, ","), config)
	if err != nil {
		log.Printf("Failed open consumer: %v", err)
		return nil,err
	}
	defer c.Close()
	return &kafkaClient{
		kafkaGroupId:kafkaGroupId,
		kafkaNodeList:kafkaNodeList,
		kafkaTopicList:kafkaTopicList,
		consumerConfig:consumerConfig,
		consumerClient:c,
	},nil
}

func (client kafkaClient) BatchConsumer() ([]EventDTO,error){
	go func() {
		for err := range client.consumerClient.Errors() {
			log.Info("Error: %s\n", err.Error())
		}
		log.Info("client.consumerClient.Errors end")
	}()

	go func() {
		for note := range client.consumerClient.Notifications() {
			log.Info("Rebalanced-------------: %v \n", note)
		}
		log.Info("client.consumerClient.ErNotificationsrors end")
	}()

	var eventArray []EventDTO

	var endTime <- chan time.Time
	endTime = time.After(client.consumerConfig.maxWaitTime)

	for {
		select {
		case msg, ok := <- client.consumerClient.Messages():
			if ok {
				//log.Debugf("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				log.Infof("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				var eventDTO EventDTO
				json.Unmarshal([]byte(string(msg.Value)), &eventDTO)
				eventArray = append(eventArray,eventDTO)
				if len(eventArray) >= client.consumerConfig.maxEvent{
					goto Loop
				}
			}
		case <-endTime:
			log.Println("Consumer timeout,return event length ",len(eventArray))
			goto Loop
		}
	}
	Loop:
	log.Info("BatchConsumer end")
	return eventArray,nil

}