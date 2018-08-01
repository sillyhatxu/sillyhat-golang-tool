package sillyhat_consul

import (
	log "sillyhat-golang-tool/sillyhat_log/logrus"
	consulapi "github.com/hashicorp/consul/api"
	consulconnect "github.com/hashicorp/consul/connect"
	"fmt"
	"math/rand"
	"encoding/json"
	"io"
)

const (
	RESPONSE_SUCCESS int = 0
	RESPONSE_ERROR   int = -1
	CONTENT_TYPE string = "application/json"
)

type ResponseEntity struct {
	Ret int `json:"ret"`
	Data interface{} `json:"data"`
	Msg string `json:"msg"`
}

func Success(data interface{}) *ResponseEntity {
	return &ResponseEntity{Ret:RESPONSE_SUCCESS,Data:data,Msg:"success"}
}

func ErrorData(ret int,data interface{},msg string) *ResponseEntity {
	return &ResponseEntity{Ret:ret,Data:data,Msg:msg}
}

func Error(ret int,msg string) *ResponseEntity {
	return &ResponseEntity{Ret:ret,Msg:msg}
}

type ConsulClient struct{
	Name,Tag,MethodURI string
	Body io.Reader
	client *consulapi.Client
}

func NewConsulClient(consulAddress, module, methodURI string, body io.Reader) (*ConsulClient,error) {
	log.Info("NewConsulClient")
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil{
		return nil,err
	}
	return &ConsulClient{
		client:    client,
		Name:      module,
		Tag:       module,
		MethodURI: methodURI,
		Body:      body,
	},nil
}

func (c ConsulClient) getService() (*consulapi.ServiceEntry, error) {
	options := &consulapi.QueryOptions{AllowStale: false}
	services, _, err := c.client.Health().Service(c.Name, c.Tag, true, options)
	if err != nil {
		return nil, err
	}
	if services != nil && len(services) < 1 {
		return nil, fmt.Errorf("error, no services available for %q with tag %q", c.Name, c.Tag)
	}
	return services[rand.Intn(len(services))],nil
}

func (c ConsulClient) GetModuleHttpSrc() (string, error) {
	service,err := c.getService()
	if err != nil{
		return "",err
	}
	return fmt.Sprintf("http://%s:%d", service.Service.Address, service.Service.Port),nil
}


func (c ConsulClient) DoGet() (*ResponseEntity, error) {
	svc, _ := consulconnect.NewService(c.Name, c.client)
	defer svc.Close()
	httpClient := svc.HTTPClient()
	httpSrc,err := c.GetModuleHttpSrc()
	if err != nil{
		return nil,err
	}
	log.Info("http : ",httpSrc)
	log.Info("MethodURI : ",c.MethodURI)
	response, err := httpClient.Get(httpSrc + c.MethodURI)
	log.Info("response : ",response)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 200{
		responseEntity := &ResponseEntity{}
		err := json.NewDecoder(response.Body).Decode(responseEntity)
		if err != nil{
			log.Error("Response to json error.",err)
			return nil,err
		}
		return responseEntity,nil
	}
	return nil,fmt.Errorf("Unknow exception.")
}

func (c ConsulClient) DoPost() (*ResponseEntity, error) {
	svc, _ := consulconnect.NewService(c.Name, c.client)
	defer svc.Close()
	httpClient := svc.HTTPClient()
	httpSrc,err := c.GetModuleHttpSrc()
	if err != nil{
		return nil,err
	}
	response, err := httpClient.Post(httpSrc + c.MethodURI,CONTENT_TYPE,c.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 200{
		responseEntity := &ResponseEntity{}
		err := json.NewDecoder(response.Body).Decode(responseEntity)
		if err != nil{
			log.Error("Response to json error.",err)
			return nil,err
		}
		return responseEntity,nil
	}
	return nil,fmt.Errorf("Unknow exception.")
}

func (c ConsulClient) DoDelete() (*ResponseEntity, error) {
	return c.DoPost()
}

func (c ConsulClient) DoPut() (*ResponseEntity, error) {
	return c.DoPost()
}