<template>
  <div class="user-management">
    <AdminTable
      ref="tableRef"
      :data="userList"
      :columns="columns"
      :pagination="paginationConfig"
      delete-item-name="用户"
      @page-change="handlePageChange"
      @delete="handleDeleteUser"
      @batch-delete="handleBatchDeleteUser"
    >
      <template #title>用户管理</template>
      <template #actionBar>
        <t-button theme="primary" @click="handleAddUser">
          <template #icon><t-icon-plus size="24" color="primary" /></template>
          新建用户
        </t-button>
      </template>
      <template #operations="{ row }">
        <t-space>
          <t-button
            theme="primary"
            variant="text"
            size="small"
            @click="handleEditUser(row)"
          >
            编辑
          </t-button>
          <t-button
            theme="danger"
            variant="text"
            size="small"
            @click="openDeleteDialog(row.id)"
          >
            删除
          </t-button>
        </t-space>
      </template>
    </AdminTable>

    <t-dialog
      v-model:visible="dialogVisible"
      :header="dialogTitle"
      width="500px"
      :footer="false"
    >
      <t-form
        ref="formRef"
        :data="formData"
        :rules="formRules"
        label-width="100px"
        @submit="handleSubmit"
      >
        <!-- 学生角色时显示学号输入 -->
        <t-form-item v-if="formData.role_id === 3" label="学号" name="id">
          <t-input v-model="formData.id" placeholder="请输入学号" />
        </t-form-item>

        <t-form-item label="用户名" name="username">
          <t-input v-model="formData.username" placeholder="请输入用户名" />
        </t-form-item>

        <t-form-item label="密码" name="password">
          <t-input v-model="formData.password" type="password" placeholder="请输入密码" />
        </t-form-item>

        <t-form-item label="角色" name="role_id">
          <t-select v-model="formData.role_id" :options="roleOptions" placeholder="请选择角色" />
        </t-form-item>

        <t-form-item style="margin-top: 20px">
          <t-space>
            <t-button type="submit" theme="primary" :loading="submitLoading">
              提交
            </t-button>
            <t-button @click="dialogVisible = false">取消</t-button>
          </t-space>
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import AdminTable from '@/components/admin/AdminTable.vue'
import UserApi from '@/apis/userApi'
import { User } from '@/models/UserModel'

interface UserItem {
  id: number
  username: string
  email: string
  phone: string
  nickname: string
  role_id: number
  status: number
  [key: string]: any
}

const searchKeyword = ref('')
const userList = ref<UserItem[]>([])

const paginationConfig = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  defaultCurrent: 1,
  defaultPageSize: 10,
  showJumper: true
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('新增用户')
const formRef = ref()
const submitLoading = ref(false)
const tableRef = ref()
const isEditMode = ref(false)

const formData = reactive({
  id: '',
  username: '',
  email: '',
  password: '',
  phone: '',
  nickname: '',
  role_id: 3,
  status: 1
})

const formRules = {
  id: [
    { required: true, message: '请输入学号', trigger: 'blur', type: 'number' }
  ],
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度应在3-20个字符之间', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
  ],
  role_id: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

const roleOptions = [
  { label: '超级管理员', value: 1 },
  { label: '教师', value: 2 },
  { label: '学生', value: 3 }
]

const columns = computed(() => [
  {
    colKey: 'row-select',
    type: 'multiple' as const,
    width: 50
  },
  {
    colKey: 'id',
    title: 'ID',
    width: 80,
    align: 'center' as const,
    fixed: 'left' as const
  },
  {
    colKey: 'username',
    title: '用户名',
    minWidth: 120,
    ellipsis: true
  },
  {
    colKey: 'nickname',
    title: '昵称',
    minWidth: 100,
    ellipsis: true
  },
  {
    colKey: 'email',
    title: '邮箱',
    minWidth: 180,
    ellipsis: true
  },
  {
    colKey: 'role_id',
    title: '角色',
    width: 100,
    align: 'center' as const,
    cell: (h, { row }) => {
      if (!row) return h('span', {}, '-')
      if (row.role && row.role.name) {
        return h('span', {}, row.role.name)
      }
      const role = roleOptions.find(r => r.value === row.role_id)
      return h('span', {}, role ? role.label : '-')
    }
  },
  {
    colKey: 'status',
    title: '状态',
    width: 80,
    align: 'center' as const,
    cell: (h, { row }) => {
      if (!row) return h('span', {}, '-')
      return h('span', {}, row.status === 1 ? '启用' : '禁用')
    }
  },
  {
    colKey: 'operations',
    title: '操作',
    width: 140,
    align: 'center' as const,
    fixed: 'right' as const
  }
])

