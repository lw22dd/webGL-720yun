<template>
  <div class="space-management">
    <AdminTable
      ref="tableRef"
      :data="spaceList"
      :columns="columns"
      :pagination="paginationConfig"
      delete-item-name="空间"
      @page-change="handlePageChange"
      @delete="handleDeleteSpace"
    >
      <template #operations="{ row }">
        <t-space>
          <t-button
            theme="primary"
            variant="text"
            size="small"
            @click="handleEditSpace(row)"
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

    <FormDialog
      ref="formDialogRef"
      :form-fields="formFields"
      :default-data="defaultFormData"
      @submit="handleSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import AdminTable from '@/components/admin/AdminTable.vue'
import FormDialog, { type FormField } from '@/components/admin/FormDialog.vue'
import SpaceApi from '@/apis/spaceApi'
import type { SpaceListItem } from '@/models/SpaceModel'

const searchKeyword = ref('')
const spaceList = ref<SpaceListItem[]>([])

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

const tableRef = ref()
const formDialogRef = ref()
const submitLoading = ref(false)

const formFields: FormField[] = [
  {
    name: 'name',
    label: '空间名称',
    type: 'input',
    required: true,
    placeholder: '请输入空间名称'
  },
  {
    name: 'slug',
    label: '空间标识',
    type: 'input',
    required: true,
    placeholder: '请输入空间标识(英文唯一)'
  },
  {
    name: 'description',
    label: '描述',
    type: 'textarea',
    placeholder: '请输入空间描述'
  },
  {
    name: 'province',
    label: '省份',
    type: 'input',
    placeholder: '请输入省份'
  },
  {
    name: 'city',
    label: '城市',
    type: 'input',
    placeholder: '请输入城市'
  },
  {
    name: 'longitude',
    label: '经度',
    type: 'input',
    placeholder: '请输入经度'
  },
  {
    name: 'latitude',
    label: '纬度',
    type: 'input',
    placeholder: '请输入纬度'
  },
  {
    name: 'zoom_level',
    label: '缩放级别',
    type: 'input',
    placeholder: '请输入缩放级别'
  },
  {
    name: 'sort_order',
    label: '排序',
    type: 'input',
    placeholder: '请输入排序值'
  },
  {
    name: 'status',
    label: '状态',
    type: 'switch'
  }
]

const defaultFormData = {
  status: 1
}

const columns = computed(() => [
  {
    colKey: 'id',
    title: 'ID',
    width: 80,
    align: 'center' as const,
    fixed: 'left' as const
  },
  {
    colKey: 'name',
    title: '空间名称',
    minWidth: 150,
    ellipsis: true
  },
  {
    colKey: 'slug',
    title: '标识',
    minWidth: 120,
    ellipsis: true
  },
  {
    colKey: 'description',
    title: '描述',
    minWidth: 180,
    ellipsis: true
  },
  {
    colKey: 'province',
    title: '省份',
    width: 100,
    ellipsis: true
  },
  {
    colKey: 'city',
    title: '城市',
    width: 100,
    ellipsis: true
  },
  {
    colKey: 'scene_count',
    title: '场景数',
    width: 80,
    align: 'center' as const
  },
  {
    colKey: 'status',
    title: '状态',
    width: 80,
    align: 'center' as const,
    cell: ({ row }: { row?: any }) => {
      if (!row) return '-'
      return row.status === 1 ? '启用' : '禁用'
    }
  },
  {
    colKey: 'created_at',
    title: '创建时间',
    width: 180,
    ellipsis: true
  },
  {
    colKey: 'operations',
    title: '操作',
    width: 140,
    align: 'center' as const,
    fixed: 'right' as const
  }
])

const loadSpaceList = async () => {
  try {
    const result = await SpaceApi.getSpaceList({
      keyword: searchKeyword.value,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    if (result.code === 200 && result.data) {
      spaceList.value = result.data.spaces
      pagination.total = result.data.page_info.total
      paginationConfig.total = result.data.page_info.total
    } else {
      MessagePlugin.error(result.msg || '获取空间列表失败')
    }
  } catch (error) {
    MessagePlugin.error('获取空间列表失败')
  }
}

const handleAddSpace = () => {
  formDialogRef.value?.openAddDialog()
}

const handleEditSpace = (space: SpaceListItem) => {
  formDialogRef.value?.openEditDialog(
    {
      name: space.name,
      slug: space.slug,
      description: space.description,
      province: space.province,
      city: space.city,
      longitude: space.longitude,
      latitude: space.latitude,
      zoom_level: space.zoom_level,
      sort_order: space.sort_order,
      status: space.status
    },
    space.id
  )
}

const handleSubmit = async (data: Record<string, any>) => {
  submitLoading.value = true
  try {
    let result
    const formDataObj = new FormData()
    Object.keys(data).forEach(key => {
      if (key !== 'id' && data[key] !== undefined && data[key] !== '') {
        formDataObj.append(key, data[key])
      }
    })

    if (data.id) {
      result = await SpaceApi.updateSpace(Number(data.id), formDataObj)
    } else {
      result = await SpaceApi.createSpace(formDataObj)
    }

    if (result.code === 200) {
      MessagePlugin.success(data.id ? '编辑空间成功' : '新增空间成功')
      formDialogRef.value?.closeDialog()
      loadSpaceList()
    } else {
      MessagePlugin.error(result.msg || '操作失败')
    }
  } catch (error) {
    MessagePlugin.error('操作失败')
  } finally {
    submitLoading.value = false
  }
}

const openDeleteDialog = (id: number) => {
  tableRef.value?.openDeleteDialog(id)
}

const handleDeleteSpace = async (id: string | number) => {
  try {
    const result = await SpaceApi.deleteSpace(Number(id))
    if (result.code === 200) {
      MessagePlugin.success('删除成功')
      loadSpaceList()
    } else {
      MessagePlugin.error(result.msg || '删除失败')
    }
  } catch (error) {
    MessagePlugin.error('删除失败')
  }
}

const handlePageChange = (context: { current: number; pageSize: number }) => {
  pagination.currentPage = context.current
  paginationConfig.current = context.current
  pagination.pageSize = context.pageSize
  paginationConfig.pageSize = context.pageSize
  loadSpaceList()
}

onMounted(() => {
  loadSpaceList()
})

defineExpose({
  handleAddSpace
})
</script>

<style scoped>
.space-management {
  width: 100%;
  height: 100%;
  min-height: 100%;
  display: flex;
  flex-direction: column;
}
</style>
