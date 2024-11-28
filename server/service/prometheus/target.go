package initprometheus

import (
	"openeuler.org/PilotGo/prometheus-plugin/server/dao"
	"openeuler.org/PilotGo/prometheus-plugin/server/model"
)

func PrometheusTargetsUpdate() ([]model.PrometheusObject, error) {
	var objects []model.PrometheusObject

	dbtarget, err := dao.GetPrometheusTarget()
	if err != nil {
		return objects, err
	}
	objects = append(objects, model.PrometheusObject{Targets: dbtarget})

	rules, err := dao.QueryRules()
	if err != nil {
		return objects, err
	}
	for _, rule := range rules {
		var result []string
		for _, target := range rule.AlertTargets {
			result = append(result, target.IP+":9100")
		}
		ob := model.PrometheusObject{
			Targets: result,
			Labels: struct {
				Group string "json:\"group\""
			}{
				Group: rule.AlertLabel,
			},
		}
		objects = append(objects, ob)
	}

	return objects, nil
}
