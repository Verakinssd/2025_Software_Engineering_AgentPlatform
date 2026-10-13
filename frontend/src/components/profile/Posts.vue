<template>
  <div class="posts-content">
    <el-empty v-if="posts.length === 0" description="暂无发布内容" />
    <v-row class="mt-4">
      <v-col
        v-for="post in posts"
        :key="post.id"
        :xxl="2" :xl="3" :lg="3" :md="4" :sm="6" :xs="12"
        class="pa-6"
      >
        <v-card
          class="post-card"
          height="100%"
          @click="navigateToPost(post.id)"
        >
          <v-img
            :aspect-ratio="16/9"
            :src="post.image"
            cover
            lazy-src="/images/default_cover.jpg"
          ></v-img>

          <v-card-title>
            {{ post.title || '未命名' }}
          </v-card-title>

          <v-card-text
            style="height: 60px; overflow: hidden; text-overflow: ellipsis;
                   display: -webkit-box; line-clamp: 3; -webkit-box-orient: vertical;"
          >
            {{ post.description }}
          </v-card-text>

          <v-divider class="mt-2"></v-divider>

          <div class="d-flex justify-space-between align-center ma-2">
            <div>
              <v-chip prepend-icon="mdi-cursor-default-click-outline" variant="text" size="small">{{ post.clicks }}</v-chip>
<!--              <v-chip :prepend-icon="post.likeStatus ? 'mdi-heart' : 'mdi-heart-outline'" variant="text" size="small">{{ post.likes }}</v-chip>-->
            </div>

            <v-select
              v-model="post.status"
              :items="statusOptions"
              density="compact"
              variant="outlined"
              hide-details
              @update:modelValue="handleStatusChange(post)"
              @click.stop
              style="max-width: 100px;"
            ></v-select>
          </div>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import server from '@/api/server'

const router = useRouter()
const posts = ref([])

const statusOptions = [
  { title: '私有', value: 'private' },
  { title: '公开', value: 'public' },
]

const fetchMyAgents = async () => {
  try {
    const res = await server.Get('/workflow/list')
    console.log(res.data)
    posts.value = res.data.workflows.map(item => ({
      id: item.id,
      title: item.name,
      image: item.cover || '/public/images/default_cover.jpg',
      status: item.status || 'private',
      name: item.name,
      description: item.description,
      cover: item.cover || '',
      clicks: item.clicks || 0,
      likes: item.likes || 0,
      likeStatus: item.likeStatus || 0
    }))
  } catch (error) {
    ElMessage.error('获取列表失败')
  }
}

const statusLabel = status => {
  return { private: '私有', public: '公开' }[status] || status
}

const handleStatusChange = async (post) => {
  try {
    const { id, name, description, cover, status } = post
    const formData = new FormData()
    formData.append('name', name)
    formData.append('description', description)
    formData.append('status', status)
    // formData.append('cover', cover)

    await server.Post(`/workflow/${id}/update`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    ElMessage.success('权限更新成功')
  } catch {
    ElMessage.error('权限更新失败')
  }
}

const navigateToPost = id => {
  router.push(`/agents/${id}/introduction`)
}

onMounted(fetchMyAgents)
</script>

<style scoped>
.posts-content {
  padding: 20px;
}

.post-card {
  transition: transform 0.2s, box-shadow 0.2s;
  cursor: pointer;
}

.post-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

:deep(.v-select .v-field) {
  padding: 0 8px !important;
  min-height: 32px !important;
}

:deep(.v-select .v-field__input) {
  padding: 0 !important;
  min-height: 32px !important;
}

:deep(.v-select .v-field__append-inner) {
  padding-top: 0 !important;
  align-items: center !important;
}
</style>
