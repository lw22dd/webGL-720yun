import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/userStore'
import index from '@/views/index.vue'
import Login from '@/views/Login.vue'
import userCenter from '@/views/userCenter.vue'
import adminLayout from '@/views/admin/adminLayout.vue'
import userDetail from '@/views/userDetail.vue'
import userManagement from '@/views/admin/userManagement.vue'
import spaceManagement from '@/views/admin/spaceManagement.vue'
import sceneManagement from '@/views/admin/sceneManagement.vue'
import SceneGraphEditor from '@/views/admin/sceneGraphEditor.vue'
import systemSettings from '@/views/admin/systemSettings.vue'
import logManagement from '@/views/admin/logManagement.vue'
import favorites from '@/views/favorites.vue'
import history from '@/views/history.vue'
import ExcelUpload from '@/components/ExcelUpload.vue'
import PanoramaViewer from '@/components/PanoramaViewer.vue'

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/login' },
  { path: '/index', component: index, name: 'home', meta: { requiresAuth: true } },
  { path: '/login', component: Login, name: 'login' },
  { path: '/user/center', component: userCenter, name: 'userCenter' },
  { path: '/user/:id', component: userDetail, name: 'userDetail' },
  { path: '/favorites', component: favorites, name: 'favorites' },
  { path: '/history', component: history, name: 'history' },
  {
    path: '/admin',
    component: adminLayout,
    name: 'admin',
    children: [
      { path: '', redirect: '/admin/users' },
      { path: 'users', component: userManagement, name: 'adminUsers' },
      { path: 'batch-register', component: ExcelUpload, name: 'adminBatchRegister' },
      { path: 'spaces', component: spaceManagement, name: 'adminSpaces' },
      { path: 'spaces/:id/graph', component: SceneGraphEditor, name: 'spaceGraph' },
      { path: 'scenes', component: sceneManagement, name: 'adminScenes' },
      { path: 'settings', component: systemSettings, name: 'adminSettings' },
      { path: 'logs', component: logManagement, name: 'adminLogs' }
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
