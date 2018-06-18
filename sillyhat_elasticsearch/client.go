package sillyhat_elasticsearch

import (
	"sync"
	"context"
	"github.com/olivere/elastic"
	"fmt"
)



type Client struct {

	url string

	elasticType string

	elasticIndex string

	ctx context.Context

	mu MutexWrap

	clientPool sync.Pool
}

type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

func (client *Client) newElastic() (*ElasticEntry,error) {
	elasticEntry, ok := client.clientPool.Get().(*ElasticEntry)
	if ok {
		return elasticEntry,nil
	}
	fmt.Println("new elastic entry")
	return NewElastic(client)
}

func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

func New() *Client {
	return &Client{
		ctx:context.Background(),
	}
}

func (client *Client) releaseEntry(elasticEntry *ElasticEntry) {
	client.clientPool.Put(elasticEntry)
}

func (client *Client) Exists() (bool,error) {
	return elasticClient.Exists()
}

func (client *Client) CreateIndex() (bool,error) {
	return elasticClient.CreateIndex()
}

func (client *Client) Index(msg string) (*elastic.IndexResponse,error) {
	elasticEntry,err := client.newElastic()
	if err != nil{
		return nil,err
	}
	client.releaseEntry(elasticEntry)
	return elasticEntry.Index(msg)
}

func (client *Client) Bulk(docArray [] interface{},idDeleteArray [] string) (*elastic.BulkResponse, error) {
	elasticEntry,err := client.newElastic()
	if err != nil{
		return nil,err
	}
	client.releaseEntry(elasticEntry)
	return elasticEntry.Bulk(docArray,idDeleteArray)
}

func (client *Client) Update(id string,doc interface{}) (*elastic.UpdateResponse, error) {
	elasticEntry,err := client.newElastic()
	if err != nil{
		return nil,err
	}
	client.releaseEntry(elasticEntry)
	return elasticEntry.Update(id,doc)
}

func (client *Client) Get(id string) (*elastic.GetResult, error) {
	elasticEntry,err := client.newElastic()
	if err != nil{
		return nil,err
	}
	client.releaseEntry(elasticEntry)
	return elasticEntry.Get(id)
}

func (client *Client) Delete(id string) (int64, error) {
	elasticEntry,err := client.newElastic()
	if err != nil{
		return 0,err
	}
	client.releaseEntry(elasticEntry)
	return elasticEntry.Delete(id)
}

func (client *Client) DeleteIndex() (bool,error) {
	elasticEntry,err := client.newElastic()
	if err != nil{
		return false,err
	}
	client.releaseEntry(elasticEntry)
	return elasticEntry.DeleteIndex()
}