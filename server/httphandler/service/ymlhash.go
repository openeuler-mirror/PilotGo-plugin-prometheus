package service

import (
	"openeuler.org/PilotGo/prometheus-plugin/utils"
)

const ymlfile = "./.prometheus-yml.data"

func CheckYMLHash() (bool, error) {

	if !utils.IsFileExist(ymlfile) {
		err := reset()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	savedContent, err := load()
	if err != nil {
		return false, err
	}

	currentContent, err := utils.FileReadString(GlobalPrometheusYml)
	if err != nil {
		return false, err
	}
	return currentContent != savedContent, nil
}

func reset() error {
	bs, err := utils.FileReadString(GlobalPrometheusYml)
	if err != nil {
		return err
	}
	err = utils.FileSaveString(ymlfile, bs)
	if err != nil {
		return err
	}
	return nil
}

func load() (string, error) {
	data, err := utils.FileReadString(ymlfile)
	if err != nil {
		return "", err
	}

	return data, nil
}
