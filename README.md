# **telegraf-huawei-plugin**

## **Overview**
the huawei plugin for telegraf to collect and process information from huawei devices

## **Installation**
### **Prerequisites**

- OS : Ubuntu, CentOS, Suse, Windows, Red Hat
- Go : go1.17.2
- Telegraf : Telegraf (1.20 recommended)
- protoc :  3.11.4
  https://github.com/protocolbuffers/protobuf/releases
- protoc-gen-go :
  go get -u github.com/protobuf/protoc-gen-go@v1.2.0

### Build From Source

Telegraf requires Go version 1.17.2 or newer, the Makefile requires GNU make.


1. Clone the Telegraf and telegraf-huawei-plugin repository（）:
   ```
   git clone https://github.com/influxdata/telegraf.git
   git clone https://github.com/HuaweiDatacomm/telegraf-huawei-plugin.git
   ```
2. Configuring the environment of telegraf ,here's an example 
   ```
   vim /etc/profile  
   export TELEGRAFROOT = the dir of telegraf
   source /etc/profile
   ```
3. Run install.sh
   ```
   chmod +x install.sh
   ./install.sh
   ```
4. get the file of proto ,then use protoc-gen-go generate the file of proto , here is an example of huawei_debug.proto . 
   ```
   cd /telegraf/plugins/parsers/huawei_grpc_gpb/telemetry_proto
   mkdir huawei_debug (put huawei-debug.proto in this dir (Note:the dir's name has "_",not "-"))
   cd huawei_debug
   protoc --go_out=plugins=grpc:. huawei-debug.proto
   vim HuaweiTelemetry.go (
   add "github.com/influxdata/telegraf/plugins/parsers/huawei_grpc_gpb/telemetry_proto/huawei_debug" in import
   add  PathKey{ProtoPath: "huawei_debug.Debug", Version: "1.0"}: []reflect.Type{reflect.TypeOf((*huawei_debug.Debug)(nil))},
   in the last function ("var pathTypeMap = map[PathKey][]reflect.Type{}"),
   like this :
   var pathTypeMap = map[PathKey][]reflect.Type{
      PathKey{ProtoPath: "huawei_debug.Debug", Version: "1.0"}: []reflect.Type{reflect.TypeOf((*huawei_debug.Debug)(nil))}, 
   }
   )
   
   ```
5. Run `make` from the source directory,you can see telegraf-plugins are installed
   ```
   cd telegraf
   make
   telegraf --input-list
   ```
## Getting Used
  
 - The TIG(Telegraf,Influxdb,Grafana) is an open-source O&M tool that collects Telemetry data sent by devices, analyzes the data, and displays the data graphically.
   The other two tools can be downloaded from the official website
 - Influxdb:https://dl.influxdata.com/influxdb/releases/influxdb-1.8.7_linux_amd64.tar.gz
 - Grafana:https://dl.grafana.com/oss/release/grafana-7.3.6.linux-amd64.tar.gz
 - 1.config huawei devices  
   https://support.huawei.com/enterprise/en/doc/EDOC1100055030/b650f7a7  
 - 2.copy telegraf.conf to /etc/telegraf
   ```
   cd /etc
   mkdir telegraf
   cp (the dir of telegraf)/etc/telegraf.conf /etc/telegraf
   vim /etc/telegraf/telegraf.conf
   ```
 - 3.configuration telegraf.conf (telegraf/ect/telegraf.conf)
   ```
   # ##################################output plugin influxdb##################################
   #[[outsputs.influxdb]]
   #urls = ["http://127.0.0.1:8086"]
   #database = "DTS"
   # username = ""
   # password = ""
   
   # ###################################input plugin telegraf-huawei###########################
   # ##telegraf-huawei is divided into two: huawei_telemetry_dialout and huawei_telemetry_dialin
   # ##huawei_telemetry_dialout
   # ##service_address = “ip:port”	 #IP address and port number enabled for static subscription
   # ##data_format = "grpc"  #reserved field. The value is fixed at grpc.
   # ##transport = "grpc"   #reserved field. The value is fixed at grpc.
  
   #[[inputs.huawei_telemetry_dialout]]
   #service_address ="ip:port"
   #data_format = "grpc"
   #transport = "grpc"

   # ##huawei_telegraf_dialin
   # ##data_format  #reserved field. The value is fixed at grpc.
   # ##address  #collection IPaddress and port
   # ##sample_interval  #sampling interval
   # ##encoding  #coding rules "json" or "gpb"
   # ##request_ip  #subscription parameters, vaule is 1~2^64-1
   # ##suppress_redundant #option of suppress_redundant: true or false
   # ##username # aaa's username
   # ##password # aaa's password(encrypted)
   # ##path #collection path
   # ##metric_match.approach #filtering metrics methods ,"include" or "exclude"
   # ##metric_match.field_filter #filtering metrics matching 
   # ##metric_match.tag # field to tag
   
   #[[inputs.huawei_telemetry_dialin]]
   #data_format = "grpc" 
   #[[inputs.huawei_telemetry_dialin.routers]]
   #address = "ip:port"
   #sample_interval = 1000
   #encoding="json"  # or "gpb" 
   #suppress_redundant = true
   #request_id = 
   #  [inputs.huawei_telemetry_dialin.routers.aaa]
   #     username = ""
   #     password = ""
   #  [[inputs.huawei_telemetry_dialin.routers.Paths]]
   #     depth = 1
   #     path = ""

   #[processors.metric_match.approach]
   #appproach = "include" # or exclude
   #[processors.metric_match.tag]
   #"telemetry" = [""]
   #[processors.metric_match.field_filter]
   #"path"=[""]
   ```
 - 4.start influxdb
   ```
   cd influxdb/usr/bin
   ./influxd
   ./influx
   create database DTS 
   ```
 - 5.start telegraf
   ```
   cd (the dir of telegraf)/cmd/telegraf
   go run ./
   ```
 - 6.use Grafana  
   1.Login  
   2.Configuration  
   3.Add data source  
   4.Choose influxdb  
   5.Fill in influxdb's configuraton,such as url,database,username,password  
   6.Save & Test  
   7.Create Dashboard  
   8.Add new panel  
 
   

  







