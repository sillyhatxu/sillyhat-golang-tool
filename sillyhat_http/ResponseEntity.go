package sillyhat_http

type ResponseBody struct {
	Ret int `json:"ret"`
	Data interface{} `json:"data"`
	Msg string `json:"msg"`
}

func Success(data interface{}) *ResponseBody {
	return &ResponseBody{Ret:RESPONSE_SUCCESS,Data:data,Msg:"success"}
}

func ErrorData(ret int,data interface{},msg string) *ResponseBody {
	return &ResponseBody{Ret:ret,Data:data,Msg:msg}
}

func Error(ret int,msg string) *ResponseBody {
	return &ResponseBody{Ret:ret,Msg:msg}
}