/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Oct 10 17:53:11 2023 +0800
 */
package utils

import (
	"io"
	"os/exec"
)

func RunCommand(s string) (int, string, string, error) {
	cmd := exec.Command("/bin/bash", "-c", "export LANG=en_US.utf8 ; "+s)

	StdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return 0, "", "", err
	}

	StderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return 0, "", "", err
	}

	exitCode := 0
	err = cmd.Start()
	if err != nil {
		return 0, "", "", err
	}

	b1, err := io.ReadAll(StdoutPipe)
	if err != nil {
		return 0, "", "", err
	}
	stdout := string(b1)

	b2, err := io.ReadAll(StderrPipe)
	if err != nil {
		return 0, "", "", err
	}
	stderr := string(b2)

	err = cmd.Wait()
	if err != nil {
		e, ok := err.(*exec.ExitError)
		if !ok {
			return 0, "", "", err
		}
		exitCode = e.ExitCode()
	}

	return exitCode, stdout, stderr, nil
}
