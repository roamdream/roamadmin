package main

// 集中注册要启用的业务模块（增删模块 = 增删一行 import）。
// 各模块在自身包的 init() 中调用 module.Register 完成编译期注册，
// 此处仅需 import 即触发注册，main 启动时由内核统一驱动其生命周期。
import (
	_ "github.com/roamdream/roamadmin/modules/demo"
	_ "github.com/roamdream/roamadmin/modules/device"
	_ "github.com/roamdream/roamadmin/modules/gateway"
	_ "github.com/roamdream/roamadmin/modules/system"
)
