<template>
  <div class="log-management">
    <t-card :bordered="false">
      <t-tabs v-model="activeTab">
        <t-tab-panel value="operation" label="操作日志">
          <div class="log-filter">
            <t-input v-model="operationFilter.keyword" placeholder="搜索操作内容" style="width: 200px;">
              <template #prefix-icon><SearchIcon /></template>
            </t-input>
            <t-select v-model="operationFilter.module" placeholder="操作模块" clearable style="width: 150px;">
              <t-option value="user" label="用户管理" />
              <t-option value="space" label="空间管理" />
              <t-option value="scene" label="场景管理" />
              <t-option value="hotspot" label="热点管理" />
            </t-select>
            <t-date-range-picker v-model="operationFilter.dateRange" clearable />
            <t-button theme="primary" @click="loadOperationLogs">查询</t-button>
            <t-button theme="default" @click="exportLogs('operation')">导出</t-button>
          </div>
          
          <t-table
            :data="operationLogs"
            :columns="operationColumns"
            :loading="operationLoading"
            :pagination="operationPagination"
            @page-change="onOperationPageChange"
            row-key="id"
            stripe
            hover
          />
        </t-tab-panel>

        <t-tab-panel value="login" label="登录日志">
          <div class="log-filter">
            <t-input v-model="loginFilter.keyword" placeholder="搜索用户名/IP" style="width: 200px;">
              <template #prefix-icon><SearchIcon /></template>
            </t-input>
            <t-select v-model="loginFilter.status" placeholder="登录状态" clearable style="width: 150px;">
              <t-option value="success" label="成功" />
              <t-option value="failed" label="失败" />
            </t-select>
            <t-date-range-picker v-model="loginFilter.dateRange" clearable />
            <t-button theme="primary" @click="loadLoginLogs">查询</t-button>
            <t-button theme="default" @click="exportLogs('login')">导出</t-button>
          </div>
          
          <t-table
            :data="loginLogs"
            :columns="loginColumns"
            :loading="loginLoading"
            :pagination="loginPagination"
            @page-change="onLoginPageChange"
            row-key="id"
            stripe
            hover
          />
        </t-tab-panel>

        <t-tab-panel value="error" label="错误日志">
          <div class="log-filter">
            <t-input v-model="errorFilter.keyword" placeholder="搜索错误信息" style="width: 200px;">
              <template #prefix-icon><SearchIcon /></template>
            </t-input>
            <t-select v-model="errorFilter.level" placeholder="错误级别" clearable style="width: 150px;">
              <t-option value="error" label="Error" />
              <t-option value="warning" label="Warning" />
              <t-option value="critical" label="Critical" />
            </t-select>
            <t-date-range-picker v-model="errorFilter.dateRange" clearable />
            <t-button theme="primary" @click="loadErrorLogs">查询</t-button>
            <t-button theme="default" @click="exportLogs('error')">导出</t-button>
          </div>
          
          <t-table
            :data="errorLogs"
            :columns="errorColumns"
            :loading="errorLoading"
            :pagination="errorPagination"
            @page-change="onErrorPageChange"
            row-key="id"
            stripe
            hover
          />
        </t-tab-panel>
      </t-tabs>
    </t-card>

    <t-dialog
      v-model:visible="detailVisible"
      header="日志详情"
      width="700px"
      :footer="false"
    >
      <div class="log-detail">
        <t-descriptions :column="1" bordered>
          <t-descriptions-item label="时间">{{ currentLog?.time }}</t-descriptions-item>
          <t-descriptions-item label="用户">{{ currentLog?.username }}</t-descriptions-item>
          <t-descriptions-item label="IP地址">{{ currentLog?.ip }}</t-descriptions-item>
          <t-descriptions-item label="操作">{{ currentLog?.action }}</t-descriptions-item>
          <t-descriptions-item label="详情">
            <pre class="detail-content">{{ currentLog?.detail }}</pre>
          </t-descriptions-item>
        </t-descriptions>
      </div>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { SearchIcon } from 'tdesign-icons-vue-next'
import { MessagePlugin } from 'tdesign-vue-next'

const activeTab = ref('operation')
const detailVisible = ref(false)
const currentLog = ref<any>(null)

const operationLoading = ref(false)
const loginLoading = ref(false)
const errorLoading = ref(false)

const operationFilter = reactive({
  keyword: '',
  module: '',
  dateRange: []
})

const loginFilter = reactive({
  keyword: '',
  status: '',
  dateRange: []
})

const errorFilter = reactive({
  keyword: '',
  level: '',
  dateRange: []
})

const operationPagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const loginPagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const errorPagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const operationColumns = [
  { colKey: 'id', title: 'ID', width: 80 },
  { colKey: 'time', title: '时间', width: 180 },
  { colKey: 'username', title: '操作用户', width: 120 },
  { colKey: 'module', title: '模块', width: 100 },
  { colKey: 'action', title: '操作', width: 150 },
  { colKey: 'ip', title: 'IP地址', width: 140 },
  { 
    colKey: 'status', 
    title: '状态', 
    width: 80,
    cell: (h: any, { row }: any) => {
      return h('span', {
        class: row.status === 'success' ? 'status-success' : 'status-failed'
      }, row.status === 'success' ? '成功' : '失败')
    }
  },
  { 
    colKey: 'operation', 
    title: '操作', 
    width: 80,
    cell: (h: any, { row }: any) => {
      return h('a', { onClick: () => showDetail(row) }, '详情')
    }
  }
]

