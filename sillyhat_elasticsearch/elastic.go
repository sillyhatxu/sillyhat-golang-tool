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

type BulkEntity struct {

	Id string

	Data interface{}

	IsDelete bool

}
type BulkDoc interface {

	InsertArray() [] interface{}

	UpdateArray() [] interface{}

	DeleteArray() [] string

}


func NewElastic(client *Client) (*ElasticEntry,error) {
	elasticClient, err := elastic.NewClient(
		elastic.SetURL(client.URL),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil{
		return nil,err
	}
	return &ElasticEntry{
		ctx:client.CTX,
		elasticClient:elasticClient,
		elasticType:client.ElasticType,
		elasticIndex:client.ElasticIndex,
	},nil
}

func (elasticEntry ElasticEntry) IndexExists() (bool,error) {
	exists, err := elasticEntry.elasticClient.IndexExists(elasticEntry.elasticIndex).Do(elasticEntry.ctx)
	if err != nil {
		return false,err
	}
	return exists,nil
}

func (elasticEntry ElasticEntry) DataExists(id string) (bool,error) {
	exists, err := elasticEntry.elasticClient.Exists().Index(elasticEntry.elasticIndex).Type(elasticEntry.elasticType).Id(id).Do(elasticEntry.ctx)
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
	exists,err := elasticEntry.IndexExists()
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

func (elasticEntry ElasticEntry) BulkAll(bulkEntityArray []BulkEntity) (*elastic.BulkResponse, error) {
	bulk := elasticEntry.elasticClient.Bulk()
	for _,bulkEntity := range bulkEntityArray{
		request := elastic.NewBulkIndexRequest().Index(elasticEntry.elasticIndex).Type(elasticEntry.elasticType).Id(bulkEntity.Id).Doc(bulkEntity.Data)
		bulk.Add(request)
	}
	bulkResponse, err := bulk.Do(elasticEntry.ctx)
	if err != nil {
		return nil,err
	}
	return bulkResponse,nil
}

func (elasticEntry ElasticEntry) Bulk(bulkEntityArray []BulkEntity) (*elastic.BulkResponse, error) {
	bulk := elasticEntry.elasticClient.Bulk()
	for _,bulkEntity := range bulkEntityArray{
		if bulkEntity.IsDelete{
			exists,err := elasticEntry.DataExists(bulkEntity.Id)
			if err != nil{
				log.Println(err.Error())
			}
			if exists{
				request := elastic.NewBulkDeleteRequest().Index(elasticEntry.elasticIndex).Type(elasticEntry.elasticType).Id(bulkEntity.Id)
				bulk.Add(request)
			}
		}else{
			exists,err := elasticEntry.DataExists(bulkEntity.Id)
			if err != nil{
				log.Println(err.Error())
			}
			if exists{
				request := elastic.NewBulkUpdateRequest().Index(elasticEntry.elasticIndex).Type(elasticEntry.elasticType).Id(bulkEntity.Id).Doc(bulkEntity.Data)
				bulk.Add(request)
			}else{
				request := elastic.NewBulkIndexRequest().Index(elasticEntry.elasticIndex).Type(elasticEntry.elasticType).Id(bulkEntity.Id).Doc(bulkEntity.Data)
				bulk.Add(request)
			}
		}
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