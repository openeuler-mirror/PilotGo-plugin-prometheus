package service

import (
	"openeuler.org/PilotGo/prometheus-plugin/server/dao"
	"openeuler.org/PilotGo/prometheus-plugin/server/model"
)

const (
	CPU               = "cpu使用率"
	Memory            = "内存使用率"
	NetWorkReceive    = "网络流入"
	NetWorkTransmit   = "网络流出"
	FileSystemStorage = "磁盘容量"
	Delay             = "服务器宕机"
	DiskIO            = "磁盘IO"
	TCPNum            = "TCP连接数"
)

type AlertInfo struct {
	Expression  string // 告警规则表达式
	Description string
	Summary     string //告警描述
}

func alertTemplate(metrics string, threshold string, label string) *AlertInfo {
	switch metrics {
	case CPU:
		return &AlertInfo{
			Expression:  "100 - (avg by (instance,group)(irate(node_cpu_seconds_total{mode='idle',group='" + label + "'}[5m])) * 100) > " + threshold,
			Description: "CPU使用大于" + threshold + "%，当前使用率{{ $value }}%. ",
			Summary:     "CPU使用率过高,请尽快处理!",
		}

	case Memory:
		return &AlertInfo{
			Expression:  "(node_memory_MemFree_bytes{group='" + label + "''}+node_memory_Cached_bytes{group='" + label + "'}+node_memory_Buffers_bytes{group='" + label + "'}) / node_memory_MemTotal_bytes{group='" + label + "'} * 100 > " + threshold,
			Summary:     "内存使用率过高，请尽快处理!",
			Description: "内存使用率超过" + threshold + "%,当前使用率{{ $value }}%.",
		}
	case NetWorkReceive:
		return &AlertInfo{
			Expression:  "((sum(rate (node_network_receive_bytes_total{device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*',group='" + label + "'}[5m])) by (instance,group,userName,departmentName)) / 100) > " + threshold,
			Summary:     "流入网络带宽过高，请尽快处理!",
			Description: "流入网络带宽持续高于" + threshold + "KB. RX带宽使用量{{$value}} KB.",
		}
	case NetWorkTransmit:
		return &AlertInfo{
			Expression:  "((sum(rate (node_network_transmit_bytes_total{device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*',group='" + label + "'}[5m])) by (instance,group,userName,departmentName)) / 100) > " + threshold,
			Summary:     "流出网络带宽过高，请尽快处理!",
			Description: "流出网络带宽持续高于" + threshold + "KB. RX带宽使用量{{$value}} KB.",
		}
	case FileSystemStorage:
		return &AlertInfo{
			Expression:  "100-(node_filesystem_avail_bytes{mountpoint='/',group='" + label + "'}/node_filesystem_size_bytes{mountpoint='/',group='" + label + "'})*100 > " + threshold,
			Summary:     "磁盘分区使用率过高，请尽快处理",
			Description: "磁盘分区使用大于" + threshold + "%，当前使用率{{ $value }}%.",
		}
	case Delay:
		return &AlertInfo{
			Expression:  "up{group='" + label + "'} == 0",
			Summary:     "服务器宕机，请尽快处理!",
			Description: "服务器宕机时间过长，当前状态{{ $value }}.",
		}
	case DiskIO:
		return &AlertInfo{
			Expression:  "avg(irate(node_disk_io_time_seconds_total{group='" + label + "'}[1m])) by(instance,job,group,userName,departmentName)* 100 > " + threshold,
			Summary:     "流入磁盘IO使用率过高,请尽快处理!",
			Description: "流入磁盘IO大于" + threshold + "%,当前使用率{{ $value }}%.",
		}
	case TCPNum:
		return &AlertInfo{
			Expression:  "node_netstat_Tcp_CurrEstab {group='" + label + "'}> " + threshold,
			Summary:     "TCP连接数过高!",
			Description: "TCP连接数大于" + threshold + ",当前连接数为{{ $value }}.",
		}
	default:
		return &AlertInfo{}
	}
}

func testDataJoinToYaml(alert *model.Rule, label string) *model.AlertRuleYaml {
	template := alertTemplate(alert.MonitorMetrics, alert.AlarmThreshold, label)
	testRule := &model.AlertRule{
		AlertName:  alert.AlertName,
		Expression: template.Expression,
		Forsearch:  alert.Forsearch + "s",
		Labels: struct {
			Severity string "yaml:\"severity\" json:\"severity\""
			Metric   string "yaml:\"metric\" json:\"metric\""
		}{
			Severity: alert.Severity,
			Metric:   alert.MonitorMetrics,
		},
		Annotations: struct {
			Description string "yaml:\"description\" json:\"description\""
			Summary     string "yaml:\"summary\" json:\"summary\""
		}{
			Summary:     template.Summary,
			Description: template.Description,
		},
	}
	testAlerts := &model.AlertRules{
		RuleName:   "test",
		AlertRules: []model.AlertRule{*testRule},
	}
	testYml := &model.AlertRuleYaml{
		Groups: []model.AlertRules{*testAlerts},
	}
	return testYml
}

