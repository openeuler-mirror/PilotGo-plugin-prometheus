package dao

import (
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
