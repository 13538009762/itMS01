import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('../views/Layout.vue'),
    redirect: '/assets',
    children: [
      {
        path: 'assets',
        name: 'AssetList',
        component: () => import('../views/AssetList.vue')
      },
      {
        path: 'applies',
        name: 'ApplyList',
        component: () => import('../views/ApplyList.vue')
      },
      {
        path: 'repairs',
        name: 'RepairList',
        component: () => import('../views/RepairList.vue')
      },
      {
        path: 'users',
        name: 'UserList',
        component: () => import('../views/UserList.vue')
      },
      {
        path: 'my-assets',
        name: 'MyAssets',
        component: () => import('../views/MyAssets.vue')
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('../views/Profile.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.name !== 'Login' && !token) {
    next({ name: 'Login' })
  } else {
    next()
  }
})

export default router