func addRuleDataJoinToYaml(alert *model.Rule, label string) (*model.AlertRuleYaml, error) {
	rules, err := dao.QueryRules()
	if err != nil {
		return &model.AlertRuleYaml{}, err
	}
	var alertRules []model.AlertRule
	for _, rule := range rules {
		template := alertTemplate(rule.MonitorMetrics, rule.AlarmThreshold, rule.AlertLabel)
		alertRule := &model.AlertRule{
			AlertName:  rule.AlertName,
			Expression: template.Expression,
			Forsearch:  rule.Forsearch + "s",
			Labels: struct {
				Severity string "yaml:\"severity\" json:\"severity\""
				Metric   string "yaml:\"metric\" json:\"metric\""
			}{
				Severity: rule.Severity,
				Metric:   rule.MonitorMetrics,
			},
			Annotations: struct {
				Description string "yaml:\"description\" json:\"description\""
				Summary     string "yaml:\"summary\" json:\"summary\""
			}{
				Summary:     template.Summary,
				Description: template.Description,
			},
		}
		alertRules = append(alertRules, *alertRule)
	}

	// 添加的规则
	template := alertTemplate(alert.MonitorMetrics, alert.AlarmThreshold, label)
	alertRule := &model.AlertRule{
		AlertName:  alert.AlertName,
		Expression: template.Expression,
		Forsearch:  alert.Forsearch + "s",
		Labels: struct {
			Severity string "yaml:\"severity\" json:\"severity\""
			Metric   string "yaml:\"metric\" json:\"metric\""
		}{
			Severity: alert.Severity,
			Metric:   alert.MonitorMetrics,
		},
		Annotations: struct {
			Description string "yaml:\"description\" json:\"description\""
			Summary     string "yaml:\"summary\" json:\"summary\""
		}{
			Summary:     template.Summary,
			Description: template.Description,
		},
	}
	alertRules = append(alertRules, *alertRule)

	alerts := &model.AlertRules{
		RuleName:   "监控规则",
		AlertRules: alertRules,
	}
	Yml := &model.AlertRuleYaml{
		Groups: []model.AlertRules{*alerts},
	}
	return Yml, nil
}

func deleteRuleDataJoinToYaml(id string) (*model.AlertRuleYaml, error) {
	rules, err := dao.QueryRulesNotIncludedId(id)
	if err != nil {
		return &model.AlertRuleYaml{}, err
	}
	var alertRules []model.AlertRule
	for _, rule := range rules {
		template := alertTemplate(rule.MonitorMetrics, rule.AlarmThreshold, rule.AlertLabel)
		alertRule := &model.AlertRule{
			AlertName:  rule.AlertName,
			Expression: template.Expression,
			Forsearch:  rule.Forsearch + "s",
			Labels: struct {
				Severity string "yaml:\"severity\" json:\"severity\""
				Metric   string "yaml:\"metric\" json:\"metric\""
			}{
				Severity: rule.Severity,
				Metric:   rule.MonitorMetrics,
			},
			Annotations: struct {
				Description string "yaml:\"description\" json:\"description\""
				Summary     string "yaml:\"summary\" json:\"summary\""
			}{
				Summary:     template.Summary,
				Description: template.Description,
			},
		}
		alertRules = append(alertRules, *alertRule)
	}
	alerts := &model.AlertRules{
		RuleName:   "监控规则",
		AlertRules: alertRules,
	}
	Yml := &model.AlertRuleYaml{
		Groups: []model.AlertRules{*alerts},
	}
	return Yml, nil
}
func updateRuleDataJoinToYaml(alert *model.Rule) (*model.AlertRuleYaml, error) {
	rules, err := dao.QueryRules()
	if err != nil {
		return &model.AlertRuleYaml{}, err
	}
	var alertRules []model.AlertRule
	for _, rule := range rules {
		if rule.ID != alert.ID {
			template := alertTemplate(rule.MonitorMetrics, rule.AlarmThreshold, rule.AlertLabel)
			alertRule := &model.AlertRule{
				AlertName:  rule.AlertName,
				Expression: template.Expression,
				Forsearch:  rule.Forsearch + "s",
				Labels: struct {
					Severity string "yaml:\"severity\" json:\"severity\""
					Metric   string "yaml:\"metric\" json:\"metric\""
				}{
					Severity: rule.Severity,
					Metric:   rule.MonitorMetrics,
				},
				Annotations: struct {
					Description string "yaml:\"description\" json:\"description\""
					Summary     string "yaml:\"summary\" json:\"summary\""
				}{
					Summary:     template.Summary,
					Description: template.Description,
				},
			}
			alertRules = append(alertRules, *alertRule)
		} else {
			template := alertTemplate(alert.MonitorMetrics, alert.AlarmThreshold, alert.AlertLabel)
			alertRule := &model.AlertRule{
				AlertName:  alert.AlertName,
				Expression: template.Expression,
				Forsearch:  alert.Forsearch + "s",
				Labels: struct {
					Severity string "yaml:\"severity\" json:\"severity\""
					Metric   string "yaml:\"metric\" json:\"metric\""
				}{
					Severity: alert.Severity,
					Metric:   alert.MonitorMetrics,
				},
				Annotations: struct {
					Description string "yaml:\"description\" json:\"description\""
					Summary     string "yaml:\"summary\" json:\"summary\""
				}{
					Summary:     template.Summary,
					Description: template.Description,
				},
			}
			alertRules = append(alertRules, *alertRule)
		}
	}
	alerts := &model.AlertRules{
		RuleName:   "监控规则",
		AlertRules: alertRules,
	}
	Yml := &model.AlertRuleYaml{
		Groups: []model.AlertRules{*alerts},
	}
	return Yml, nil
}
