/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Thu Oct 19 11:33:37 2023 +0800
 */
package model

type PrometheusTarget struct {
	ID       uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UUID     string `json:"uuid"`
	TargetIP string `json:"targetIp"`
	Port     string `json:"port"`
}

type PrometheusObject struct {
	Targets []string `json:"targets"`
	Labels  struct {
		Group string `json:"group"`
	} `json:"labels"`
}
