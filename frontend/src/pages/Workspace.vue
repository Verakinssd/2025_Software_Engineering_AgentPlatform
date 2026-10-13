<template>
  <v-container>
    <v-row class="mb-4" style="border-bottom: 1px solid #e0e0e0;">
      <v-col cols="12" class="d-flex justify-space-between align-center">
        <h1 class="text-h4" data-guide="workspace-title">工作空间</h1>
        <v-btn color="primary" @click="dialog = true" data-guide="create-workflow">创建工作流</v-btn>
      </v-col>
    </v-row>


    <v-select
      style="width: 150px; margin-top: 36px;"
      v-model="filterStatus"
      :items="statusOptions"
      item-title="text"
      item-value="value"
      variant="solo"
      density="compact"
    />

    <v-row v-if="filteredWorkflows.length === 0" justify="center">
      <v-col cols="12" class="text-center">
        <v-icon size="48">mdi-inbox</v-icon>
        <div>暂无已发布工作流</div>
      </v-col>
    </v-row>

    <v-row v-else>
      <v-col :xxl="2" :xl="3" :lg="3" :md="4" :sm="6" :xs="12" class="pa-6" v-for="(item, index) in filteredWorkflows" :key="item.id">
        <v-card :key="index"
          height="100%"
        >
          <v-img
            :key="index"
            :aspect-ratio="16/9"
            :src="item.cover || '/images/default_cover.jpg'"
            cover
            lazy-src="/images/default_cover.jpg"
          ></v-img>
          <v-card-title>
            {{ item.name || '未命名' }}
            <p style="font-size: 12px; color: gray" class="mt-1">
              <v-avatar size="16" class="mb-1 mr-1">
                <v-img :src="item.belongsTo.avatar || '/images/default_avatar.jpg'"></v-img>
              </v-avatar>
              {{ item.belongsTo.username }}
            </p>
          </v-card-title>
          <v-card-text style="height: 60px; overflow: hidden; text-overflow: ellipsis; display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical;">
            {{ item.description }}
          </v-card-text>
          <v-divider class="mt-2"></v-divider>
          <v-card-actions>
            <v-btn color="primary" size="small" @click="editWorkflowInfo(item)" prepend-icon="mdi-pencil">简介</v-btn>
            <v-btn color="primary" size="small" @click="editWorkflow(item.id)" prepend-icon="mdi-eye">工作流</v-btn>
            <v-btn color="error" size="small" @click="confirmDelete(item.id)" prepend-icon="mdi-delete">删除</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <v-dialog v-model="dialog" max-width="600px" persistent>
      <v-card>
        <v-card-title>新建工作流</v-card-title>
        <v-card-text>
          <v-form ref="formRef" @submit.prevent="submitWorkflow">
            <v-text-field
              v-model="form.title"
              label="标题"
              :rules="rules.title"
              required
            />
            <v-textarea
              v-model="form.description"
              label="简介"
              :rules="rules.description"
              auto-grow
              rows="3"
              required
            />
            <v-select
              v-model="form.permission"
              label="权限"
              :items="[
                { value: 'public', title: '公开' },
                { value: 'private', title: '私有' },
              ]"
              :rules="rules.permission"
              required
            />
            <v-file-input
              v-model="form.coverFile"
              label="封面图片"
              accept="image/*"
              prepend-icon="mdi-image"
              show-size
              @change="handleFilePreview"
            />
            <v-img v-if="coverPreview" :src="coverPreview" max-width="100%" class="mt-2" />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn text @click="dialog = false">取消</v-btn>
          <v-btn color="primary" @click="submitWorkflow">提交</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="deleteDialog" max-width="400px">
      <v-card>
        <v-card-title class="text-h6">确认删除</v-card-title>
        <v-card-text>确定要删除该工作流吗？此操作不可恢复。</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn text @click="deleteDialog = false">取消</v-btn>
          <v-btn color="error" text @click="deleteConfirmed">确认</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="editInfoDialog" max-width="600px" persistent>
      <v-card>
        <v-card-title>修改工作流信息</v-card-title>
        <v-card-text>
          <v-form>
            <v-text-field
              v-model="editingWorkflow.name"
              label="标题"
              required
            />
            <v-textarea
              v-model="editingWorkflow.description"
              label="简介"
              auto-grow
              rows="3"
              required
            />
            <v-select
              v-model="editingWorkflow.status"
              label="权限"
              :items="[ { value: 'public', title: '公开' }, { value: 'private', title: '私有' } ]"
              required
            />
            <v-file-input
              v-model="editingWorkflow.coverFile"
              label="封面图片"
              accept="image/*"
              prepend-icon="mdi-image"
              show-size
              @change="handleEditCoverPreview"
            />
            <v-img v-if="editCoverPreview" :src="editCoverPreview" max-width="100%" class="mt-2" />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn text @click="editInfoDialog = false">取消</v-btn>
          <v-btn color="primary" @click="submitEditInfo">保存</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import server from '@/api/server'

const router = useRouter()

const dialog = ref(false)
const deleteDialog = ref(false)
const toDeleteId = ref<string | null>(null)

const formRef = ref()
const form = ref({
  title: '',
  description: '',
  permission: '',
  coverFile: null as File | null
})
const coverPreview = ref('')

