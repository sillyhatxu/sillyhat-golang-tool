package main

import (
	"github.com/bsm/sarama-cluster"
	"time"
	"github.com/Shopify/sarama"
	"strings"
	"log"
	"sillyhat-golang-tool/sillyhat_kafka"
)

func get()  {
	log.Println("start kafka service")
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	//config.Consumer.Offsets.Initial = sarama.OffsetNewest //初始从最新的offset开始
	config.Consumer.Fetch.Max = 10
	config.Consumer.MaxWaitTime = 5 * time.Second

	groupID := "ocb-syncer-11"
	nodeList := "172.28.2.22:9092,172.28.2.22:9091,172.28.2.22:9090"
	topicList := "DataRefreshed"

	c, err := cluster.NewConsumer(strings.Split(nodeList, ","), groupID, strings.Split(topicList, ","), config)
	if err != nil {
		log.Printf("Failed open consumer: %v", err)
		return
	}
	defer c.Close()

	subscriptions := c.Subscriptions()
	log.Println("----------")
	log.Println(subscriptions)
	log.Println("----------")

	go func() {
		for err := range c.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	go func() {
		for note := range c.Notifications() {
			log.Printf("Rebalanced-------------: %v \n", note)
		}
	}()

	//2018/06/22 22:45:59 start kafka service
	//2018/06/22 22:45:59 Rebalanced-------------: &{rebalance start map[] map[] map[]}
	//2018/06/22 22:46:00 message start
	//2018/06/22 22:46:22 Rebalanced-------------: &{rebalance OK map[DataRefreshed:[0 1]] map[] map[DataRefreshed:[0 1]]}

	//time.Sleep(1*time.Second)
	i := 0
	log.Println("message start")
	for msg := range c.Messages() {
		log.Println("----------message----------")
		log.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		log.Println()
		//concumer(string(msg.Value))
		//c.MarkOffset(msg, "") //MarkOffset 并不是实时写入kafka，有可能在程序crash时丢掉未提交的offset
		i++
		if i == 10 {
			break
		}
	}
	sub := c.Subscriptions()
	log.Println("----------")
	log.Println(sub)
	log.Println("----------")

	log.Println("message end")
}

func commit(){
	log.Println("start kafka service")

	groupID := "ocb-syncer-11"
	nodeList := "172.28.2.22:9092,172.28.2.22:9091,172.28.2.22:9090"
	topicList := "DataRefreshed"

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	cu, err := cluster.NewConsumer(strings.Split(nodeList, ","), groupID, strings.Split(topicList, ","), config)
	if err != nil {
		log.Printf("Failed open consumer: %v", err)
		return
	}

	defer cu.Close()

	go func() {
		for err := range cu.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	go func() {
		for note := range cu.Notifications() {
			log.Printf("Rebalanced-------------: %v \n", note)
		}
	}()

	time.Sleep(10*time.Second)

	//concumer.HighWaterMarks()
	subscriptions := cu.Subscriptions()
	log.Println("----------")
	log.Println(subscriptions)
	log.Println("----------")

	cu.MarkPartitionOffset(topicList,0,258,"")
	cu.MarkPartitionOffset(topicList,1,198,"")
	cu.CommitOffsets()

	log.Println("commit")


}

//subscriptionsOf(cs).Should(Equal(map[string][]int32{
//"topic-a": {0, 1, 2, 3}},
//))

//func (c *cluster.Consumer) subscriptionsOf cluster.GomegaAsyncAssertion {
//	return Eventually(func() map[string][]int32 {
//		return c.Subscriptions()
//	}, "10s", "100ms")
//}

func testBatchConsumer()  {
	groupID := "ocb-syncer-11"
	nodeList := "172.28.2.22:9092,172.28.2.22:9091,172.28.2.22:9090"
	topicList := "DataRefreshed"
	config := sillyhat_kafka.NewKafkaConsumerConfig(60 * time.Second,10000)
	client,err := sillyhat_kafka.NewKafkaClient(groupID,nodeList,topicList,config)
	if err != nil{
		log.Println(err)
		return
	}
	eventArray,err := client.BatchConsumer()
	for _,event := range eventArray{
		log.Println(event)
	}

}
func main() {
	get()
	//commit()
	//testBatchConsumer()
	//done := make(chan os.Signal, 1)
	//signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	//<-done
	//:0, Offset:257	0, Offset:266 93267
	//1, Offset:197	:1, Offset:203,	1023
}
