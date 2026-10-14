package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/kernel/tsdb"
	"github.com/roamdream/roamadmin/modules/device/internal/consts"
	"github.com/roamdream/roamadmin/modules/device/internal/model"
)

// Ingest 消化一条 MQTT 上行消息（由消息总线回调，无 HTTP/租户上下文）。
//
// 主题约定：ra/{productKey}/{deviceKey}/telemetry 或 ra/{productKey}/{deviceKey}/shadow/reported。
// 未知设备 / 下行主题 / 非法负载视为可忽略情况（返回 nil）；真正的落库失败才返回 error，
// 供消息总线据此转死信（redis-streams 模式）。
func Ingest(ctx context.Context, topic string, payload []byte) error {
	productKey, deviceKey, kind, ok := parseTopic(topic)
	if !ok {
		return nil
	}
	dev, err := GetDeviceByKey(ctx, productKey, deviceKey)
	if err != nil {
		g.Log().Warningf(ctx, "[device] ingest ignore topic %s: device %s/%s not found", topic, productKey, deviceKey)
		return nil
	}
	var fields map[string]any
	if err := json.Unmarshal(payload, &fields); err != nil || fields == nil {
		g.Log().Warningf(ctx, "[device] ingest ignore topic %s: bad payload", topic)
		return nil
	}
	switch kind {
	case consts.TopicTelemetry:
		if err := reportTelemetry(ctx, dev, fields); err != nil {
			g.Log().Errorf(ctx, "[device] ingest telemetry failed: %v", err)
			return err
		}
	case consts.TopicShadowReported:
		if err := upsertReported(ctx, dev.Id, fields); err != nil {
			g.Log().Errorf(ctx, "[device] ingest shadow reported failed: %v", err)
			return err
		}
	default:
		// command / shadow/desired 为下行主题，忽略。
	}
	return nil
}

// ReportTelemetry HTTP 模拟上报（租户校验后复用 reportTelemetry）。
func ReportTelemetry(ctx context.Context, deviceID int64, fields map[string]any) error {
	dev, err := GetDevice(ctx, deviceID)
	if err != nil {
		return err
	}
	return reportTelemetry(ctx, dev, fields)
}

// reportTelemetry 写入遥测数据点并更新影子 reported。
func reportTelemetry(ctx context.Context, dev *model.Device, fields map[string]any) error {
	point := &tsdb.Point{
		Measurement: measurementOf(dev.Id),
		Tags:        map[string]string{"productKey": dev.ProductKey, "deviceKey": dev.DeviceKey},
		Fields:      fields,
		Time:        time.Now().UTC(),
	}
	if err := ra.TSDB().Write(ctx, point); err != nil {
		return err
	}
	return upsertReported(ctx, dev.Id, fields)
}

// QueryTelemetry 查询设备最近遥测数据点（时间倒序）。
func QueryTelemetry(ctx context.Context, deviceID int64, limit int) ([]*tsdb.Point, error) {
	dev, err := GetDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	q, ok := ra.TSDB().(tsdb.RecentQuerier)
	if !ok {
		return []*tsdb.Point{}, nil
	}
	return q.QueryRecent(ctx, measurementOf(dev.Id), limit)
}

// PublishCommand 向设备下发命令（MQTT 下行），返回主题与负载。
func PublishCommand(ctx context.Context, deviceID int64, data map[string]any) (topic, payload string, err error) {
	dev, err := GetDevice(ctx, deviceID)
	if err != nil {
		return "", "", err
	}
	topic = strings.Join([]string{consts.TopicRoot, dev.ProductKey, dev.DeviceKey, consts.TopicCommand}, "/")
	b, _ := json.Marshal(data)
	if err := ra.MQTT().Publish(topic, b, 1); err != nil {
		return "", "", err
	}
	return topic, string(b), nil
}

// parseTopic 解析 MQTT 主题，返回 productKey/deviceKey/类型。
func parseTopic(topic string) (productKey, deviceKey, kind string, ok bool) {
	segs := strings.Split(strings.Trim(topic, "/"), "/")
	if len(segs) < 4 || segs[0] != consts.TopicRoot {
		return "", "", "", false
	}
	productKey, deviceKey = segs[1], segs[2]
	if segs[3] == consts.TopicShadow && len(segs) >= 5 {
		return productKey, deviceKey, segs[4], true
	}
	return productKey, deviceKey, segs[3], true
}

// measurementOf 计算设备遥测 measurement。
func measurementOf(deviceID int64) string {
	return consts.MeasurementDevicePrefix + strconv.FormatInt(deviceID, 10)
}
