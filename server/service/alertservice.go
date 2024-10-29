package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"openeuler.org/PilotGo/prometheus-plugin/config"
	"openeuler.org/PilotGo/prometheus-plugin/model"
)

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
