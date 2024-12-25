/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jul 26 16:42:38 2023 +0800
 */
import request from './request'
// 获取指标列表
export function getPromRules() {
  return request({
    url: '/plugin/prometheus/api/targets',
    method: 'get',
  })
}

// 获取prome某一时间点的数据
export function getPromeCurrent(data: object) {
  return request({
    url: '/plugin/prometheus/api/query',
    method: 'get',
    params: data
  })
}

// 获取prome某一时间段的数据
export function getPromeRange(data: object) {
  return request({
    url: '/plugin/prometheus/api/query_range',
    method: 'get',
    params: data
  })
}

// 获取监控主机ip
export function getMacIp() {
  return request({
    url: '/plugin_manage/info',
    method: 'get',
  })
}

// 获取全部配置规则
export function getRuleList(data: Object) {
  return request({
    url: "/plugin/prometheus/api/ruleQuery",
    method: "get",
    params: data,
  });
}
// 添加告警告警规则
export function addConfigRule(data: object) {
  return request({
    url: "/plugin/prometheus/api/ruleAdd",
    method: "post",
    data,
  });
}

// 编辑告警告警规则
export function updateConfigRule(data: object) {
  return request({
    url: " /plugin/prometheus/api/ruleUpdate",
    method: "post",
    data,
  });
}


// 删除告警告警规则
export function delConfigRule(data: { id: number }) {
  return request({
    url: "/plugin/prometheus/api/ruleDelete",
    method: "get",
    params: data,
  });
}

// 获取所有监控指标
export function getMetrics() {
  return request({
    url: "/plugin/prometheus/api/ruleMetrics",
    method: "get",
  });
}

// 已安装监控组件的主机列表分页
export function getExporterList(data: object) {
  return request({
    url: "/plugin/prometheus/api/monitorlist",
    method: "get",
    params: data,
  });
}

// 获取所有历史告警
export function getHistoryAlerts(data: Object) {
  return request({
    url: "/plugin/prometheus/api/alertQuery",
    method: "get",
    params: data,
  });
}

// 变更告警状态
export function updateAlertState(data: Object) {
  return request({
    url: "/plugin/prometheus/api/alertUpdateState",
    method: "post",
    data,
  });
}