package sillyhat_elasticsearch

import (
	"sync"
	"context"
	"github.com/olivere/elastic"
)

type Client struct {

	URL string

	ElasticType string

	ElasticIndex string

	CTX context.Context

	MU MutexWrap

	ClientPool sync.Pool
}

type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

func (client *Client) newElastic() (*ElasticEntry,error) {
	elasticEntry, ok := client.ClientPool.Get().(*ElasticEntry)
	if ok {
		return elasticEntry,nil
	}
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
		CTX:context.Background(),
	}
}

func (client *Client) releaseEntry(elasticEntry *ElasticEntry) {
	client.ClientPool.Put(elasticEntry)
}

func (client *Client) IndexExists() (bool,error) {
	elasticEntry,err := client.newElastic()
	if err != nil{
		return false,err
	}
	client.releaseEntry(elasticEntry)
	return elasticEntry.IndexExists()
}

func (client *Client) CreateIndex() (bool,error) {
	elasticEntry,err := client.newElastic()
	if err != nil{
		return false,err
	}
	client.releaseEntry(elasticEntry)
	return elasticEntry.CreateIndex()
}

func (client *Client) Index(msg string) (*elastic.IndexResponse,error) {
	elasticEntry,err := client.newElastic()
	if err != nil{
		return nil,err
	}
	client.releaseEntry(elasticEntry)
	return elasticEntry.Index(msg)
}

func (client *Client) Bulk(bulkEntityArray []BulkEntity) (*elastic.BulkResponse, error) {
	elasticEntry,err := client.newElastic()
	if err != nil{
		return nil,err
	}
	client.releaseEntry(elasticEntry)
	return elasticEntry.Bulk(bulkEntityArray)
}

func (client *Client) BulkAll(bulkEntityArray []BulkEntity) (*elastic.BulkResponse, error) {
	elasticEntry,err := client.newElastic()
	if err != nil{
		return nil,err
	}
	client.releaseEntry(elasticEntry)
	return elasticEntry.BulkAll(bulkEntityArray)
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


func (client *Client) MultiGet(idArray []string) (*elastic.MgetResponse, error) {
	elasticEntry,err := client.newElastic()
	if err != nil{
		return nil,err
	}
	client.releaseEntry(elasticEntry)
	return elasticEntry.MultiGet(idArray)
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