import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user.store'
import index from '@/views/index.vue'
import Login from '@/views/Login.vue'
import UserCenter from '@/views/user/UserCenter.vue'
import AdminLayout from '@/layouts/AdminLayout.vue'
import UserDetail from '@/views/user/UserDetail.vue'
import UserManagement from '@/views/admin/UserManagement.vue'
import SpaceManagement from '@/views/admin/SpaceManagement.vue'
import SpaceMapEditor from '@/views/admin/SpaceMapEditor.vue'
import SystemSettings from '@/views/admin/SystemSettings.vue'
import LogManagement from '@/views/admin/LogManagement.vue'
import Favorites from '@/views/user/Favorites.vue'
import History from '@/views/user/History.vue'
import ExcelUpload from '@/components/admin/ExcelUpload.vue'
import PanoramaViewer from '@/components/PanoramaViewer.vue'
import SpaceDetail from '@/views/SpaceDetail.vue'

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/login' },
  { path: '/index', component: index, name: 'home', meta: { requiresAuth: true } },
  { path: '/login', component: Login, name: 'login' },
  { path: '/user/center', component: UserCenter, name: 'userCenter' },
  { path: '/user/:id', component: UserDetail, name: 'userDetail' },
  { path: '/favorites', component: Favorites, name: 'favorites' },
  { path: '/history', component: History, name: 'history' },
  { path: '/space/:id', component: SpaceDetail, name: 'spaceDetail', meta: { requiresAuth: true } },
  {
    path: '/admin',
    component: AdminLayout,
    name: 'admin',
    children: [
      { path: '', component: UserManagement, name: 'adminHome' },
      { path: 'users', component: UserManagement, name: 'adminUsers' },
      { path: 'batch-register', component: ExcelUpload, name: 'adminBatchRegister' },
      { path: 'spaces', component: SpaceManagement, name: 'adminSpaces' },
      { path: 'spaces/:id/map', component: SpaceMapEditor, name: 'spaceMapEditor' },
      { path: 'settings', component: SystemSettings, name: 'adminSettings' },
      { path: 'logs', component: LogManagement, name: 'adminLogs' }
    ]
  },
  { path: '/panorama', component: PanoramaViewer, name: 'panorama' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()
  const token = userStore.accessToken

  if (to.meta.requiresAuth && !token) {
    next({ name: 'login' })
  } else if (to.name === 'login' && token) {
    next({ path: '/index' })
  } else {
    next()
  }
})

export default router
