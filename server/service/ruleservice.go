package service

import (
	"errors"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/google/uuid"
	"openeuler.org/PilotGo/prometheus-plugin/server/dao"
	"openeuler.org/PilotGo/prometheus-plugin/server/model"
	prometheus "openeuler.org/PilotGo/prometheus-plugin/server/service/prometheus"
)

func AddRule(alert *model.Rule) error {
	alertLabel := uuid.New().String()

	err := prometheus.TestingUpdateRule(testDataJoinToYaml(alert, alertLabel)) //验证新增配置是否有误
	if err != nil {
		return err
	}
	// 更新yaml文件
	yamlData, err := addRuleDataJoinToYaml(alert, alertLabel)
	if err != nil {
		return err
	}
	err = prometheus.UpdateAlertYml(yamlData)
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

func DeleteRuleList(id string) error {
	yamlData, err := deleteRuleDataJoinToYaml(id)
	if err != nil {
		return err
	}

	err = prometheus.UpdateAlertYml(yamlData)
	if err != nil {
		return err
	}

	err = dao.DeleteRule(id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateRule(a *model.Rule) error {
	err := prometheus.TestingUpdateRule(testDataJoinToYaml(a, a.AlertLabel)) //验证编辑的配置是否有误
	if err != nil {
		return err
	}
	// 更新yaml文件
	yamlData, err := updateRuleDataJoinToYaml(a)
	if err != nil {
		return err
	}
	err = prometheus.UpdateAlertYml(yamlData)
	if err != nil {
		return err
	}

	updateAlert := &model.Rule{
		AlertName:      a.AlertName,
		CustomDesc:     a.CustomDesc,
		MonitorMetrics: a.MonitorMetrics,
		AlarmThreshold: a.AlarmThreshold,
		Forsearch:      a.Forsearch,
		Severity:       a.Severity,
		AlertTargets:   a.AlertTargets,
	}

	if err := dao.UpdateRule(a.ID, updateAlert); err != nil {
		if strings.Contains(err.Error(), "重复键违反唯一约束") {
			return errors.New("请确认该规则是否存在或已做修改")
		} else {
			return err
		}
	}
	return nil
}
func SearchRules(search string, query *response.PaginationQ) ([]*model.Rule, int, error) {
	var rules []*model.Rule
	var total int64
	var err error

	if len(search) != 0 {
		rules, total, err = dao.SearchRules(search, query)
		if err != nil {
			return []*model.Rule{}, 0, err
		}
	} else {
		rules, total, err = dao.QueryRulesPage(query)
		if err != nil {
			return []*model.Rule{}, 0, err
		}
	}
	return rules, int(total), nil
}
func GetRuleLevel() []string {
	var levels = []string{"紧急", "严重", "中等严重", "警告"}
	result, err := dao.GetRuleLevel()
	if err != nil {
		return levels
	}
	return removeDuplicates(levels, result)
}
func removeDuplicates(a, b []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, s := range a {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	for _, s := range b {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	return result
}
