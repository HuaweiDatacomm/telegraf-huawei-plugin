package huawei_grpc_gpb

import (
    "fmt"
    "github.com/golang/protobuf/proto"
    "github.com/influxdata/telegraf/plugins/common/telemetry_proto"
    "github.com/influxdata/telegraf/plugins/common/telemetry_proto/huawei_ifm"
    "github.com/influxdata/telegraf/testutil"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "strconv"
    "testing"
)

func marshalTelemIfm(index uint64) ([]byte , error) {
    ifmByte, err := marshalIfm(index)
    //var testN uint32 = 23
    //debugByte, err := marshalDebug(testN)
    if err == nil {
        row := &telemetry.TelemetryRowGPB{
            //Content: debugByte,
            Content: ifmByte,
            Timestamp: uint64(1614160758),
            //ifmByte,
        }
        rows := []*telemetry.TelemetryRowGPB{row}
        gpbTable := telemetry.TelemetryGPBTable{
            Row:       rows,
            Delete:    nil,
            Generator: nil,
        }
        //telemetryGPBTable := &telemetry.TelemetryGPBTable{}
        //teleStructTest := &telemetry.Telemetry{
        //    SensorPath: "sensormine",
        //}
        teleStruct := &telemetry.Telemetry{
            NodeIdStr:           "nodeid",
            SubscriptionIdStr:   "xx",
            SensorPath:          "huawei-ifm:ifm/interfaces/interface",
            //SensorPath:          "huawei-debug:debug/interfaces/interface",
            ProtoPath:           "huawei_ifm.Ifm",
            CollectionId:        index,
            CollectionStartTime: 0,
            MsgTimestamp:        uint64(1614160758),
            DataGpb:             &gpbTable,
            CollectionEndTime:   0,
            CurrentPeriod:       0,
            ExceptDesc:          "desc",
            ProductName:         "pro",
            Encoding:            0,
            DataStr:             "dd",
            NeId:                "ss233333333333333333",
            SoftwareVersion:     "20.1",
        }

        return proto.Marshal(teleStruct)
    }
     return nil,nil
}

func marshalIfm (index uint64) ([]byte, error){
    ipv4ConflictEnable := huawei_ifm.Ifm_Global_Ipv4ConflictEnable{PreemptEnable: true}
    ipv6ConflictEnable := huawei_ifm.Ifm_Global_Ipv6ConflictEnable{PreemptEnable: false}
    fimIfmGlobal := huawei_ifm.Ifm_Global_FimIfmGlobal{
        TrunkDelaysendTime: 110,
    }
    global := huawei_ifm.Ifm_Global {
        StatisticInterval:    388,
        Ipv4IgnorePrimarySub: false,
        Ipv4ConflictEnable:   &ipv4ConflictEnable,
        Ipv6ConflictEnable:   &ipv6ConflictEnable,
        FimIfmGlobal:         &fimIfmGlobal,
        FimTrunkLocalfwd:     nil,
        VeGroups:             nil,
    }
    len := 2
    intefs := [2]*huawei_ifm.Ifm_Interfaces_Interface{}
    for j := 1; j<=len;j++ {
        interface2 := huawei_ifm.Ifm_Interfaces_Interface{
            Name:            "interface"+strconv.Itoa(int(index)),
            Index:           uint32(j),
            Position:        "huaweiggg2",
            ParentName:      "nanjinghuawei2",
            Number:          "2234",
            Description:     "woyonglaiceshi2",
            AggregationName: "",
            IsL2Switch:      false,
            MacAddress:      "45sddfs5wr12",
            VsName:          "vsss2",
        }
        intefs[j-1]=&interface2
    }
    interfaces := huawei_ifm.Ifm_Interfaces{Interface: intefs[0:]}
    ifm := &huawei_ifm.Ifm{
        Global:                &global,
        Interfaces:            &interfaces,
        Damp:                  nil,
        AutoRecoveryTimes:     nil,
        StaticDimensionRanges: nil,
        Ipv4InterfaceCount:    nil,
        RemoteInterfaces:      nil,
        HdlcDamp:              nil,
    }
    return proto.Marshal(ifm)
}

func TestDataPublish(t *testing.T) {
    telebyte, err := marshalTelemIfm(uint64(1))
    if err != nil {
        fmt.Println("ERR to connect: "+strconv.Itoa((1)))
    }
    tes := testutil.Logger{}
    parser := &Parser{Log:tes}
    metrics, err := parser.Parse([]byte(telebyte))
    require.NoError(t, err)
    assert.Equal(t, "huawei-ifm:ifm/interfaces/interface", metrics[0].Name())
}