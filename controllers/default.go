package controllers

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)




type MainController struct {
	beego.Controller
}


func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "aabbcc@gmail.com"
	c.TplName = "index.tpl"
}


// 注册ProxyController
type ProxyApiController struct {
	beego.Controller
}

func (c *ProxyApiController) AllMethod() {
	request_from_wsdl := beego.AppConfig.String("request_from_wsdl")
	beego.Info(request_from_wsdl)
	// 判断请求来源及处理逻辑
	// 如果请求来源为wsdl，则采用将soap协议转换为正常http协议逻辑
	// 如果请求来源非wsdl，则将请求封装为soap协议，并转发至wsdl
	if  request_from_wsdl == "false" {
		// Get WSDL TOKEN from app.conf
		token := beego.AppConfig.String("token")
		beego.Info("WSDL token is", token)
		// Get Request Method for string
		Method := c.Ctx.Request.Method
		beego.Info("request method is", Method)
		// Get Request Body for string
		request_body := string([]byte(c.Ctx.Input.RequestBody))
		beego.Info("request body is", request_body)
		// Get Request Header for string
		// init headers map
		headers := make(map[string]string)
		// Get header's k,v from Request
		for k, v := range c.Ctx.Request.Header{
			// apend header to map headers
			headers[k] = v[0]
		}
		// from map to string
		headers_str, err := json.Marshal(headers)
		if err != nil {
			fmt.Println(err)
		}
		beego.Info("request headers is", string(headers_str))
		// Get Request Uri
		//request_domain := c.Ctx.Input.Domain()
		//request_port := strconv.Itoa(c.Ctx.Input.Port())
		request_uri := c.Ctx.Request.RequestURI
		beego.Info("request uri is", request_uri)
		// Make wsdl Request data
		req_data := make(map[string]string)
		req_data["headers"] = string(headers_str)
		req_data["method"] = Method
		req_data["body"] = request_body
		req_data["request_uri"] = request_uri
		req_data_str, err := json.Marshal(req_data)
		// encode Request data
		encode_req_data_str := EncodingData(req_data_str)
		// Parse Xml data
		wsdl_data := Parserequestxml(encode_req_data_str, token)
		//fmt.Println(wsdl_data)
		// Request wsdl with pared xml data
		wsdl_resq_code, wsdl_resq_body := RequestWSDL(wsdl_data)
		//c.Ctx.WriteString(strconv.Itoa(wsdl_resq_code))
		//c.Ctx.WriteString(wsdl_resq_body)
		if wsdl_resq_code == 200 {
			resq_reslut := Decoderesponsexml([]byte(wsdl_resq_body))
			resq_data := resq_reslut.SoapenvBody.EbsCsgMessage.Return.Message
			resq_code := resq_reslut.SoapenvBody.EbsCsgMessage.Return.Code
			if resq_code == "00" {
				beego.Info("response from wsdl is: ", resq_data)
				c.Ctx.WriteString(resq_data)
			} else {
				beego.Error("respose from wsdl is: ", resq_data)
				c.Ctx.WriteString(resq_data)
			}
		} else {
			c.Ctx.ResponseWriter.WriteHeader(500)
		}
	} else {
		// Get Request Body for string
		request_body := string([]byte(c.Ctx.Input.RequestBody))
		beego.Info("get request from wsdl")
		beego.Info("requset body from wsdl is", request_body)
		// Decode xml
		wsdl_data := Decoderequestxml([]byte(request_body))
		// Get request data
		wsdl_req_data := wsdl_data.SoapenvBody.EbsCsgMessage.Data
		//c.Ctx.WriteString(request_body)
		//c.Ctx.WriteString(wsdl_req_data)
		beego.Info("data from wsdl is: ", wsdl_req_data)
		// Decode data by base64
		response_code, response_body := RequestHTTP(DecodingData([]byte(wsdl_req_data)))
		if response_code == "00" {
			wsdl_resp_data := Parseresponsexml(response_code, response_body)
			beego.Info("response code:", response_code)
			beego.Info("response body:", response_body)
			c.Ctx.WriteString(wsdl_resp_data)
		} else {
			wsdl_resp_data := Parseresponsexml(response_code, response_body)
			beego.Error("response code:", response_code)
			beego.Error("response body:", response_body)
			c.Ctx.WriteString(wsdl_resp_data)
		}
	}
}

// Encoding data by base64
func EncodingData(data []byte) (string){
	result := base64.StdEncoding.EncodeToString(data)
	return result
}

// Decoding data by base64
func DecodingData(data []byte) (string){
	result, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		beego.Error("can not decode to string,", data)
	}
	return string(result)
}

// Parse xml for request to wsdl server
func Parserequestxml(data string, token string) (string){
	result := `<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ebs="http://ebs.soa.csg.cn">
   <soapenv:Header/>
   <soapenv:Body>
       <ebs:csg_message>
           <data>REQ_DATA_STR</data>
           <token>TOKEN</token>
       </ebs:csg_message>
   </soapenv:Body>
</soapenv:Envelope>`
	result = strings.Replace(result, "REQ_DATA_STR", string(data), -1)
	result = strings.Replace(result, "TOKEN", token, -1)
	return result
}

