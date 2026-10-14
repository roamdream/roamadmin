package service

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/gogf/gf/v2/os/gtime"

	"github.com/roamdream/roamadmin/kernel/errcode"
	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/modules/device/internal/consts"
	"github.com/roamdream/roamadmin/modules/device/internal/model"
)

// GetShadow 查询设备影子（reported / desired / delta）。
func GetShadow(ctx context.Context, deviceID int64) (*model.ShadowState, error) {
	if _, err := GetDevice(ctx, deviceID); err != nil {
		return nil, err
	}
	reported, desired, version, err := loadShadowMaps(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	return &model.ShadowState{
		DeviceId: deviceID,
		Reported: reported,
		Desired:  desired,
		Version:  version,
	}, nil
}

// UpdateDesired 合并更新影子期望态，返回新版本号。
func UpdateDesired(ctx context.Context, deviceID int64, desired map[string]any) (int64, error) {
	if _, err := GetDevice(ctx, deviceID); err != nil {
		return 0, err
	}
	lockKey := consts.CachePrefixDeviceKey + "shadow:" + itoa(deviceID)
	unlock, err := ra.Lock().Lock(ctx, lockKey, 5*time.Second)
	if err != nil {
		return 0, errcode.Wrap(errcode.CodeInternal, err)
	}
	defer unlock()

	cur, err := loadShadow(ctx, deviceID)
	if err != nil {
		return 0, err
	}
	var merged map[string]any
	if cur != nil && cur.Desired != "" {
		if err := json.Unmarshal([]byte(cur.Desired), &merged); err != nil {
			return 0, errcode.Wrap(errcode.CodeDB, err)
		}
	} else {
		merged = map[string]any{}
	}
	for k, v := range desired {
		merged[k] = v
	}
	nextVersion := int64(1)
	if cur != nil {
		nextVersion = cur.Version + 1
	}
	return nextVersion, upsertShadow(ctx, deviceID, "", marshal(merged), nextVersion, cur == nil)
}

// upsertReported 合并更新影子上报态（MQTT 上行 / HTTP 模拟共用，无租户上下文）。
func upsertReported(ctx context.Context, deviceID int64, reported map[string]any) error {
	lockKey := consts.CachePrefixDeviceKey + "shadow:" + itoa(deviceID)
	unlock, err := ra.Lock().Lock(ctx, lockKey, 5*time.Second)
	if err != nil {
		return errcode.Wrap(errcode.CodeInternal, err)
	}
	defer unlock()

	cur, err := loadShadow(ctx, deviceID)
	if err != nil {
		return err
	}
	var merged map[string]any
	if cur != nil && cur.Reported != "" {
		if err := json.Unmarshal([]byte(cur.Reported), &merged); err != nil {
			return errcode.Wrap(errcode.CodeDB, err)
		}
	} else {
		merged = map[string]any{}
	}
	for k, v := range reported {
		merged[k] = v
	}
	nextVersion := int64(1)
	if cur != nil {
		nextVersion = cur.Version + 1
	}
	return upsertShadow(ctx, deviceID, marshal(merged), "", nextVersion, cur == nil)
}

// loadShadow 读取影子行（不存在返回 nil）。
func loadShadow(ctx context.Context, deviceID int64) (*model.DeviceShadow, error) {
	var list []*model.DeviceShadow
	if err := ra.DB().Model(consts.TableDeviceShadow).Where("device_id", deviceID).Scan(&list); err != nil {
		return nil, errcode.Wrap(errcode.CodeDB, err)
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

// loadShadowMaps 解析影子为 map 形态，并计算 delta。
func loadShadowMaps(ctx context.Context, deviceID int64) (reported, desired map[string]any, version int64, err error) {
	cur, err := loadShadow(ctx, deviceID)
	if err != nil {
		return nil, nil, 0, err
	}
	if cur == nil {
		return map[string]any{}, map[string]any{}, 0, nil
	}
	reported = map[string]any{}
	desired = map[string]any{}
	if cur.Reported != "" {
		_ = json.Unmarshal([]byte(cur.Reported), &reported)
	}
	if cur.Desired != "" {
		_ = json.Unmarshal([]byte(cur.Desired), &desired)
	}
	return reported, desired, cur.Version, nil
}

// upsertShadow 插入或更新影子行（isNew 用于区分 Insert / Update）。
func upsertShadow(ctx context.Context, deviceID int64, reported, desired string, version int64, isNew bool) error {
	data := map[string]any{"version": version, "updated_at": gtime.Now()}
	if reported != "" {
		data["reported"] = reported
	}
	if desired != "" {
		data["desired"] = desired
	}
	if isNew {
		data["id"] = ra.DB().NextID()
		data["device_id"] = deviceID
		if _, ok := data["reported"]; !ok {
			data["reported"] = "{}"
		}
		if _, ok := data["desired"]; !ok {
			data["desired"] = "{}"
		}
		if _, err := ra.DB().Model(consts.TableDeviceShadow).Data(data).Insert(); err != nil {
			return errcode.Wrap(errcode.CodeDB, err)
		}
		return nil
	}
	if _, err := ra.DB().Model(consts.TableDeviceShadow).Where("device_id", deviceID).Data(data).Update(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// ComputeDelta 计算 delta（desired 中与 reported 不一致的键）。
func ComputeDelta(reported, desired map[string]any) map[string]any {
	delta := map[string]any{}
	for k, dv := range desired {
		rv, ok := reported[k]
		if !ok || !reflect.DeepEqual(rv, dv) {
			delta[k] = dv
		}
	}
	return delta
}

func marshal(v map[string]any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func itoa(n int64) string {
	return fmt.Sprintf("%d", n)
}
