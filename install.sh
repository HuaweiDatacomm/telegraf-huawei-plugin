#!/bin/sh
echo "please enter telegraf path: (please use '/' replace '\'):" 
read telegraf_dir
huaweiplugin_dir=$(pwd)
echo $telegraf_dir 
echo $huaweiplugin_dir
cp -r $huaweiplugin_dir/huawei-plugin $telegraf_dir

cp -r $telegraf_dir/huawei-plugin/huawei_telemetry_dialin $telegraf_dir/telegraf/plugins/inputs
cp -r $telegraf_dir/huawei-plugin/huawei_telemetry_dialout $telegraf_dir/telegraf/plugins/inputs
cp -r $telegraf_dir/huawei-plugin/huawei_grpc_json $telegraf_dir/telegraf/plugins/parsers
cp -r $telegraf_dir/huawei-plugin/huawei_grpc_gpb $telegraf_dir/telegraf/plugins/parsers
cp -r $telegraf_dir/huawei-plugin/metric_match $telegraf_dir/telegraf/plugins/processors
str1='_ "github.com/influxdata/telegraf/plugins/inputs/huawei_telemetry_dialin"'
str2='_ "github.com/influxdata/telegraf/plugins/inputs/huawei_telemetry_dialout"'
str3='"github.com/influxdata/telegraf/plugins/parsers/huawei_grpc_json"'
str4='"github.com/influxdata/telegraf/plugins/parsers/huawei_grpc_gpb"'
str5='_ "github.com/influxdata/telegraf/plugins/processors/metric_match"'
sed -i "5i ${str1}" $telegraf_dir/telegraf/plugins/inputs/all/all.go
sed -i "5i ${str2}" $telegraf_dir/telegraf/plugins/inputs/all/all.go
sed -i "5i ${str3}" $telegraf_dir/telegraf/plugins/parsers/registry.go
sed -i "5i ${str4}" $telegraf_dir/telegraf/plugins/parsers/registry.go
sed -i "5i ${str5}" $telegraf_dir/telegraf/plugins/processors/all/all.go

cat>> $telegraf_dir/telegraf/plugins/parsers/registry.go <<EOF


func NewHuaweiGrpcGpbParser() (Parser, error) {
	tags := map[string]string{"parsers": "huawei_grpc_gpb_parser"}
	grpcRegister := selfstat.Register("huawei_grpc_gpb_parser", "errors", tags)
	logger := models.NewLogger("parsers", "huawei_grpc_gpb_parser", "")
	logger.OnErr(func() {
		grpcRegister.Incr(1)
	})
	parser, err := huawei_grpc_gpb.New()
	if err != nil {
		return nil, err
	}
	models.SetLoggerOnPlugin(parser, logger)
	return parser, err
}

func NewHuaweiGrpcJsonParser() (Parser, error) {
	tags := map[string]string{"parsers": "huawei_grpc_json_parser"}
	grpcRegister := selfstat.Register("huawei_grpc_json_parser", "errors", tags)
	logger := models.NewLogger("parsers", "huawei_grpc_json_parser", "")
	logger.OnErr(func() {grpcRegister.Incr(1)})
	parser, err := huawei_grpc_json.New()
	models.SetLoggerOnPlugin(parser, logger)
	return parser, err
}
EOF

echo "install huawei-plugin successfully"

