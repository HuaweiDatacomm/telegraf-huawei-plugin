package grpc

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/models"
	"github.com/influxdata/telegraf/plugins/processors"
	"github.com/influxdata/telegraf/selfstat"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const sampleConfig = ``
const sensorPathKey = "sensor_path"
const TelemetryKey = "telemetry"
const pointConfig = "."

type MetricMatch struct {
	Tag         map[string][]string `toml:"tag"`
	FieldFilter map[string][]string `toml:"field_filter"`
	Approach    map[string]string   `toml:"approach"`
	Log         telegraf.Logger
}

func (m *MetricMatch) SampleConfig() string {
	return sampleConfig
}

func (m *MetricMatch) Description() string {
	return "metric match"
}

func (m *MetricMatch) Apply(in ...telegraf.Metric) []telegraf.Metric {

	// get telemetry header field_filter and tag
	//var useWay string
	var res []telegraf.Metric
	headerFilter := m.FieldFilter[TelemetryKey]
	headerTag := m.Tag[TelemetryKey]
	headerWay := m.Approach
	var approach string
	for _, v := range headerWay {
		approach = v
	}

	if approach == "include" {

		// include filter field

		for _, eachMetric := range in {
			sensorPath, ok := eachMetric.GetField(sensorPathKey)
			if ok {
				fieldFilters := m.FieldFilter[sensorPath.(string)]
				if len(fieldFilters) == 0 {
					m.Log.Warnf("the %s 's field filters is empty...", sensorPath)
				}
				allKeys := make([]string, 0)
				needKeys := make([]string, 0)
				for _, v := range eachMetric.FieldList() {
					if !strings.Contains(v.Key, pointConfig) {
						needKeys = append(needKeys, v.Key)
					}
					allKeys = append(allKeys, v.Key)
				}

				fieldFilters = append(fieldFilters, headerFilter...)
				for _, filter := range fieldFilters {

					if ok, matchKeys := matchField(filter, eachMetric.FieldList()); ok {

						for _, realKey := range matchKeys {
							needKeys = append(needKeys, realKey)
							//eachMetric.AddField(realKey,eachMetric)
						}
					}
				}

				for _, needKeysV := range needKeys {
					for k, allKeysV := range allKeys {
						if allKeysV == needKeysV {
							allKeys = append(allKeys[:k], allKeys[k+1:]...)
						}
					}
				}

				for _, v := range allKeys {
					eachMetric.RemoveField(v)
				}

			}
		}

	} else {
		// exclude filter field
		for _, eachMetric := range in {
			sensorPath, ok := eachMetric.GetField(sensorPathKey)
			if ok {
				fieldFilters := m.FieldFilter[sensorPath.(string)]
				if len(fieldFilters) == 0 {
					m.Log.Warnf("the %s 's field filters is empty...", sensorPath)
				}
				fieldFilters = append(fieldFilters, headerFilter...)
				for _, filter := range fieldFilters {
					if ok, matchKeys := matchField(filter, eachMetric.FieldList()); ok {
						for _, realKey := range matchKeys {
							eachMetric.RemoveField(realKey)
						}
					}
				}
			}
		}

	}

	// field to tag
	for _, eachMetric := range in {
		sensorPath, ok := eachMetric.GetField(sensorPathKey)
		if ok {
			tags := m.Tag[sensorPath.(string)]
			if len(tags) == 0 {
				m.Log.Warnf("the %s 's tag is empty...", sensorPath)
			}
			tags = append(tags, headerTag...)
			for _, tag := range tags {
				if ok, matchKeys := matchField(tag, eachMetric.FieldList()); ok {
					for _, realKey := range matchKeys {
						value, ok := eachMetric.GetField(realKey)
						if ok {
							typeOfV := reflect.TypeOf(value)
							if typeOfV.Name() != "string" {
								if typeOfV.Name() != "int64" {
									m.Log.Errorf("wrong with metric tag [%s %s], it's type is %s", sensorPath.(string), tag, typeOfV.Name())
									m.stop()
								} else {
									value = strconv.FormatInt(value.(int64), 10)
								}
							}
							eachMetric.AddTag(realKey, value.(string))
							eachMetric.RemoveField(tag)
						}
					}
				}
			}
		}
		res = append(res, eachMetric)
	}

	return res
}

func matchField(key string, fields []*telegraf.Field) (bool, []string) {
	var matches []string
	for _, field := range fields {
		//if strings.HasSuffix(field.Key, key) {
		//   matches = append(matches, field.Key)
		//}
		if field != nil {
			m := regexp.MustCompile(".*[^.]" + key + "$")
			ok := m.MatchString(string(field.Key))
			if ok {
				matches = append(matches, field.Key)
			}
		}

	}
	if len(matches) > 0 {
		return true, matches
	} else {
		return false, matches
	}
}

func init() {
	processors.Add("metric_match", func() telegraf.Processor {
		tags := map[string]string{"processor": "metric_match"}
		grpcRegister := selfstat.Register("metric_match", "errors", tags)
		logger := models.NewLogger("processors", "metric_match", "")
		logger.OnErr(func() {
			grpcRegister.Incr(1)
		})
		return &MetricMatch{
			Log: logger,
		}
	})
}

func (c *MetricMatch) stop() {
	log.SetOutput(os.Stderr)
	log.Printf("I! telegraf stopped because error.")
	os.Exit(1)
}
