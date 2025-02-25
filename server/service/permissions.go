/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Feb 25 16:23:18 2025 +0800
 */
package service

import (
	"gitee.com/openeuler/PilotGo/sdk/common"
	"openeuler.org/PilotGo/prometheus-plugin/server/plugin"
)

func AddPermissions() {
	var pe []common.Permission
	p1 := common.Permission{
		Resource: "monitor_operate",
		Operate:  "button",
	}
	p2 := common.Permission{
		Resource: "monitor",
		Operate:  "menu",
	}

	p := append(pe, p1, p2)
	plugin.Client.RegisterPermission(p)
}
