import { createRouter, createWebHistory } from 'vue-router'
import index from '@/views/index.vue' 
import userList from '@/views/userList.vue' 
import userDetail from '@/views/userDetail.vue' 

const routes = [
  { path: '/', component: index, name: 'home' },
  // 用户管理路由
  { path: '/user/list', component: userList, name: 'userList' },
  { path: '/user/:id', component: userDetail, name: 'userDetail' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router