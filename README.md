# gosoap
http => soap; soap => http

### 配置说明  
appname = easemob-gosoap #服务名称  
httpport = 8080 #监听端口  
runmode = dev #运行模式  
copyrequestbody = true #默认配置  
log_level = debug #日志级别  
log_path = /Users/Joey_Wu/go/src/easemob-gosoap/log/easemob-gosoap.log #日志目录及日志文件名称  

#判断服务请求从何而来  
#true表示从请求协议转换至soap协议，发送给webservice  
#false表示从请求协议转换至http协议，发送给后端http服务  
request_from_wsdl = "false"  
  
#当request_from_wsdl = "false"时，请求协议转换至soap并转发至以下webservice  
wsdl_server = http://aaa.bbb.ccc.ddd:13201/vss/2.0/workflow/CSGServiceServiceSync?wsdl  
#当request_from_wsdl = "true"时，请求协议转换http并转发至以下http server  
http_server = http://abc.abc.com  

#webservice平台认证token  
token = xxxxxxxxxxxxxxxxxxxxxxxxxx  