const loginColumns = [
  { colKey: 'id', title: 'ID', width: 80 },
  { colKey: 'time', title: '时间', width: 180 },
  { colKey: 'username', title: '用户名', width: 120 },
  { colKey: 'ip', title: 'IP地址', width: 140 },
  { colKey: 'location', title: '登录地点', width: 150 },
  { colKey: 'browser', title: '浏览器', width: 120 },
  { colKey: 'os', title: '操作系统', width: 100 },
  { 
    colKey: 'status', 
    title: '状态', 
    width: 80,
    cell: (h: any, { row }: any) => {
      return h('span', {
        class: row.status === 'success' ? 'status-success' : 'status-failed'
      }, row.status === 'success' ? '成功' : '失败')
    }
  }
]

const errorColumns = [
  { colKey: 'id', title: 'ID', width: 80 },
  { colKey: 'time', title: '时间', width: 180 },
  { colKey: 'level', title: '级别', width: 100 },
  { colKey: 'message', title: '错误信息', ellipsis: true },
  { colKey: 'file', title: '文件', width: 200, ellipsis: true },
  { colKey: 'line', title: '行号', width: 80 },
  { 
    colKey: 'operation', 
    title: '操作', 
    width: 80,
    cell: (h: any, { row }: any) => {
      return h('a', { onClick: () => showDetail(row) }, '详情')
    }
  }
]

const operationLogs = ref<any[]>([])
const loginLogs = ref<any[]>([])
const errorLogs = ref<any[]>([])

const generateMockOperationLogs = () => {
  const modules = ['user', 'space', 'scene', 'hotspot', 'system']
  const actions = ['创建', '更新', '删除', '查看', '导出']
  const users = ['admin', 'teacher1', 'student1']
  
  return Array.from({ length: 10 }, (_, i) => ({
    id: (operationPagination.current - 1) * 10 + i + 1,
    time: new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000).toLocaleString(),
    username: users[Math.floor(Math.random() * users.length)],
    module: modules[Math.floor(Math.random() * modules.length)],
    action: actions[Math.floor(Math.random() * actions.length)],
    ip: `192.168.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}`,
    status: Math.random() > 0.1 ? 'success' : 'failed',
    detail: '操作详情信息...'
  }))
}

const generateMockLoginLogs = () => {
  const users = ['admin', 'teacher1', 'student1', 'student2']
  const browsers = ['Chrome', 'Firefox', 'Safari', 'Edge']
  const osList = ['Windows 10', 'macOS', 'Linux', 'Android']
  
  return Array.from({ length: 10 }, (_, i) => ({
    id: (loginPagination.current - 1) * 10 + i + 1,
    time: new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000).toLocaleString(),
    username: users[Math.floor(Math.random() * users.length)],
    ip: `192.168.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}`,
    location: '北京市',
    browser: browsers[Math.floor(Math.random() * browsers.length)],
    os: osList[Math.floor(Math.random() * osList.length)],
    status: Math.random() > 0.2 ? 'success' : 'failed'
  }))
}

const generateMockErrorLogs = () => {
  const levels = ['error', 'warning', 'critical']
  const messages = [
    '数据库连接超时',
    '文件上传失败',
    'API请求超时',
    '内存溢出',
    '权限验证失败'
  ]
  
  return Array.from({ length: 10 }, (_, i) => ({
    id: (errorPagination.current - 1) * 10 + i + 1,
    time: new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000).toLocaleString(),
    level: levels[Math.floor(Math.random() * levels.length)],
    message: messages[Math.floor(Math.random() * messages.length)],
    file: `/app/src/handlers/${['user', 'space', 'scene'][Math.floor(Math.random() * 3)]}.go`,
    line: Math.floor(Math.random() * 500) + 1,
    detail: '详细错误堆栈信息...'
  }))
}

const loadOperationLogs = () => {
  operationLoading.value = true
  setTimeout(() => {
    operationLogs.value = generateMockOperationLogs()
    operationPagination.total = 100
    operationLoading.value = false
  }, 500)
}

const loadLoginLogs = () => {
  loginLoading.value = true
  setTimeout(() => {
    loginLogs.value = generateMockLoginLogs()
    loginPagination.total = 100
    loginLoading.value = false
  }, 500)
}

const loadErrorLogs = () => {
  errorLoading.value = true
  setTimeout(() => {
    errorLogs.value = generateMockErrorLogs()
    errorPagination.total = 100
    errorLoading.value = false
  }, 500)
}

const onOperationPageChange = (pageInfo: any) => {
  operationPagination.current = pageInfo.current
  loadOperationLogs()
}

const onLoginPageChange = (pageInfo: any) => {
  loginPagination.current = pageInfo.current
  loadLoginLogs()
}

const onErrorPageChange = (pageInfo: any) => {
  errorPagination.current = pageInfo.current
  loadErrorLogs()
}

const showDetail = (log: any) => {
  currentLog.value = log
  detailVisible.value = true
}

const exportLogs = (type: string) => {
  MessagePlugin.success(`${type === 'operation' ? '操作' : type === 'login' ? '登录' : '错误'}日志导出成功`)
}

onMounted(() => {
  loadOperationLogs()
})
</script>

<style scoped>
.log-management {
  padding: 24px;
}

.log-filter {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.status-success {
  color: #10B981;
}

.status-failed {
  color: #EF4444;
}

.log-detail {
  padding: 16px 0;
}

.detail-content {
  background: #f5f5f5;
  padding: 12px;
  border-radius: 4px;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 300px;
  overflow-y: auto;
  margin: 0;
}
</style>
