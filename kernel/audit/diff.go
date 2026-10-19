package audit

import (
	"encoding/json"
	"reflect"
)

// Diff 计算 before/after 之间发生变化的字段，返回 JSON 字符串：
// {"email":{"from":"a@x.com","to":"b@x.com"}}；无变化返回空字符串。
// before/after 可为 struct、map 等任意可 JSON 序列化的值，字段以 JSON key 对齐比较。
func Diff(before, after any) string {
	bm := toJSONMap(before)
	am := toJSONMap(after)
	if bm == nil || am == nil {
		return ""
	}
	changes := map[string]any{}
	for k, av := range am {
		bv, ok := bm[k]
		if !ok || !reflect.DeepEqual(bv, av) {
			changes[k] = map[string]any{"from": bv, "to": av}
		}
	}
	if len(changes) == 0 {
		return ""
	}
	b, err := json.Marshal(changes)
	if err != nil {
		return ""
	}
	return string(b)
}

func toJSONMap(v any) map[string]any {
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	return m
}
