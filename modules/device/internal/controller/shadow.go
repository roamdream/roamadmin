package controller

import (
	"context"

	"github.com/roamdream/roamadmin/kernel/utils"
	"github.com/roamdream/roamadmin/modules/device/api"
	"github.com/roamdream/roamadmin/modules/device/internal/service"
)

// Shadow 设备影子 + 遥测控制器。
type Shadow struct{}

// NewShadow 创建影子控制器。
func NewShadow() *Shadow { return &Shadow{} }

// GetShadow 查询设备影子。
func (c *Shadow) GetShadow(ctx context.Context, req *api.ShadowGetReq) (*api.ShadowGetRes, error) {
	state, err := service.GetShadow(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.ShadowGetRes{
		DeviceId: state.DeviceId,
		Reported: state.Reported,
		Desired:  state.Desired,
		Delta:    service.ComputeDelta(state.Reported, state.Desired),
		Version:  state.Version,
	}, nil
}

// UpdateDesired 更新影子期望态。
func (c *Shadow) UpdateDesired(ctx context.Context, req *api.ShadowUpdateDesiredReq) (*api.ShadowUpdateDesiredRes, error) {
	version, err := service.UpdateDesired(ctx, req.Id, req.Desired)
	if err != nil {
		return nil, err
	}
	return &api.ShadowUpdateDesiredRes{Version: version}, nil
}

// GetDelta 查询影子增量。
func (c *Shadow) GetDelta(ctx context.Context, req *api.ShadowDeltaReq) (*api.ShadowDeltaRes, error) {
	state, err := service.GetShadow(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.ShadowDeltaRes{Delta: service.ComputeDelta(state.Reported, state.Desired)}, nil
}

// ReportTelemetry 上报遥测（模拟设备上行）。
func (c *Shadow) ReportTelemetry(ctx context.Context, req *api.TelemetryReportReq) (*api.TelemetryReportRes, error) {
	if err := service.ReportTelemetry(ctx, req.Id, req.Fields); err != nil {
		return nil, err
	}
	return &api.TelemetryReportRes{Ok: true}, nil
}

// QueryTelemetry 查询最近遥测。
func (c *Shadow) QueryTelemetry(ctx context.Context, req *api.TelemetryQueryReq) (*api.TelemetryQueryRes, error) {
	list, err := service.QueryTelemetry(ctx, req.Id, req.Limit)
	if err != nil {
		return nil, err
	}
	out := make([]*api.TelemetryPoint, 0, len(list))
	for _, p := range list {
		out = append(out, &api.TelemetryPoint{
			Measurement: p.Measurement,
			Tags:        p.Tags,
			Fields:      p.Fields,
			Time:        utils.TimeRFC3339T(p.Time),
		})
	}
	return &api.TelemetryQueryRes{List: out}, nil
}
