package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"openeuler.org/PilotGo/prometheus-plugin/dao"
	"openeuler.org/PilotGo/prometheus-plugin/model"
	initprometheus "openeuler.org/PilotGo/prometheus-plugin/service/prometheus"
)

func AddRule(alert *model.Rule) error {
	alertLabel := uuid.New().String()

	err := initprometheus.TestingUpdateRule(testDataJoinToYaml(alert, alertLabel)) //验证新增配置是否有误
	if err != nil {
		return err
	}
	// 更新yaml文件
	yamlData, err := addRuleDataJoinToYaml(alert, alertLabel)
	if err != nil {
		return err
	}
	err = initprometheus.UpdateAlertYml(yamlData)
	if err != nil {
		return err
	}

	err = dao.SaveRuleList(&model.Rule{
		AlertName:      alert.AlertName,
		CustomDesc:     alert.CustomDesc,
		MonitorMetrics: alert.MonitorMetrics,
		AlarmThreshold: alert.AlarmThreshold,
		Forsearch:      alert.Forsearch,
		Severity:       alert.Severity,
		AlertTargets:   alert.AlertTargets,
		AlertLabel:     alertLabel,
	})
	if err != nil {
		if strings.Contains(err.Error(), "重复键违反唯一约束") {
			return errors.New("请勿重复添加告警规则")
		} else {
			return err
		}
	}
	return nil
}