// Parse xml for response to wsdl server
func Parseresponsexml(code string, message string) (string){
	current_time := time.Now().Format("20060102150405")
	fmt.Println(current_time)
	result := `<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ebs="http://ebs.soa.csg.cn">
    <soapenv:Header/>
    <soapenv:Body>
        <ebs:csg_message>
            <return>
                <code>RESP_CODE</code>
                <message>RESP_MESSAGE</message>
                <time>RESP_TIME</time>
            </return>
        </ebs:csg_message>
    </soapenv:Body>
</soapenv:Envelope>`
	result = strings.Replace(result, "RESP_CODE", string(code), -1)
	result = strings.Replace(result,"RESP_MESSAGE", string(message), -1)
	result = strings.Replace(result, "RESP_TIME", string(current_time), -1)
	return result
}


// Decode xml from request
type XReqEbsCsgMessage struct {
	Data string `xml:"data"`
	Token string `xml:"token"`
}
type XReqSoapenvBody struct {
	EbsCsgMessage XReqEbsCsgMessage `xml:"csg_message"`
}
type ReqResult struct {
	XMLName xml.Name `xml:"Envelope"`
	SoapenvHeader string `xml:"header"`
	SoapenvBody XReqSoapenvBody `xml:"Body"`
}
func Decoderequestxml(data []byte) ReqResult{
	var result ReqResult
	if err := xml.Unmarshal(data, &result);err == nil {
		beego.Info("decoding xml from wsdl server")
		return result
	} else {
		beego.Error(err)
		panic(err)
	}
}


// Decode xml from response
type XRespEbsCsgMessage struct {
	Code string `xml:"code"`
	Message string `xml:"message"`
}
type XRespReturn struct {
	Return XRespEbsCsgMessage `xml:"return"`
}
type XRespSoapenvBody struct {
	EbsCsgMessage XRespReturn `xml:"csg_message"`
}
type RespResult struct {
	XMLName xml.Name `xml:"Envelope"`
	SoapenvBody XRespSoapenvBody `xml:"Body"`
}
func Decoderesponsexml(data []byte) RespResult{
	var result RespResult
	if err := xml.Unmarshal(data, &result);err == nil {
		beego.Info("decoding xml from http server")
		return result
	} else {
		beego.Error(err)
		panic(err)
	}
}


// Request WSDL with parsed xml data
func RequestWSDL(data string) (http_code int, body string){
	wsdl_server := beego.AppConfig.String("wsdl_server")
	beego.Info("request to", wsdl_server)
	beego.Info("request to wsdl server with data\n", data)
	req, err := http.NewRequest("POST", wsdl_server, strings.NewReader(data))
	if err != nil {
		fmt.Println(err)
	}
	defer req.Body.Close()
	client := &http.Client{Timeout: 3 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		beego.Error("Failed request to wsdl server")
		return 500, "Failed request to wsdl server"
	}
	defer resp.Body.Close()
	http_code = resp.StatusCode
	result, _ := ioutil.ReadAll(resp.Body)
	body = string(result)
	return http_code, body
}

// Request http with decoded xml data
func RequestHTTP(data string) (wsdl_code string, body string){
	type Request struct {
		Body string `json:"body"`
		Headers string `json:"headers"`
		Method string `json:"method"`
		RequestURI string `json:"request_uri"`
	}
	var request Request
	if err := json.Unmarshal([]byte(data), &request); err == nil{
		fmt.Println(request.Method)
		fmt.Println(request.Body)
		fmt.Println(request.Headers)
		fmt.Println(request.RequestURI)
	} else {
		fmt.Println(err)
	}
	http_server := beego.AppConfig.String("http_server")
	beego.Info("request to", http_server)
	beego.Info("request to http server with data\n", data)
	// GET request
	if request.Method == "GET" {
		requesturl := http_server + request.RequestURI
		req, err := http.NewRequest(request.Method, requesturl, strings.NewReader(request.Body))
		if err != nil {
			beego.Error(err)
		}
		headers := make(map[string]string)
		herr := json.Unmarshal([]byte(request.Headers), &headers)
		if herr != nil {
			beego.Error(herr)
		}
		for k, v := range headers {
			req.Header.Add(k, v)
		}
		defer req.Body.Close()
		client := &http.Client{Timeout: 3 * time.Second}
		resp, error := client.Do(req)
		if error != nil {
			beego.Error("Failed request to http server")
			return "99","Failed request to http server"
		}
		defer resp.Body.Close()
		result, _ := ioutil.ReadAll(resp.Body)
		return "00", string(result)
	} else {
		requesturl := http_server + request.RequestURI
		req, err := http.NewRequest(request.Method, requesturl, strings.NewReader(request.Body))
		if err != nil {
			beego.Error(err)
		}
		headers := make(map[string]string)
		herr := json.Unmarshal([]byte(request.Headers), &headers)
		if herr != nil {
			beego.Error(herr)
		}
		for k, v := range headers {
			req.Header.Add(k, v)
		}
		defer req.Body.Close()
		client := &http.Client{Timeout: 3 * time.Second}
		resp, error := client.Do(req)
		if error != nil {
			beego.Error("Failed request to http server")
			return "99","Failed request to http server"
		}
		defer resp.Body.Close()
		result, _ := ioutil.ReadAll(resp.Body)
		return "00", string(result)
	}
}