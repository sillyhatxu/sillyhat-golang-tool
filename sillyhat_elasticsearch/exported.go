package sillyhat_elasticsearch

import "github.com/olivere/elastic"

var (
	// elasticClient is the name of the standard elasticsearch
	elasticClient = New()
)

func SetURL(url string) {
	elasticClient.MU.Lock()
	defer elasticClient.MU.Unlock()
	elasticClient.URL = url
}

func SetType(elasticType string) {
	elasticClient.MU.Lock()
	defer elasticClient.MU.Unlock()
	elasticClient.ElasticType = elasticType
}

func SetIndex(elasticIndex string) {
	elasticClient.MU.Lock()
	defer elasticClient.MU.Unlock()
	elasticClient.ElasticIndex = elasticIndex
}

func Bulk(bulkEntityArray []BulkEntity) (*elastic.BulkResponse, error) {
	return elasticClient.Bulk(bulkEntityArray)
}

func BulkAll(bulkEntityArray []BulkEntity) (*elastic.BulkResponse, error) {
	return elasticClient.BulkAll(bulkEntityArray)
}

func IndexExists() (bool,error) {
	return elasticClient.IndexExists()
}

func CreateIndex() (bool,error) {
	return elasticClient.CreateIndex()
}

func Index(msg string) (*elastic.IndexResponse,error) {
	return elasticClient.Index(msg)
}

func Update(id string,doc interface{}) (*elastic.UpdateResponse, error) {
	return elasticClient.Update(id,doc)
}

func Get(id string) (*elastic.GetResult, error) {
	return elasticClient.Get(id)
}

func MultiGet(idArray [] string) (*elastic.MgetResponse, error) {
	return elasticClient.MultiGet(idArray)
}

func Delete(id string) (int64, error){
	return elasticClient.Delete(id)
}

func DeleteIndex() (bool,error){
	return elasticClient.DeleteIndex()
}

func Search(msg string) (*elastic.IndexResponse,error) {
	return elasticClient.Index(msg)
}