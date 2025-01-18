import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('../views/Home.vue'),
    },
    {
      path: '/realty-agency',
      component: () => import('../views/RealtyAgency.vue'),
    },
    {
      path: '/trading-platform',
      component: () => import('../views/TradingPlatform.vue'),
    },
    {
      path: '/bank',
      component: () => import('../views/Bank.vue'),
    },
  ],
})

export default router 