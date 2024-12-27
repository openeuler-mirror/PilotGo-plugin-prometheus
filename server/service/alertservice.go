/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Oct 29 16:42:12 2024 +0800
 */
package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"openeuler.org/PilotGo/prometheus-plugin/server/config"
	"openeuler.org/PilotGo/prometheus-plugin/server/dao"
	"openeuler.org/PilotGo/prometheus-plugin/server/model"
	prometheus "openeuler.org/PilotGo/prometheus-plugin/server/service/prometheus"
)

const (
	/*
	*告警各状态：活跃、待处理、已处理
	 */
	Firing      = "活跃"
	WaitProcess = "待处理"
	Completed   = "已处理"
)

func SearchAlerts(alertName, ip, level, handleState, alertState string, alertStart, alertEnd model.AlertTime) ([]model.Alert, int, error) {
	switch alertState {
	case Firing:
		alerts, total, err := dao.SearchAlertsFiring(alertName, ip, level, handleState, alertStart, alertEnd)
		if err != nil {
			return []model.Alert{}, 0, err
		}
		return alerts, int(total), nil
	case WaitProcess:
		alerts, total, err := dao.SearchAlertsWaitProcess(alertName, ip, level, handleState, alertStart, alertEnd)
		if err != nil {
			return []model.Alert{}, 0, err
		}
		return alerts, int(total), nil
	case Completed:
		alerts, total, err := dao.SearchAlertsCompleted(alertName, ip, level, handleState, alertStart, alertEnd)
		if err != nil {
			return []model.Alert{}, 0, err
		}

		return alerts, int(total), nil
	default:
		alerts, total, err := dao.SearchAlerts(alertName, ip, level, handleState, alertStart, alertEnd)
		if err != nil {
			return []model.Alert{}, 0, err
		}

		return alerts, int(total), nil
	}
}

func UpdateHandleState(id int, state string) error {
	if state == "已确认" {
		err := dao.UpdateHandleState(id, &model.Alert{
			ConfirmTime: time.Now().Format("2006-01-02 15:04:05"),
			HandleState: state,
		})
		return err
	}
	if state == "已完成" {
		err := dao.UpdateHandleState(id, &model.Alert{
			CompleteTime: time.Now().Format("2006-01-02 15:04:05"),
			HandleState:  state,
		})
		return err
	}
	return nil
}

func GetAlertLevel() []string {
	var levels = []string{}
	result, err := dao.GetAlertLevel()
	if err != nil {
		return levels
	}
	return result
}

func PullAlert() error {
	var previousAlerts []model.AlertResponse

	daoAlert, err := dao.QueryCompleteAlerts()
	if err != nil {
		return err
	}

	previousAlerts, err = pullAlert()
	if err != nil {
		return err
	}
	if err = processAlerts(daoAlert, previousAlerts); err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if prometheus.SetDelayTicker {
				prometheus.SetDelayTicker = false
				continue
			}
			alerts, err := pullAlert()
			if err != nil {
				logger.Error("Error pull alerts from prometheus: %v", err.Error())
				continue
			}
			processAlerts(previousAlerts, alerts)
			previousAlerts = alerts
		}
	}()
	return nil
}
func pullAlert() ([]model.AlertResponse, error) {
	remote := "http://" + config.Config().PrometheusServer.Addr + "/api/v1/alerts"

	request, err := http.NewRequest("GET", remote, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Transport: &http.Transport{}}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer client.CloseIdleConnections()

	var alerts model.AlertsResponse
	err = json.Unmarshal(body, &alerts)
	if err != nil {
		return nil, err
	}
	if alerts.Status != "success" {
		return nil, err
	}
	var data []model.AlertResponse
	for _, alert := range alerts.Data.Alerts {
		if alert.State == "firing" {
			data = append(data, alert)
		}
	}
	return data, nil
}

func processAlerts(old interface{}, new []model.AlertResponse) error {

	arrType := reflect.TypeOf(old)
	if arrType.Kind() != reflect.Slice {
		return fmt.Errorf("Unknown")
	}
	elemType := arrType.Elem()

	var oldMap map[string]interface{}
	if elemType == reflect.TypeOf(model.AlertResponse{}) {
		oldMap = make(map[string]interface{})
		for _, a := range old.([]model.AlertResponse) {
			oldMap[getAlertKey(a)] = a
		}
	} else if elemType == reflect.TypeOf(model.Alert{}) {
		oldMap = make(map[string]interface{})
		for _, o := range old.([]model.Alert) {
			oldMap[getDBAlertKey(o)] = o
		}
	}

	newMap := make(map[string]model.AlertResponse)
	for _, a := range new {
		newMap[getAlertKey(a)] = a
	}

	for key := range oldMap {
		if _, ok := newMap[key]; ok {
			delete(newMap, key)
		} else {
			if elemType == reflect.TypeOf(model.Alert{}) {
				oldAlert := oldMap[key].(model.Alert)
				err := dao.UpdateAlert(oldAlert.AlertName, oldAlert.GroupId, oldAlert.IP, &model.Alert{
					AlertEndTime: time.Now().Format("2006-01-02 15:04:05"),
				})
				if err != nil {
					logger.Error("Error update alert end time for alertname=%s,group=%s,ip=%s :%v", oldAlert.AlertName, oldAlert.GroupId, oldAlert.IP, err.Error())
				}
			} else if elemType == reflect.TypeOf(model.AlertResponse{}) {
				oldAlert := oldMap[key].(model.AlertResponse)
				err := dao.UpdateAlert(oldAlert.Labels.AlertName, oldAlert.Labels.Group, strings.Split(oldAlert.Labels.Instance, ":")[0], &model.Alert{
					AlertEndTime: time.Now().Format("2006-01-02 15:04:05"),
				})
				if err != nil {
					logger.Error("Error update alert end time for alertname=%s,group=%s,ip=%s :%v", oldAlert.Labels.AlertName, oldAlert.Labels.Group, strings.Split(oldAlert.Labels.Instance, ":")[0], err.Error())
				}
			}
		}
	}
	for _, n := range newMap {
		alertTime, err := utcTimeToAsia(n.ActiveAt)
		if err != nil {
			return err
		}
		err = dao.SaveAlertList(&model.Alert{
			AlertName:   n.Labels.AlertName,
			GroupId:     n.Labels.Group,
			IP:          strings.Split(n.Labels.Instance, ":")[0],
			AlertLevel:  n.Labels.Severity,
			Summary:     n.Annotations.Summary,
			Metric:      n.Labels.Metric,
			Description: n.Annotations.Description,
			AlertTime:   alertTime,
		})
		if err != nil {
			logger.Error("保存失败：%v", err.Error())
		}
	}
	return nil
}

func getDBAlertKey(alert model.Alert) string {
	return fmt.Sprintf("%s--%s--%s", alert.AlertName, alert.GroupId, alert.IP)
}
func getAlertKey(alert model.AlertResponse) string {
	return fmt.Sprintf("%s--%s--%s", alert.Labels.AlertName, alert.Labels.Group, strings.Split(alert.Labels.Instance, ":")[0])
}
func utcTimeToAsia(utcTimeStr string) (string, error) {
	utcTime, err := time.Parse(time.RFC3339Nano, utcTimeStr)
	if err != nil {
		return "", err
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("加载时区出错:", err)
		return "", err
	}

	shanghaiTime := utcTime.In(loc)
	return shanghaiTime.Format("2006-01-02 15:04:05"), nil
}
