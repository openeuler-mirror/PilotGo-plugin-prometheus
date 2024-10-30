package dao

import (
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gorm.io/gorm"
	"openeuler.org/PilotGo/prometheus-plugin/db"
	"openeuler.org/PilotGo/prometheus-plugin/model"
	"openeuler.org/PilotGo/prometheus-plugin/utils"
)

func SaveAlertList(a *model.Alert) error {
	err := db.MySQL.Create(&a).Error
	return err
}
func QueryCompleteAlerts() ([]model.Alert, error) {
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

func QueryAlerts() ([]model.Alert, int, error) {
	var alert []model.Alert
	if err := db.MySQL.Order("id desc").Find(&alert).Error; err != nil {
		return nil, 0, nil
	}

	var total int64
	if err := db.MySQL.Model(&alert).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return alert, int(total), nil
}

func SearchAlertsFiring(alertName, ip, level, handleState string, alertStart, alertEnd model.AlertTime) ([]model.Alert, int64, error) {
	var alert []model.Alert
	dm := db.MySQL.Order("id desc").Where("alert_end_time = ''")
	dm, err := queryFromDB(dm, alertName, ip, level, handleState, alertStart, alertEnd)
	if err != nil {
		logger.Error("时间戳转换失败：%v", err.Error())
	}
	if err := dm.Find(&alert).Error; err != nil {
		return nil, 0, nil
	}

	var total int64
	d := db.MySQL.Where("alert_end_time = ''")
	d, err = queryFromDB(d, alertName, ip, level, handleState, alertStart, alertEnd)
	if err != nil {
		logger.Error("时间戳转换失败：%v", err.Error())
	}
	if err := d.Model(&alert).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return alert, total, nil
}

func SearchAlertsWaitProcess(alertName, ip, level, handleState string, alertStart, alertEnd model.AlertTime) ([]model.Alert, int64, error) {
	var alert []model.Alert

	dm := db.MySQL.Order("id desc").Where("alert_end_time != '' AND handle_state = '已确认'")
	dm, err := queryFromDB(dm, alertName, ip, level, handleState, alertStart, alertEnd)
	if err != nil {
		logger.Error("时间戳转换失败：%v", err.Error())
	}

	subQuery := db.MySQL.Where("alert_end_time != '' AND handle_state = '待确认'")
	if len(alertName) > 0 {
		subQuery = subQuery.Where("alert_name LIKE ? ", "%"+alertName+"%")
	}
	if len(ip) > 0 {
		subQuery = subQuery.Where("ip LIKE ? ", "%"+ip+"%")
	}
	if len(level) > 0 {
		subQuery = subQuery.Where("alert_level = ? ", level)
	}
	if len(handleState) > 0 {
		subQuery = subQuery.Where("handle_state = ? ", handleState)
	}

	if len(alertStart.Start) != 0 && len(alertStart.End) != 0 {
		start, err := utils.UnixTimeToShanghai(alertStart.Start)
		if err != nil {
			return nil, 0, err
		}
		end, err := utils.UnixTimeToShanghai(alertStart.End)
		if err != nil {
			return nil, 0, err
		}
		subQuery = subQuery.Where("to_timestamp(alert_time, 'YYYY-MM-DD HH24:MI:SS') >= ? AND to_timestamp(alert_time, 'YYYY-MM-DD HH24:MI:SS') <= ?", start, end)
	}
	if len(alertEnd.Start) != 0 && len(alertEnd.End) != 0 {
		start, err := utils.UnixTimeToShanghai(alertEnd.Start)
		if err != nil {
			return nil, 0, err
		}
		end, err := utils.UnixTimeToShanghai(alertEnd.End)
		if err != nil {
			return nil, 0, err
		}
		subQuery = subQuery.Where("to_timestamp(alert_end_time, 'YYYY-MM-DD HH24:MI:SS') >= ? AND to_timestamp(alert_end_time, 'YYYY-MM-DD HH24:MI:SS') <= ?", start, end)
	}
	dm = dm.Or(subQuery)
	if err := dm.Find(&alert).Error; err != nil {
		return nil, 0, nil
	}

	var total int64
	d := db.MySQL.Where("alert_end_time != '' AND handle_state = '已确认'")
	d, err = queryFromDB(d, alertName, ip, level, handleState, alertStart, alertEnd)
	if err != nil {
		logger.Error("时间戳转换失败：%v", err.Error())
	}

	sub := db.MySQL.Where("alert_end_time != '' AND handle_state = '待确认'")
	if len(alertName) > 0 {
		sub = sub.Where("alert_name LIKE ? ", "%"+alertName+"%")
	}
	if len(ip) > 0 {
		sub = sub.Where("ip LIKE ? ", "%"+ip+"%")
	}
	if len(level) > 0 {
		sub = sub.Where("alert_level = ? ", level)
	}
	if len(handleState) > 0 {
		sub = sub.Where("handle_state = ? ", handleState)
	}
	if len(alertStart.Start) != 0 && len(alertStart.End) != 0 {
		start, err := utils.UnixTimeToShanghai(alertStart.Start)
		if err != nil {
			return nil, 0, err
		}
		end, err := utils.UnixTimeToShanghai(alertStart.End)
		if err != nil {
			return nil, 0, err
		}
		sub = sub.Where("to_timestamp(alert_time, 'YYYY-MM-DD HH24:MI:SS') >= ? AND to_timestamp(alert_time, 'YYYY-MM-DD HH24:MI:SS') <= ?", start, end)
	}
	if len(alertEnd.Start) != 0 && len(alertEnd.End) != 0 {
		start, err := utils.UnixTimeToShanghai(alertEnd.Start)
		if err != nil {
			return nil, 0, err
		}
		end, err := utils.UnixTimeToShanghai(alertEnd.End)
		if err != nil {
			return nil, 0, err
		}
		sub = sub.Where("to_timestamp(alert_end_time, 'YYYY-MM-DD HH24:MI:SS') >= ? AND to_timestamp(alert_end_time, 'YYYY-MM-DD HH24:MI:SS') <= ?", start, end)
	}
	d = d.Or(sub)
	if err := d.Model(&alert).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return alert, total, nil
}

func SearchAlertsCompleted(alertName, ip, level, handleState string, alertStart, alertEnd model.AlertTime) ([]model.Alert, int64, error) {
	var alert []model.Alert
	dm := db.MySQL.Order("id desc").Where("handle_state = '已完成'")
	dm, err := queryFromDB(dm, alertName, ip, level, handleState, alertStart, alertEnd)
	if err != nil {
		logger.Error("时间戳转换失败：%v", err.Error())
	}
	if err := dm.Find(&alert).Error; err != nil {
		return nil, 0, nil
	}

	var total int64
	d := db.MySQL.Where("handle_state = '已完成'")
	d, err = queryFromDB(d, alertName, ip, level, handleState, alertStart, alertEnd)
	if err != nil {
		logger.Error("时间戳转换失败：%v", err.Error())
	}
	if err := d.Model(&alert).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return alert, total, nil
}

func SearchAlerts(alertName, ip, level, handleState string, alertStart, alertEnd model.AlertTime) ([]model.Alert, int64, error) {
	var alert []model.Alert
	dm := db.MySQL.Order("id desc")
	dm, err := queryFromDB(dm, alertName, ip, level, handleState, alertStart, alertEnd)
	if err != nil {
		logger.Error("时间戳转换失败：%v", err.Error())
	}
	if err := dm.Find(&alert).Error; err != nil {
		return nil, 0, nil
	}

	var total int64
	d := db.MySQL
	d, err = queryFromDB(d, alertName, ip, level, handleState, alertStart, alertEnd)
	if err != nil {
		logger.Error("时间戳转换失败：%v", err.Error())
	}
	if err := d.Model(&alert).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return alert, total, nil
}

func UpdateHandleState(id int, alert *model.Alert) error {
	var a model.Alert
	err := db.MySQL.Model(&a).Where("id = ?", id).Updates(alert).Error
	return err
}

func queryFromDB(db *gorm.DB, alertName, ip, level, handleState string, alertStart, alertEnd model.AlertTime) (*gorm.DB, error) {
	if len(alertName) > 0 {
		db = db.Where("alert_name LIKE ? ", "%"+alertName+"%")
	}
	if len(ip) > 0 {
		db = db.Where("ip LIKE ? ", "%"+ip+"%")
	}
	if len(level) > 0 {
		db = db.Where("alert_level = ? ", level)
	}
	if len(handleState) > 0 {
		db = db.Where("handle_state = ? ", handleState)
	}
	if len(alertStart.Start) != 0 && len(alertStart.End) != 0 {
		start, err := utils.UnixTimeToShanghai(alertStart.Start)
		if err != nil {
			return db, err
		}
		end, err := utils.UnixTimeToShanghai(alertStart.End)
		if err != nil {
			return db, err
		}
		db = db.Where("to_timestamp(alert_time, 'YYYY-MM-DD HH24:MI:SS') >= ? AND to_timestamp(alert_time, 'YYYY-MM-DD HH24:MI:SS') <= ?", start, end)
	}
	if len(alertEnd.Start) != 0 && len(alertEnd.End) != 0 {
		start, err := utils.UnixTimeToShanghai(alertEnd.Start)
		if err != nil {
			return db, err
		}
		end, err := utils.UnixTimeToShanghai(alertEnd.End)
		if err != nil {
			return db, err
		}
		db = db.Where("to_timestamp(alert_end_time, 'YYYY-MM-DD HH24:MI:SS') >= ? AND to_timestamp(alert_end_time, 'YYYY-MM-DD HH24:MI:SS') <= ?", start, end)
	}
	return db, nil
}
