package model

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
