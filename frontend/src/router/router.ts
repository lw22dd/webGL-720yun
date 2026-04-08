import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import index from '@/views/index.vue'
import userCenter from '@/views/userCenter.vue'
import adminLayout from '@/views/admin/adminLayout.vue'
import userDetail from '@/views/userDetail.vue'
import userManagement from '@/views/admin/userManagement.vue'
import spaceManagement from '@/views/admin/spaceManagement.vue'
import ExcelUpload from '@/components/ExcelUpload.vue'
import PanoramaViewer from '@/components/PanoramaViewer.vue'

const routes: RouteRecordRaw[] = [
  { path: '/', component: index, name: 'home' },
  { path: '/user/center', component: userCenter, name: 'userCenter' },
  { path: '/user/:id', component: userDetail, name: 'userDetail' },
  {
    path: '/admin',
    component: adminLayout,
    name: 'admin',
    children: [
      { path: '', redirect: '/admin/users' },
      { path: 'users', component: userManagement, name: 'adminUsers' },
      { path: 'batch-register', component: ExcelUpload, name: 'adminBatchRegister' },
      { path: 'spaces', component: spaceManagement, name: 'adminSpaces' }
    ]
  },
  { path: '/panorama', component: PanoramaViewer, name: 'panorama' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