const loadUserList = async () => {
  try {
    const result = await UserApi.getUserList({
      keyword: searchKeyword.value,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    if (result.code === 200) {
      console.log(result.data.users)
      userList.value = result.data.users
      pagination.total = result.data.total
      paginationConfig.total = result.data.total
    } else {
      MessagePlugin.error(result.msg || '获取用户列表失败')
    }
  } catch (error) {
    MessagePlugin.error('获取用户列表失败')
  }
}

const handleEditUser = (user: UserItem) => {
  dialogTitle.value = '编辑用户'
  isEditMode.value = true
  Object.assign(formData, {
    id: user.id,
    username: user.username,
    email: user.email,
    phone: user.phone || '',
    nickname: user.nickname || '',
    role_id: user.role_id,
    status: user.status
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  console.log(formData)

  const valid = await (formRef.value as any).validate()
  if (valid) {
    submitLoading.value = true
    try {
      let result
      if (isEditMode.value) {
        console.log("进入更新用户:")
        result = await UserApi.updateUser(formData.id.toString(), formData)
      } else {
        const { status, created_at, updated_at, role, id, ...createData } = formData as any;
        if (id && id !== '') {
          createData.id = Number(id);
        }
        result = await UserApi.addUser(createData)
        console.log("新增用户结果:", result)
      }

      if (result.code === 200) {
        MessagePlugin.success(isEditMode.value ? '编辑用户成功' : '新增用户成功')
        dialogVisible.value = false
        loadUserList()
      } else {
        MessagePlugin.error(result.msg || '操作失败')
      }
    } catch (error: any) {
      console.error('提交失败:', error)
      const errorMessage = error?.message || error?.msg || '操作失败'
      MessagePlugin.error(`操作失败: ${errorMessage}`)
    } finally {
      submitLoading.value = false
    }
  }
}

const openDeleteDialog = (id: number) => {
  tableRef.value?.openDeleteDialog(id)
}

const handleAddUser = () => {
  dialogTitle.value = '新增用户'
  isEditMode.value = false
  Object.assign(formData, {
    id: '',
    username: '',
    email: '',
    password: '',
    phone: '',
    nickname: '',
    role_id: 3,
    status: 1
  })
  dialogVisible.value = true
}

const handleDeleteUser = async (id: string | number) => {
  try {
    const result = await UserApi.deleteUser(id.toString())
    if (result.code === 200) {
      MessagePlugin.success('删除成功')
      loadUserList()
    } else {
      MessagePlugin.error(result.msg || '删除失败')
    }
  } catch (error) {
    MessagePlugin.error('删除失败')
  }
}

const handleBatchDeleteUser = async (ids: (string | number)[]) => {
  try {
    const result = await UserApi.deleteUserBatch({ ids: ids.map(id => Number(id)) })
    if (result.code === 200) {
      MessagePlugin.success('批量删除成功')
      loadUserList()
    } else {
      MessagePlugin.error(result.msg || '批量删除失败')
    }
  } catch (error) {
    MessagePlugin.error('批量删除失败')
  }
}

const handlePageChange = (context: any) => {
  pagination.currentPage = context.current
  paginationConfig.current = context.current
  pagination.pageSize = context.pageSize
  paginationConfig.pageSize = context.pageSize
  loadUserList()
}

onMounted(() => {
  loadUserList()
})
</script>

<style scoped>
.user-management {
  width: 100%;
  height: 100%;
  min-height: 100%;
  display: flex;
  flex-direction: column;
}
</style>
