import { createRouter, createWebHistory } from 'vue-router'
import index from '@/views/index.vue'
import userCenter from '@/views/userCenter.vue'
import adminLayout from '@/views/adminLayout.vue'
import userDetail from '@/views/userDetail.vue'
import PanoramaViewer from '@/components/PanoramaViewer.vue'

const routes = [
  { path: '/', component: index, name: 'home' },
  { path: '/user/center', component: userCenter, name: 'userCenter' },
  { path: '/user/:id', component: userDetail, name: 'userDetail' },
  {
    path: '/admin',
    component: adminLayout,
    name: 'admin',
    children: [
      { path: '', redirect: '/admin/users' },
      { path: 'users', name: 'adminUsers' },
      { path: 'batch-register', name: 'adminBatchRegister' }
    ]
  },
  { path: '/panorama', component: PanoramaViewer, name: 'panorama' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
