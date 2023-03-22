import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [{
  path: '/login',
  component: () => import('@/views/login/index'),
  hidden: true
},

{
  path: '/404',
  component: () => import('@/views/404'),
  hidden: true
},

{
  path: '/',
  component: Layout,
  redirect: '/realestate',
  children: [{
    path: 'realestate',
    name: 'Realestate',
    component: () => import('@/views/realestate/list/index'),
    meta: {
      title: '房产信息',
      icon: 'realestate'
    }
  }]
}
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const asyncRoutes = [
  {
    path: '/selling',
    component: Layout,
    redirect: '/selling/all',
    name: 'Selling',
    alwaysShow: true,
    meta: {
      title: '销售',
      icon: 'selling'
    },
    children: [{
      path: 'all',
      name: 'SellingAll',
      component: () => import('@/views/selling/all/index'),
      meta: {
        title: '所有销售',
        icon: 'sellingAll'
      }
    },
    {
      path: 'me',
      name: 'SellingMe',
      component: () => import('@/views/selling/me/index'),
      meta: {
        roles: ['editor'],
        title: '我发起的',
        icon: 'sellingMe'
      }
    }, {
      path: 'buy',
      name: 'SellingBuy',
      component: () => import('@/views/selling/buy/index'),
      meta: {
        roles: ['editor'],
        title: '我购买的',
        icon: 'sellingBuy'
      }
    }
    ]
  },
  {
    path: '/donating',
    component: Layout,
    redirect: '/donating/all',
    name: 'Donating',
    alwaysShow: true,
    meta: {
      title: '捐赠',
      icon: 'donating'
    },
    children: [{
      path: 'all',
      name: 'DonatingAll',
      component: () => import('@/views/donating/all/index'),
      meta: {
        title: '所有捐赠',
        icon: 'donatingAll'
      }
    },
    {
      path: 'donor',
      name: 'DonatingDonor',
      component: () => import('@/views/donating/donor/index'),
      meta: {
        roles: ['editor'],
        title: '我发起的捐赠',
        icon: 'donatingDonor'
      }
    }, {
      path: 'grantee',
      name: 'DonatingGrantee',
      component: () => import('@/views/donating/grantee/index'),
      meta: {
        roles: ['editor'],
        title: '我收到的受赠',
        icon: 'donatingGrantee'
      }
    }
    ]
  },
  {
    path: '/addRealestate',
    component: Layout,
    meta: {
      roles: ['admin']
    },
    children: [{
      path: '/addRealestate',
      name: 'AddRealestate',
      component: () => import('@/views/realestate/add/index'),
      meta: {
        title: '新增房产',
        icon: 'addRealestate'
      }
    }]
  },

  // 404 page must be placed at the end !!!
  {
    path: '*',
    redirect: '/404',
    hidden: true
  }
]

const createRouter = () => new Router({
  base: '/',
  // mode: 'history', // require service support
  scrollBehavior: () => ({
    y: 0
  }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
