/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Oct 29 16:42:12 2024 +0800
 */
package model

type Alert struct {
	ID           int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	AlertName    string `gorm:"uniqueIndex:idx_alerts, length:100" json:"alertName"`
	GroupId      string `gorm:"uniqueIndex:idx_alerts, length:100" json:"groupId"`
	IP           string `gorm:"uniqueIndex:idx_alerts, length:100" json:"ip"`
	AlertLevel   string `json:"level"`
	Summary      string `json:"summary"`
	Metric       string `json:"metric"`
	Description  string `json:"description"`
	AlertTime    string `json:"alertTime"`
	AlertEndTime string `gorm:"uniqueIndex:idx_alerts, length:100" json:"alertEndTime"`
	ConfirmTime  string `json:"confirmTime"`
	CompleteTime string `json:"completeTime"`
	HandleState  string `json:"handleState"`
}

type AlertTime struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type AlertsResponse struct {
	Status string `json:"status"`
	Data   struct {
		Alerts []AlertResponse `json:"alerts"`
	} `json:"data"`
}

type AlertResponse struct {
	Labels struct {
		AlertName      string `json:"alertname"`
		Group          string `json:"group"`
		UserName       string `json:"userName"`
		DepartmentName string `json:"departmentName"`
		Instance       string `json:"instance"`
		Job            string `json:"job"`
		Severity       string `json:"severity"`
		Metric         string `json:"metric"`
	} `json:"labels"`
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
	State    string `json:"state"`
	ActiveAt string `json:"activeAt"`
	Value    string `json:"value"`
}