const rules = {
  title: [
    (v: string) => !!v || '请输入标题',
    (v: string) => (v?.length >= 1 && v.length <= 30) || '长度在1到30个字符'
  ],
  description: [
    (v: string) => !!v || '请输入简介',
    (v: string) => (v?.length >= 2 && v.length <= 200) || '长度在2到200个字符'
  ],
  permission: [
    (v: string) => !!v || '请选择权限'
  ]
}

const beforeCoverUpload = (file: File) => {
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过5MB!')
    return false
  }
  return true
}

const handleCoverChange = (file: any) => {
  if (file && file.raw) {
    const fileObj = file.raw as File
    const reader = new FileReader()
    reader.onload = (e) => {
      coverPreview.value = e.target?.result as string
      form.value.coverFile = fileObj
    }
    reader.readAsDataURL(fileObj)
  }
}

const editCoverPreview = ref('')
const editCoverFile = ref<File | null>(null)

const handleEditCoverPreview = (file: any) => {
  if (file && file.raw) {
    const fileObj = file.raw as File
    const isImage = fileObj.type.startsWith('image/')
    const isLt5M = fileObj.size / 1024 / 1024 < 5

    if (!isImage) {
      ElMessage.error('只能上传图片文件')
      editingWorkflow.value.coverFile = null
      editCoverPreview.value = ''
      return
    }
    if (!isLt5M) {
      ElMessage.error('图片大小不能超过5MB')
      editingWorkflow.value.coverFile = null
      editCoverPreview.value = ''
      return
    }

    const reader = new FileReader()
    reader.onload = e => {
      editCoverPreview.value = e.target?.result as string
      editingWorkflow.value.coverFile = fileObj
    }
    reader.readAsDataURL(fileObj)
  }
}

const workflows = ref<any[]>([])
const filterStatus = ref('all')
const statusOptions = [
  { text: '全部', value: 'all' },
  { text: '已发布', value: 'public' },
  { text: '未发布', value: 'private' }
]

const filteredWorkflows = computed(() => {
  if (filterStatus.value === 'all') {
    return workflows.value
  } else if (filterStatus.value === 'public') {
    return workflows.value.filter(w => w.status === 'public')
  } else if (filterStatus.value === 'private') {
    return workflows.value.filter(w =>
      ['private', 'team', 'offline'].includes(w.status)
    )
  } else {
    return []
  }
})

async function fetchWorkflows() {
  try {
    const res = await server.Get('/workflow/list')
    workflows.value = res.data.workflows
  } catch {
    ElMessage.error('获取列表失败')
  }
}

onMounted(fetchWorkflows)

async function submitWorkflow() {
  const valid = await formRef.value.validate()
  if (!valid) return

  try {
    const { title, description, permission, coverFile } = form.value
    const createRes = await server.Post('/workflow/create', {
      name: title,
      description
    })
    const workflowId = createRes.data.workflow_id

    const formData = new FormData()
    formData.append('name', title)
    formData.append('description', description)
    formData.append('status', permission)
    if (coverFile) {
      formData.append('cover', coverFile)
    }

    await server.Post(`/workflow/${workflowId}/update`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    ElMessage.success('创建成功')
    dialog.value = false
    await fetchWorkflows()
    router.push({
      path: `/agents/${workflowId}/workflow`,
    })
  } catch (err: any) {
    ElMessage.error(err.message || '创建失败')
  }
}

const editInfoDialog = ref(false)
const editingWorkflow = ref<any>(null)

function editWorkflowInfo(w: any) {
  editingWorkflow.value = { ...w }
  editCoverPreview.value = w.cover
  editCoverFile.value = null
  editInfoDialog.value = true
}

async function submitEditInfo() {
  try {
    const { id, name, description, status, coverFile } = editingWorkflow.value
    const formData = new FormData()
    formData.append('name', name)
    formData.append('description', description)
    formData.append('status', status)
    if (coverFile) {
      formData.append('cover', coverFile)
    }

    const updateRes = await server.Post(`/workflow/${id}/update`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    if (updateRes.status === 200) {
      ElMessage.success('修改成功')
      editInfoDialog.value = false
      await fetchWorkflows()
    } else {
      ElMessage.error('修改失败')
    }
  } catch (err: any) {
    ElMessage.error(err.message || '修改失败')
  }
}

function confirmDelete(id: string) {
  toDeleteId.value = id
  deleteDialog.value = true
}

async function deleteConfirmed() {
  const res = await server.Delete(`/workflow/${toDeleteId.value}`)
  if (res.status === 200) {
    workflows.value = workflows.value.filter(w => w.id !== toDeleteId.value)
    ElMessage.success('已删除')
    deleteDialog.value = false
  } else {
    ElMessage.error('删除失败')
  }
}

function editWorkflow(id: string) {
  router.push({
    path: `/agents/${id}/workflow`,
  })
}

</script>

<style scoped>
.v-card {
  min-height: 280px;
  display: flex;
  flex-direction: column;
}
.v-card-title,
.v-card-subtitle {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.v-card-text {
  flex-grow: 1;
  overflow: hidden;
  text-overflow: ellipsis;
}
.v-card-actions {
  justify-content: flex-end;
}
.text-truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.v-card {
  min-height: auto !important;
}
</style>
