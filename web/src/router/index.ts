/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jul 26 16:42:38 2023 +0800
 */
import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import RuleList from '../views/ruleList.vue'
import AlertList from '../views/alertList.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'home',
    component: HomeView
  },
  {
    path: '/rule',
    name: 'rule',
    component: RuleList
  },
  {
    path: '/alert',
    name: 'alert',
    component: AlertList
  },
]

const router = createRouter({
  history: createWebHashHistory('/plugin/prometheus'),
  routes
})

export default router
