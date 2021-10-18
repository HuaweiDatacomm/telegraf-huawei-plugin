package huawei_telemetry_dialin

import (
    "context"
    huawei_dialin"github.com/influxdata/telegraf/plugins/inputs/huawei_telemetry_dialin/huawei_dialin"
    "log"
    "net"
    "testing"

    telemetry "github.com/cisco-ie/nx-telemetry-proto/telemetry_bis"
    "github.com/golang/protobuf/proto"
    "github.com/influxdata/telegraf/testutil"
    "github.com/stretchr/testify/require"
    "google.golang.org/grpc"
)

func (dialin *HuaweiTelemetryDialin) Subscribe(stream huawei_dialin.GRPCConfigOper_SubscribeServer) error {
    return nil
}

func TestHuaweiTelemetryDialin_Start(t *testing.T) {
    // 1. prepare telemetry struct data, and marshal to []byte
    // 2. start one grpc server with huawei-grpc-dialin.proto in go routine
    //
    // 3. one go routine opened to run Start() in huawei_telemetry_dialin.go
    // 4.
    // start the mock dialin server
    // initialize test object
    hin := &HuaweiTelemetryDialin{Log: testutil.Logger{}, Transport: "grpc"}
    acc := &testutil.Accumulator{}
    // listen on address with transport
    transport := "grpc"
    address := "127.0.0.1:5431"
    listener, err := net.Listen(transport, address)
    if err != nil {
        log.Fatalf("net.Listen err: %v", err)
    }
    grpcServer := grpc.NewServer()
    // registry our service on grpc server
    huawei_dialin.RegisterGRPCConfigOperServer(grpcServer, hin)
    // Use the server Serve() method and our port information area to block and wait until the process is killed or Stop() is called
    err = grpcServer.Serve(listener)
    if err != nil {
        log.Fatalf("grpcServer.Serve err: %v", err)
    }

    // test the dialin client by testing Start
    errStart := hin.Start(acc)
    // error is expected since we are passing in dummy transport
    require.Error(t, errStart)

    telemetry := &telemetry.Telemetry{
        MsgTimestamp: 1543236572000,
        EncodingPath: "type:model/some/path",
        NodeId:       &telemetry.Telemetry_NodeIdStr{NodeIdStr: "hostname"},
        Subscription: &telemetry.Telemetry_SubscriptionIdStr{SubscriptionIdStr: "subscription"},
        DataGpbkv: []*telemetry.TelemetryField{
            {
                Fields: []*telemetry.TelemetryField{
                    {
                        Name: "keys",
                        Fields: []*telemetry.TelemetryField{
                            {
                                Name:        "name",
                                ValueByType: &telemetry.TelemetryField_StringValue{StringValue: "str"},
                            },
                            {
                                Name:        "uint64",
                                ValueByType: &telemetry.TelemetryField_Uint64Value{Uint64Value: 1234},
                            },
                        },
                    },
                    {
                        Name: "content",
                        Fields: []*telemetry.TelemetryField{
                            {
                                Name:        "bool",
                                ValueByType: &telemetry.TelemetryField_BoolValue{BoolValue: true},
                            },
                        },
                    },
                },
            },
            {
                Fields: []*telemetry.TelemetryField{
                    {
                        Name: "keys",
                        Fields: []*telemetry.TelemetryField{
                            {
                                Name:        "name",
                                ValueByType: &telemetry.TelemetryField_StringValue{StringValue: "str2"},
                            },
                        },
                    },
                    {
                        Name: "content",
                        Fields: []*telemetry.TelemetryField{
                            {
                                Name:        "bool",
                                ValueByType: &telemetry.TelemetryField_BoolValue{BoolValue: false},
                            },
                        },
                    },
                },
            },
        },
    }
    data, _ := proto.Marshal(telemetry)

    hin.handleTelemetry(data)
    require.Empty(t, acc.Errors)

    tags := map[string]string{"path": "type:model/some/path", "name": "str", "uint64": "1234", "source": "hostname", "subscription": "subscription"}
    fields := map[string]interface{}{"bool": true}
    acc.AssertContainsTaggedFields(t, "alias", fields, tags)

    tags = map[string]string{"path": "type:model/some/path", "name": "str2", "source": "hostname", "subscription": "subscription"}
    fields = map[string]interface{}{"bool": false}
    acc.AssertContainsTaggedFields(t, "alias", fields, tags)
}

func mockTelemetryMessage() *telemetry.Telemetry {
    return &telemetry.Telemetry{
        MsgTimestamp: 1543236572000,
        EncodingPath: "type:model/some/path",
        NodeId:       &telemetry.Telemetry_NodeIdStr{NodeIdStr: "hostname"},
        Subscription: &telemetry.Telemetry_SubscriptionIdStr{SubscriptionIdStr: "subscription"},
        DataGpbkv: []*telemetry.TelemetryField{
            {
                Fields: []*telemetry.TelemetryField{
                    {
                        Name: "keys",
                        Fields: []*telemetry.TelemetryField{
                            {
                                Name:        "name",
                                ValueByType: &telemetry.TelemetryField_StringValue{StringValue: "str"},
                            },
                        },
                    },
                    {
                        Name: "content",
                        Fields: []*telemetry.TelemetryField{
                            {
                                Name:        "value",
                                ValueByType: &telemetry.TelemetryField_Sint64Value{Sint64Value: -1},
                            },
                        },
                    },
                },
            },
        },
    }
}

func (dialin *HuaweiTelemetryDialin) Subscribe(args *huawei_dialin.SubsArgs, server huawei_dialin.GRPCConfigOper_SubscribeServer) error {
    panic("implement me")
}

func (dialin *HuaweiTelemetryDialin) Cancel(ctx context.Context, args *huawei_dialin.CancelArgs) (*huawei_dialin.CancelReply, error) {
    panic("implement me")
}
