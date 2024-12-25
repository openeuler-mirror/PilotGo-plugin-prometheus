/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jul 26 16:42:38 2023 +0800
 */
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
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


let baseRoute = '';
if (window.__MICRO_APP_ENVIRONMENT__) {
  console.log('在微前端环境中,baseRoute:',baseRoute)
  baseRoute = window.__MICRO_APP_BASE_ROUTE__ || '/';
} else {
  baseRoute = import.meta.env.BASE_URL;
}
const router = createRouter({
  history: createWebHistory(baseRoute),
  routes
})

export default router
