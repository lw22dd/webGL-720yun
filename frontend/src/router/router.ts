import { createRouter, createWebHistory } from 'vue-router'
import index from '@/views/index.vue' 
import userManagement from '@/views/userManagement.vue' 
import userDetail from '@/views/userDetail.vue' 
import PanoramaViewer from '@/components/PanoramaViewer.vue' 

const routes = [
  { path: '/', component: index, name: 'home' },
  // 用户管理路由
  { path: '/user/list', component: userManagement, name: 'userList' },
  { path: '/user/:id', component: userDetail, name: 'userDetail' },
  // 全景图路由
  { path: '/panorama', component: PanoramaViewer, name: 'panorama' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router