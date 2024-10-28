import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import RuleList from '../views/ruleList.vue'

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
]

const router = createRouter({
  history: createWebHashHistory('/plugin/prometheus'),
  routes
})

export default router
