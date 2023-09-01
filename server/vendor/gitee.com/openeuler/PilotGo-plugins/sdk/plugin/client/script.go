package client

import (
	"encoding/base64"
	"encoding/json"

	"gitee.com/openeuler/PilotGo-plugins/sdk/common"
	"gitee.com/openeuler/PilotGo-plugins/sdk/utils/httputils"
)

type CmdResult struct {
	MachineUUID string
	MachineIP   string
	RetCode     int
	Stdout      string
	Stderr      string
}

func (c *Client) RunCommand(batch *common.Batch, cmd string) ([]*CmdResult, error) {
	url := c.Server + "/api/v1/pluginapi/run_command"

	p := &struct {
		Batch   *common.Batch `json:"batch"`
		Command string        `json:"command"`
	}{
		Batch:   batch,
		Command: base64.StdEncoding.EncodeToString([]byte(cmd)),
	}

	r, err := httputils.Post(url, &httputils.Params{
		Body: p,
	})
	if err != nil {
		return nil, err
	}

	res := &struct {
		Code    int          `json:"code"`
		Message string       `json:"msg"`
		Data    []*CmdResult `json:"data"`
	}{}
	if err := json.Unmarshal(r.Body, res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (c *Client) RunScript(batch *common.Batch, script string, params []string) ([]*CmdResult, error) {
	url := c.Server + "/api/v1/pluginapi/run_script"

	p := &struct {
		Batch  *common.Batch `json:"batch"`
		Script string        `json:"script"`
		Params []string      `json:"params"`
	}{
		Batch:  batch,
		Script: base64.StdEncoding.EncodeToString([]byte(script)),
		Params: params,
	}

	r, err := httputils.Post(url, &httputils.Params{
		Body: p,
	})
	if err != nil {
		return nil, err
	}

	res := &struct {
		Code    int          `json:"code"`
		Message string       `json:"msg"`
		Data    []*CmdResult `json:"data"`
	}{}
	if err := json.Unmarshal(r.Body, res); err != nil {
		return nil, err
	}

	return res.Data, nil
}
