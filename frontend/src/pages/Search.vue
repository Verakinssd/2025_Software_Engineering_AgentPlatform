<template>
  <div class="search-page">
    <div class="carousel-container">
      <el-carousel :interval="4000" type="card" height="300px" indicator-position="outside">
        <el-carousel-item v-for="agent in featuredAgents" :key="agent.id">
          <div class="carousel-item" @click="goToPost(agent.id)">
            <img :src="agent.image" class="carousel-image" />
            <div class="carousel-title">{{ agent.title }}</div>
          </div>
        </el-carousel-item>
      </el-carousel>
    </div>

    <!-- 顶部搜索框 -->
    <div class="search-header">
      <el-input
        v-model="searchQuery"
        placeholder="搜索智能体..."
        prefix-icon="el-icon-search"
        clearable
        class="search-input"
      />
      <el-button type="primary" @click="drawerVisible = true" class="create-btn">
        测试上传
      </el-button>
      <el-button type="primary" @click="getWorkflow" class="create-btn">
        测试获取
      </el-button>
    </div>

    <el-drawer
      v-model="drawerVisible"
      title="创建智能体"
      direction="rtl"
      size="40%"
    >
      <el-form :model="workflowForm" :rules="rules" ref="workflowFormRef" label-position="top">
        <el-form-item label="标题" prop="title">
          <el-input v-model="workflowForm.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="简介" prop="description">
          <el-input
            v-model="workflowForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入简介"
            style="white-space: pre-wrap;"
          />
        </el-form-item>
        <el-form-item label="封面">
          <el-upload
            class="cover-uploader"
            action="#"
            :show-file-list="false"
            :auto-upload="false"
            :on-change="handleCoverChange"
            accept="image/*"
            :before-upload="beforeCoverUpload"
          >
            <img v-if="coverPreview" :src="coverPreview" class="cover-image" />
            <el-icon v-else class="cover-uploader-icon"><Plus /></el-icon>
          </el-upload>
        </el-form-item>
        <el-form-item label="权限" prop="permission">
          <el-select
            v-model="workflowForm.permission"
            placeholder="请选择权限"
          >
            <el-option
              v-for="option in permissionOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitWorkflow">提交</el-button>
          <el-button @click="drawerVisible = false">取消</el-button>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-empty v-if="filteredAgents.length === 0" description="未找到相关智能体" />

    <div v-else class="agent-grid">
      <div v-for="agent in filteredAgents" :key="agent.id" class="agent-item">
        <el-card shadow="hover" class="agent-card" @click="goToPost(agent.id)">
          <img :src="agent.image" class="agent-image" />
          <div class="agent-info">
            <span class="agent-title">{{ agent.title }}</span>
            <div class="agent-bottom">
              <el-tag size="small">{{ agent.type }}</el-tag>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus } from "@element-plus/icons-vue"
import { ElMessage, FormInstance } from "element-plus"
import server from "@/api/server"

const router = useRouter()

const agents = ref([])

// 提取前5个作为轮播图 featured
const featuredAgents = computed(() => agents.value.slice(0, 5))

const searchQuery = ref('')

const filteredAgents = computed(() => {
  if (!searchQuery.value) return agents.value
  return agents.value.filter(agent =>
    agent.title.toLowerCase().includes(searchQuery.value.toLowerCase()))
})

const permissionOptions = [
  { label: '私有', value: 'private' },
  { label: '公开', value: 'public' }
]

const goToPost = (agent_id) => {
  router.push(`/agents/${agent_id}/introduction`)
}

const fetchAgents = async (page = 1, pageSize = 50) => {
  try {
    const res = await server.Get('/workflow/public', {
      params: {
        page: page.toString(),
        pageSize: pageSize.toString()
      }
    });

    agents.value = res.data.workflows.map(item => ({
      id: item.id,
      title: item.name,
      description: item.description,
      image: item.cover || '/public/images/default_cover.jpg',
      type: item.status
    }));
  } catch (error) {
    ElMessage.error('获取智能体列表失败');
    console.error(error);
  }
};

onMounted(() => {
  fetchAgents(1, 50)
})

const drawerVisible = ref(false)
const workflowFormRef = ref<FormInstance>()

const workflowForm = reactive({
  title: '',
  description: '',
  cover: null as File | null,
  permission: 'private'
})

const rules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { min: 1, max: 30, message: '长度在1到30个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入简介', trigger: 'blur' },
    { min: 2, max: 200, message: '长度在2到200个字符', trigger: 'blur' }
  ],
  permission: [
    { required: true, message: '请选择权限', trigger: 'change' }
  ]
}

const coverPreview = ref('')

const beforeCoverUpload = (file) => {
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

const handleCoverChange = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    coverPreview.value = e.target?.result as string
    workflowForm.cover = file.raw
  }
  reader.readAsDataURL(file.raw)
}

const getWorkflow = async () => {
  const res = await server.Get('/workflow/list')
  console.log(res.data)
}

const submitWorkflow = async () => {
  try {
    await workflowFormRef.value.validate()

    const createResponse = await server.Post('/workflow/create', {
      name: workflowForm.title,
      description: workflowForm.description
    })

    const workflowId = createResponse.data.workflow_id
    ElMessage.success('工作流创建成功')

    const formData = new FormData()
    formData.append('name', workflowForm.title)
    formData.append('description', workflowForm.description)
    formData.append('status', workflowForm.permission)
    if (workflowForm.cover) {
      formData.append('cover', workflowForm.cover)
    }

    await server.Post(`/workflow/${workflowId}/update`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    ElMessage.success('更新成功')
    drawerVisible.value = false

    // 重新拉取更新后的数据
    fetchAgents(1, 50)
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
    console.error('Error:', error)
  }
}
</script>

<style scoped>
.search-page {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.carousel-container {
  margin: 0 auto 30px;
  width: 90%;
}

.carousel-item {
  position: relative;
  height: 100%;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
}

.carousel-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.carousel-title {
  position: absolute;
  left: 0;
  bottom: 0;
  width: 100%;
  padding: 15px;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
  color: white;
  font-size: 18px;
  font-weight: bold;
}

.search-header {
  display: flex;
  justify-content: center;
  margin-bottom: 30px;
}

.search-input {
  width: 400px;
}

.agent-card {
  transition: transform 0.3s;
  height: 100%;
  border-radius: 8px;
  cursor: pointer;
}

.agent-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.agent-image {
  width: 100%;
  height: 180px;
  object-fit: cover;
  border-radius: 8px 8px 0 0;
}

.agent-info {
  padding: 15px;
}

.agent-title {
  font-weight: bold;
  font-size: 16px;
  display: block;
  margin-bottom: 8px;
}

.agent-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

@media (max-width: 768px) {
  .carousel-container {
    width: 100%;
  }
}

@media (max-width: 576px) {
  .search-input {
    width: 100%;
  }
}

.agent-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
  padding: 10px;
}

.agent-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.agent-image {
  height: 180px;
  object-fit: cover;
}

.cover-uploader {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  width: 320px;
  height: 180px;
}

.cover-uploader:hover {
  border-color: var(--el-color-primary);
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cover-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 320px;
  height: 180px;
  text-align: center;
  line-height: 148px;
}
</style>
