package model

type PrometheusTarget struct {
	ID       uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UUID     string `json:"uuid"`
	TargetIP string `json:"targetIp"`
	Port     string `json:"port"`
	ID_idx   string `gorm:"index"`
}
