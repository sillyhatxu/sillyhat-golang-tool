package sillyhat_http

import (
	log "sillyhat-golang-tool/sillyhat_log/logrus"
	"time"
	"io/ioutil"
	"io"
	"net/http"
	"bytes"
	"strings"
	"errors"
)

const (
	timeout = time.Duration(30 * time.Second)
)
type HttpDTO struct{
	URL string
	Body io.Reader
	Header http.Header
}

func NewHttpClient(url string,body io.Reader,header http.Header) *HttpDTO {
	return &HttpDTO{URL:url,Body:body,Header:header}
}

func (dto HttpDTO)Get() (*Response,error) {
	client := &http.Client{Timeout: timeout}
	request, err := http.NewRequest("GET", dto.URL, dto.Body)
	if err != nil {
		log.Error("New request error.",err)
		return nil,err
	}
	request.Header = dto.Header
	//reqest.Header.Add("uid", uid)
	//reqest.Header.Add("User-Agent", "xxx")
	//reqest.Header.Add("X-Requested-With", "xxxx")
	//处理返回结果
	res, err := client.Do(request)
	if err != nil {
		log.Error("client.Do(request) error.",err)
		return nil,err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("response.Body to []byte error.",err)
		return nil,err
	}
	defer res.Body.Close()
	response := &Response{Status:res.StatusCode,Headers:res.Header,Body:NewJSONReader(body)}
	return response,nil
}

type Response struct {
	Headers map[string][]string
	Body    *JSONReader
	Status  int
}

type JSONReader struct {
	*bytes.Reader
}

func NewJSONReader(outBytes []byte) *JSONReader {
	jr := new(JSONReader)
	jr.Reader = bytes.NewReader(outBytes)
	return jr
}

func (js JSONReader) MarshalBody() (string, error) {
	data, err := ioutil.ReadAll(js.Reader)
	if err != nil {
		return "", err
	}
	return string(data), nil
}


func (js JSONReader) MarshalJSON() ([]byte, error) {
	data, err := ioutil.ReadAll(js.Reader)
	if err != nil {
		return nil, err
	}
	return data,nil
	//log.Info(string(data))
	//data = []byte(`"` + string(data) + `"`)
	//return data, nil
}

// UnmarshalJSON sets *jr to a copy of data.
func (jr *JSONReader) UnmarshalJSON(data []byte) error {
	if jr == nil {
		return errors.New("json.JSONReader: UnmarshalJSON on nil pointer")
	}
	if data == nil {
		return nil
	}
	data = []byte(strings.Trim(string(data), "\""))
	jr.Reader = bytes.NewReader(data)
	return nil
}