/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: wanghaohao <wanghaohao@kylinos.cn>
 * Date: Mon Oct 28 17:58:01 2024 +0800
 */
export interface ConfigRule {
  id?: number;
  metrics: string;
  ips?: string[];
  batches: any[];
  batches_str: string;
  alertTargets: any[];
  alertHostIds?: string;
  alertLabel: string;
  alertName: string;
  duration: number | string;
  severity: string;
  input_severity?: string;
  threshold: number | string;
  desc: string;
  [key: string]: unknown;
}
