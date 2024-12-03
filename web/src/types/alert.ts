/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: wanghaohao <wanghaohao@kylinos.cn>
 * Date: Mon Oct 28 17:43:48 2024 +0800
 */
export interface Alert {
  id: number;
  ip: string;
  alertName: string;
  level: string;
  alertTime: string;
  alertEndTime: string;
  handleState: string;
  summary: string;
  description: string;
}