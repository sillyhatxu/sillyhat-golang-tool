package sillyhat_utils

import "encoding/json"

func ToJson(obj interface{}) (string,error) {
	result,err := json.Marshal(obj)
	if err != nil {
		return "",err
	}
	return string(result),nil

}
