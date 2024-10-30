package model

type Alert struct {
	ID           int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	AlertName    string `gorm:"uniqueIndex:idx_alerts, length:100" json:"alertName"`
	GroupId      string `gorm:"uniqueIndex:idx_alerts, length:100" json:"groupId"`
	IP           string `gorm:"uniqueIndex:idx_alerts, length:100" json:"ip"`
	AlertLevel   string `json:"level"`
	Summary      string `json:"summary"`
	Metric       string `json:"metric"`
	Description  string `json:"description"`
	AlertTime    string `json:"alertTime"`
	AlertEndTime string `gorm:"uniqueIndex:idx_alerts, length:100" json:"alertEndTime"`
	ConfirmTime  string `json:"confirmTime"`
	CompleteTime string `json:"completeTime"`
	HandleState  string `json:"handleState"`
}

type AlertsResponse struct {
	Status string `json:"status"`
	Data   struct {
		Alerts []AlertResponse `json:"alerts"`
	} `json:"data"`
}

type AlertResponse struct {
	Labels struct {
		AlertName      string `json:"alertname"`
		Group          string `json:"group"`
		UserName       string `json:"userName"`
		DepartmentName string `json:"departmentName"`
		Instance       string `json:"instance"`
		Job            string `json:"job"`
		Severity       string `json:"severity"`
		Metric         string `json:"metric"`
	} `json:"labels"`
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
	State    string `json:"state"`
	ActiveAt string `json:"activeAt"`
	Value    string `json:"value"`
}
