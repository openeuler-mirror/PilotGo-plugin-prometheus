/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: wanghaohao <wanghaohao@kylinos.cn>
 * Date: Mon Oct 28 17:58:01 2024 +0800
 */
export interface Host {
  hostId: number,
  targetIp: string,
  deptName: string,
  agentStatus: string,
  operatingSystem: string,
  monitorAgentStatus: string,
  version: string,
  monitorVersion: string,
  architecture: string,
  registTime: string
}