/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Oct 23 16:52:16 2024 +0800
 */
package initprometheus

import (
	"openeuler.org/PilotGo/prometheus-plugin/server/utils"
)

const savedYml = "./.prometheus-yml.data"

func checkPrometheusYmlConsistency(httpaddr string) (bool, error) {

	if !utils.IsFileExist(savedYml) {
		err := resetPrometheusYml(httpaddr)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	savedContent, err := loadPrometheusYml()
	if err != nil {
		return false, err
	}

	currentContent, err := utils.FileReadString(PrometheusYml)
	if err != nil {
		return false, err
	}
	return currentContent != savedContent, nil
}

func resetPrometheusYml(httpaddr string) error {
	if err := initYml(httpaddr); err != nil {
		return err
	}
	bs, err := utils.FileReadString(PrometheusYml)
	if err != nil {
		return err
	}
	err = utils.FileSaveString(savedYml, bs)
	if err != nil {
		return err
	}
	return nil
}

func loadPrometheusYml() (string, error) {
	data, err := utils.FileReadString(savedYml)
	if err != nil {
		return "", err
	}

	return data, nil
}
