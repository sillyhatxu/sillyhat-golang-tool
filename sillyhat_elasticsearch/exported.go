package sillyhat_elasticsearch

import "github.com/olivere/elastic"

var (
	// elasticClient is the name of the standard elasticsearch
	elasticClient = New()
)

func SetURL(url string) {
	elasticClient.mu.Lock()
	defer elasticClient.mu.Unlock()
	elasticClient.url = url
}

func SetType(elasticType string) {
	elasticClient.mu.Lock()
	defer elasticClient.mu.Unlock()
	elasticClient.elasticType = elasticType
}

func SetIndex(elasticIndex string) {
	elasticClient.mu.Lock()
	defer elasticClient.mu.Unlock()
	elasticClient.elasticIndex = elasticIndex
}

func Bulk(docArray [] interface{},idDeleteArray [] string) (*elastic.BulkResponse, error) {
	return elasticClient.Bulk(docArray,idDeleteArray)
}

func Exists() (bool,error) {
	return elasticClient.Exists()
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

func Delete(id string) (int64, error){
	return elasticClient.Delete(id)
}

func DeleteIndex() (bool,error){
	return elasticClient.DeleteIndex()
}

func Search(msg string) (*elastic.IndexResponse,error) {
	return elasticClient.Index(msg)
}