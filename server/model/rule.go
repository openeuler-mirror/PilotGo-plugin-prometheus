package model

type Rule struct {
	ID             int          `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	AlertName      string       `gorm:"uniqueIndex:idx_alert, length:100" json:"alertName"`                 // 告警名称
	CustomDesc     string       `gorm:"uniqueIndex:idx_alert, length:100" json:"desc"`                      // 自定义描述
	MonitorMetrics string       `gorm:"uniqueIndex:idx_alert, length:100" json:"metrics"`                   // 监控指标
	AlarmThreshold string       `gorm:"uniqueIndex:idx_alert, length:100" json:"threshold"`                 // 告警阈值
	Forsearch      string       `gorm:"uniqueIndex:idx_alert, length:100" json:"duration"`                  // 触发时长
	Severity       string       `gorm:"uniqueIndex:idx_alert, length:100" json:"severity"`                  // 告警级别
	AlertLabel     string       `json:"alertLabel"`                                                         // 告警标签
	AlertTargets   []RuleTarget `gorm:"foreignKey:RuleId;constraint:OnDelete:CASCADE;" json:"alertTargets"` // 告警机器                                                 // 部门名称
}
type RuleTarget struct {
	ID     int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	RuleId int    `gorm:"uniqueIndex:idx_rule_ip" json:"ruleId"`
	IP     string `gorm:"uniqueIndex:idx_rule_ip" json:"ip"`
	UUID   int    `gorm:"uniqueIndex:idx_rule_ip" json:"uuid"`
}

type AlertRuleYaml struct {
	Groups []AlertRules `yaml:"groups" json:"groups"`
}

type AlertRules struct {
	RuleName   string      `yaml:"name" json:"name"`
	AlertRules []AlertRule `yaml:"rules" json:"rules"`
}

type AlertRule struct {
	AlertName  string      `yaml:"alert" json:"name"`   // 告警名称
	Expression string      `yaml:"expr" json:"query"`   // 告警规则表达式
	Forsearch  interface{} `yaml:"for" json:"duration"` // 触发时长
	Labels     struct {
		Severity string `yaml:"severity" json:"severity"` // 告警级别
		Metric   string `yaml:"metric" json:"metric"`     //告警指标
	} `yaml:"labels" json:"labels"`
	Annotations struct {
		Description string `yaml:"description" json:"description"`
		Summary     string `yaml:"summary" json:"summary"`
	} `yaml:"annotations" json:"annotations"` // 告警通知
}
