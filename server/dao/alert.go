package dao

import (
	"openeuler.org/PilotGo/prometheus-plugin/db"
	"openeuler.org/PilotGo/prometheus-plugin/model"
)

func SaveAlertList(a *model.Alert) error {
	err := db.MySQL.Create(&a).Error
	return err
}
func QueryAlerts() ([]model.Alert, error) {
	var alerts []model.Alert
	err := db.MySQL.Where("alert_end_time = ''").Find(&alerts).Error
	return alerts, err
}
func UpdateAlert(alertName, groupId, ip string, alert *model.Alert) error {
	var a model.Alert
	var maxID uint
	db.MySQL.Model(&a).Where("alert_name = ? AND group_id = ? AND ip = ?", alertName, groupId, ip).Order("id desc").Limit(1).Pluck("id", &maxID)
	err := db.MySQL.Model(&a).Where("id = ? AND alert_name = ? AND group_id = ? AND ip = ?", maxID, alertName, groupId, ip).Updates(alert).Error
	return err
}
