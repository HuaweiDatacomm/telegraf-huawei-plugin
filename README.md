# **telegraf-huawei-plugin**

## **Overview**
the huawei plugin for telegraf to collect and process information from huawei devices

## **Installation**
### **Prerequisites**

- OS : Ubuntu, CentOS, Suse, Windows, Red Hat
- Go : go1.17.2
- Telegraf : Telegraf (1.20 recommended)
- Glibc : https://www.gnu.org/software/libc/libc.html
- Make : https://www.gnu.org/software/make/
- Protoc-gen-go :https://github.com/golang/protobuf




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
   ./install.sh
   ```
4. get the file of proto ,then use protoc-gen-go generate the file of proto , here is an example of huawei_debug.proto
   ```
   cd /telegraf/plugins/parsers/huawei_grpc_gpb/telemetry_proto
   mkdir huawei_debug (put huawei_debug.proto in this dir (this dir's name must be same of proto ))
   protoc --go_out=plugins=grpc:. huawei_debug.proto
   vim HuaweiTelemetry.go (
   add "github.com/influxdata/telegraf/plugins/parsers/huawei_grpc_gpb/telemetry_proto/huawei_debug" in import
   add  PathKey{ProtoPath: "huawei_debug.Debug", Version: "1.0"}: []reflect.Type{reflect.TypeOf((*huawei_debug.Debug)(nil))},
   in the last function ("var pathTypeMap = map[PathKey][]reflect.Type{}"),
   like this :
   var pathTypeMap = map[PathKey][]reflect.Type{
      PathKey{ProtoPath: "huawei_debug.Debug", Version: "1.0"}: []reflect.Type{reflect.TypeOf((*huawei_debug.Debug)(nil))}, 
   }
   )
   
   This step needs to be improved. 
   ```
5. Run `make` from the source directory,you can see telegraf-huawei-plugin
   ```
   cd telegraf
   make
   telegraf --input-list
   ```
## Getting Used
  
 - The TIG(telegraf,influxdb,grafana) is an open-source O&M tool that collects Telemetry data sent by devices, analyzes the data, and displays the data graphically.
   The other two tools can be downloaded from the official website
 - 1.configuration telegraf.conf (telegraf/ect/telegraf.conf)
   ```
   [[outsputs.influxdb]]
   urls = ["http://127.0.0.1:8086"]
   database = ""
   
   [[inputs.huawei_telemetry_dialout]]
   service_address ="ip:port"
   data_format = "grpc"
   transport = "grpc"

   [[inputs.huawei_telemetry_dialin]]
   data_format = "grpc" 
   [[inputs.huawei_telemetry_dialin.routers]]
   address = "ip:port"
   sample_interval = 
   encoding="json"  # or "gpb" 
   request_id = 
     [inputs.huawei_telemetry_dialin.routers.aaa]
        username = ""
        password = ""
     [[inputs.huawei_telemetry_dialin.routers.Paths]]
        depth = 1
        path = ""

   [processors.metric_match.approach]
   appproach = "include" # or exclude
   [processors.metric_match.tag]
   "telemetry" = [""]
   [processors.metric_match.field_filter]
   ""=[""]
   ```
 - 2.start influxdb
   ```
   cd influxdb/usr/bin
   ./influxd
   ```
 - 3.start telegraf
 - 4.use grafana tool

  







