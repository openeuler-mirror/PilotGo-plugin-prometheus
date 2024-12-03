/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Oct 23 20:21:56 2024 +0800
 */
package dao

import (
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"openeuler.org/PilotGo/prometheus-plugin/server/db"
	"openeuler.org/PilotGo/prometheus-plugin/server/model"
)

func SaveRuleList(a *model.Rule) error {
	return db.MySQL.Transaction(func(tx *gorm.DB) error {
		r := model.Rule{
			AlertName:      a.AlertName,
			CustomDesc:     a.CustomDesc,
			MonitorMetrics: a.MonitorMetrics,
			AlarmThreshold: a.AlarmThreshold,
			Forsearch:      a.Forsearch,
			Severity:       a.Severity,
			AlertLabel:     a.AlertLabel,
		}
		if err := tx.Create(&r).Error; err != nil {
			return err
		}

		for _, target := range a.AlertTargets {
			target.RuleId = r.ID
			if err := tx.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&target).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func QueryRulesPage(query *response.PaginationQ) ([]*model.Rule, int64, error) {
	var alert []*model.Rule
	if err := db.MySQL.Limit(query.PageSize).Offset((query.Page - 1) * query.PageSize).Order("id desc").Preload("AlertTargets").Find(&alert).Error; err != nil {
		return alert, 0, nil
	}
	// 确保 AlertTargets 是空数组而不是 null
	for i := range alert {
		if alert[i].AlertTargets == nil {
			alert[i].AlertTargets = []model.RuleTarget{}
			logger.Info("%v", alert[i].AlertTargets)
		}
	}

	var total int64
	if err := db.MySQL.Model(&alert).Count(&total).Error; err != nil {
		return alert, 0, err
	}
	return alert, total, nil
}
func SearchRules(search string, query *response.PaginationQ) ([]*model.Rule, int64, error) {
	var alert []*model.Rule
	if err := db.MySQL.Limit(query.PageSize).Offset((query.Page-1)*query.PageSize).Order("id desc").Where("severity = ? ", search).Preload("AlertTargets").Find(&alert).Error; err != nil {
		return alert, 0, nil
	}
	// 确保 AlertTargets 是空数组而不是 null
	for i := range alert {
		if alert[i].AlertTargets == nil {
			alert[i].AlertTargets = []model.RuleTarget{}
			logger.Info("%v", alert[i].AlertTargets)
		}
	}

	var total int64
	if err := db.MySQL.Where("severity = ?", search).Model(&alert).Count(&total).Error; err != nil {
		return alert, 0, err
	}
	return alert, total, nil
}
func DeleteRule(id string) error {
	err := db.MySQL.Where("id = ?", id).Delete(&model.Rule{}).Error
	return err
}

func UpdateRule(id int, alert *model.Rule) error {
	tx := db.MySQL.Begin()

	if err := tx.Model(&model.Rule{}).Where("id = ?", id).Select("alert_name", "custom_desc", "monitor_metrics", "alarm_threshold", "forsearch", "severity", "batches").Updates(alert).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("rule_id = ?", id).Delete(&model.RuleTarget{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, target := range alert.AlertTargets {
		target.RuleId = id
		if err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&target).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func QueryRules() ([]model.Rule, error) {
	var alerts []model.Rule
	err := db.MySQL.Order("id desc").Preload("AlertTargets").Find(&alerts).Error
	return alerts, err
}

func QueryRulesNotIncludedId(id string) ([]model.Rule, error) {
	var alerts []model.Rule
	err := db.MySQL.Order("id desc").Where("id <> ?", id).Preload("AlertTargets").Find(&alerts).Error
	return alerts, err
}
func GetRuleLevel() ([]string, error) {
	var levels []string
	err := db.MySQL.Model(&model.Rule{}).Distinct("severity").Pluck("severity", &levels).Error
	return levels, err
}
