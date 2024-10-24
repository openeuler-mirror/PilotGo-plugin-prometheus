package model

type PrometheusTarget struct {
	ID       uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UUID     string `json:"uuid"`
	TargetIP string `json:"targetIp"`
	Port     string `json:"port"`
}

type PrometheusObject struct {
	Targets []string `json:"targets"`
	Labels  struct {
		Group string `json:"group"`
	} `json:"labels"`
}
