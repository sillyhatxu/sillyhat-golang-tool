package sillyhat_elasticsearch

import (
	"sync"
	"bytes"
	"github.com/olivere/elastic"
	"context"
	"time"
	"os"
	"log"
	"github.com/pkg/errors"
)

var bufferPool *sync.Pool

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}

type ElasticEntry struct {

	ctx context.Context

	elasticClient *elastic.Client

	elasticType string

	elasticIndex string

}

type BulkableDoc interface {
	Source() ([]string, error)
}


func NewElastic(client *Client) (*ElasticEntry,error) {
	elasticClient, err := elastic.NewClient(
		elastic.SetURL(client.url),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil{
		return nil,err
	}
	return &ElasticEntry{
		ctx:client.ctx,
		elasticClient:elasticClient,
		elasticType:client.elasticType,
		elasticIndex:client.elasticIndex,
	},nil
}

func (elasticEntry ElasticEntry) Exists() (bool,error) {
	exists, err := elasticEntry.elasticClient.IndexExists(elasticEntry.elasticIndex).Do(elasticEntry.ctx)
	if err != nil {
		return false,err
	}
	return exists,nil
}

func (elasticEntry ElasticEntry) CreateIndex() (bool,error) {
	index, err := elasticEntry.elasticClient.CreateIndex(elasticEntry.elasticIndex).Do(elasticEntry.ctx)
	if err != nil {
		return false,err
	}
	if !index.Acknowledged {
		return false,errors.New("Not acknowledged")
	}
	return true,nil
}

func (elasticEntry ElasticEntry) Index(json string) (*elastic.IndexResponse,error) {
	exists,err := elasticEntry.Exists()
	if err != nil {
		return nil,err
	}
	if !exists {
		_,err := elasticEntry.CreateIndex()
		if err != nil{
			return nil,err
		}
	}
	index, err := elasticEntry.elasticClient.Index().Index(elasticEntry.elasticIndex).Type(elasticEntry.elasticType).BodyJson(json).Do(elasticEntry.ctx)
	if err != nil {
		return nil,err
	}
	return index,nil
}

func (elasticEntry ElasticEntry) Bulk(docArray [] interface{},idDeleteArray [] string) (*elastic.BulkResponse, error) {
	bulk := elasticEntry.elasticClient.Bulk()
	for _,doc := range docArray{
		request := elastic.NewBulkIndexRequest().Index(elasticEntry.elasticIndex).Type(elasticEntry.elasticType).Doc(doc)
		bulk.Add(request)
	}
	//for _,doc := range docUpdateArray{
	//	request := elastic.NewBulkUpdateRequest().Index(elasticEntry.elasticIndex).Type(elasticEntry.elasticType).Doc(doc)
	//	bulk.Add(request)
	//}
	for _,id := range idDeleteArray{
		request := elastic.NewBulkDeleteRequest().Index(elasticEntry.elasticIndex).Type(elasticEntry.elasticType).Id(id)
		bulk.Add(request)
	}
	bulkResponse, err := bulk.Do(elasticEntry.ctx)
	if err != nil {
		return nil,err
	}
	return bulkResponse,nil
}

func (elasticEntry ElasticEntry) Get(id string) (*elastic.GetResult, error) {
	getResult, err := elasticEntry.elasticClient.Get().Index(elasticEntry.elasticIndex).Type(elasticEntry.elasticType).Id(id).Do(elasticEntry.ctx)
	if err != nil {
		return nil,err
	}
	return getResult,nil
}

func (elasticEntry ElasticEntry) Update(id string,doc interface{}) (*elastic.UpdateResponse, error) {
	//.Script(elastic.NewScriptInline("ctx._source.retweets += params.num").Lang("painless").Param("num", 1)).
	//update, err := elasticEntry.elasticClient.Update().Index(elasticEntry.elasticIndex).Type(elasticEntry.elasticType).Id(id).Upsert(map[string]interface{}{"retweets": 0}).Do(elasticEntry.ctx)
	update, err := elasticEntry.elasticClient.Update().Index(elasticEntry.elasticIndex).Type(elasticEntry.elasticType).Id(id).Upsert(doc).Do(elasticEntry.ctx)
	if err != nil {
		return nil,err
	}
	return update,nil
}

func (elasticEntry ElasticEntry) Delete(id string) (int64, error) {
	termQuery := elastic.NewTermQuery("id", id)
	deleteResponse, err := elasticEntry.elasticClient.DeleteByQuery(elasticEntry.elasticIndex).Query(termQuery).Do(elasticEntry.ctx)
	if err != nil {
		return 0,err
	}
	return deleteResponse.Deleted,nil
}

func (elasticEntry ElasticEntry) DeleteIndex() (bool, error) {
	deleteIndex, err := elasticEntry.elasticClient.DeleteIndex(elasticEntry.elasticIndex).Do(elasticEntry.ctx)
	if err != nil {
		return false,err
	}
	if !deleteIndex.Acknowledged {
		return false,errors.New("Not acknowledged")
	}
	return true,nil
}