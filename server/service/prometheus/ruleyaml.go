package initprometheus

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
	"openeuler.org/PilotGo/prometheus-plugin/model"
	"openeuler.org/PilotGo/prometheus-plugin/utils"
)

const ruleYaml = "/etc/prometheus/rules.yaml"

func TestingUpdateRule(alertYml *model.AlertRuleYaml) error { // 生成测试文件，验证这个文件的合规性
	filename := "/etc/prometheus/test.yaml"
	file, err := os.Create(filename)
	if err != nil {
		return errors.New("创建测试文件出错：" + err.Error())
	}
	defer file.Close()

	err = updateAlertRuleYaml(filename, alertYml)
	if err != nil {
		return err
	}

	err = os.Remove(filename)
	if err != nil {
		return errors.New("删除测试文件出错：" + err.Error())
	}
	return nil
}
func UpdateAlertYml(alertYml *model.AlertRuleYaml) error { // 更新文件，验证这个文件的合规性
	if err := backupPromAlertRuleYaml(); err != nil {
		return err
	}
	if err := updateAlertRuleYaml(ruleYaml, alertYml); err != nil {
		return rollBackPromAlertRuleYaml()
	}
	if err := reloadPrometheus(); err != nil {
		return rollBackPromAlertRuleYaml()
	}
	return nil
}

func updateAlertRuleYaml(ruleyaml string, alertYml *model.AlertRuleYaml) error {
	f, err := os.OpenFile(ruleyaml, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	yaml.FutureLineWrap()
	encoder := yaml.NewEncoder(f)

	err = encoder.Encode(alertYml)
	if err != nil {
		return err
	}

	err = validityAlertRuleYaml(ruleyaml)
	if err != nil {
		return err
	}
	return nil
}
func validityAlertRuleYaml(ruleyaml string) error {
	cmd := "promtool check config /etc/prometheus/prometheus.yml"
	_, _, stderr, err := utils.RunCommand(cmd)
	if err != nil {
		return err
	}
	if len(stderr) != 0 {
		return errors.New("prometheus配置文件有误:" + stderr)
	}

	cmdRule := "promtool check rules " + ruleyaml
	_, _, stderr, err = utils.RunCommand(cmdRule)
	if err != nil {
		return err
	}
	if len(stderr) != 0 {
		return errors.New("告警规则配置文件有误:" + stderr)
	}
	return nil
}

func backupPromAlertRuleYaml() error {
	cmd := "cp /etc/prometheus/rules.yaml /etc/prometheus/rules.yaml.bak"
	_, _, stderr, err := utils.RunCommand(cmd)
	if err != nil {
		return err
	}
	if len(stderr) != 0 {
		return errors.New("备份告警规则文件有误:" + stderr)
	}
	return nil
}
func rollBackPromAlertRuleYaml() error {
	cmd := "cp /etc/prometheus/rules.yaml.bak /etc/prometheus/rules.yaml"
	_, _, stderr, err := utils.RunCommand(cmd)
	if err != nil {
		return err
	}
	if len(stderr) != 0 {
		return errors.New("恢复规则配置文件有误:" + stderr)
	}
	return nil
}
func reloadPrometheus() error {
	cmd := "systemctl stop prometheus && systemctl start prometheus"
	_, _, stderr, err := utils.RunCommand(cmd)
	if err != nil {
		return err
	}
	if len(stderr) != 0 {
		return errors.New("重启kylin-monitor失败:" + stderr)
	}
	return nil
}
