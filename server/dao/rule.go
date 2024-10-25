package dao

import (
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"openeuler.org/PilotGo/prometheus-plugin/db"
	"openeuler.org/PilotGo/prometheus-plugin/model"
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
