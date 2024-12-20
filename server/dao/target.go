/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Fri Oct 20 11:41:45 2023 +0800
 */
package dao

import (
	"errors"

	"openeuler.org/PilotGo/prometheus-plugin/server/db"
	"openeuler.org/PilotGo/prometheus-plugin/server/model"
)

func GetPrometheusTarget() ([]string, error) {
	var ips []model.PrometheusTarget
	err := db.MySQL.Raw("SELECT * FROM prometheus_target ORDER BY id DESC").Scan(&ips).Error
	if err != nil {
		return []string{}, err
	}

	if len(ips) == 0 {
		return []string{}, errors.New("ip targets is null")
	}
	var targets []string
	for _, ip := range ips {
		target := ip.TargetIP + ":" + ip.Port
		targets = append(targets, target)
	}
	return targets, nil
}

func QueryPrometheusTargets() ([]model.PrometheusTarget, error) {
	var targets []model.PrometheusTarget
	err := db.MySQL.Raw("SELECT * FROM prometheus_target ORDER BY id DESC").Scan(&targets).Error
	if err != nil {
		return targets, err
	}
	return targets, nil
}

func AddPrometheusTarget(pt *model.PrometheusTarget) error {
	t := model.PrometheusTarget{
		UUID:     pt.UUID,
		TargetIP: pt.TargetIP,
		Port:     pt.Port,
	}
	err := db.MySQL.Save(&t).Error
	if err != nil {
		return err
	}
	return nil
}

func DeletePrometheusTarget(pt *model.PrometheusTarget) error {
	var t model.PrometheusTarget
	err := db.MySQL.Where("uuid = ?", pt.UUID).Unscoped().Delete(t).Error
	if err != nil {
		return err
	}
	return nil
}

func IsExistTargetUUID(uuid string) (bool, error) {
	var r model.PrometheusTarget
	err := db.MySQL.Where("uuid = ?", uuid).Find(&r).Error
	if err != nil {
		return false, errors.New("查询数据库失败：" + err.Error())
	}
	return r.ID != 0, nil
}
