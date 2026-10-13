<template>
  <div class="favorites-content">
    <el-empty v-if="favorites.length === 0" description="暂无收藏内容" />
    <v-row class="mt-4">
      <v-col
        v-for="item in favorites"
        :key="item.id"
        :xxl="2" :xl="3" :lg="3" :md="4" :sm="6" :xs="12"
        class="pa-6"
      >
        <v-card
          class="favorite-card"
          height="100%"
          @click="navigateToPost(item.id)"
        >
          <v-img
            :aspect-ratio="16/9"
            :src="item.image"
            cover
            lazy-src="/images/default_cover.jpg"
          ></v-img>

          <v-card-title>
            {{ item.title || '未命名' }}
          </v-card-title>

          <v-card-text
            style="height: 60px; overflow: hidden; text-overflow: ellipsis;
                   display: -webkit-box; line-clamp: 3; -webkit-box-orient: vertical;"
          >
            {{ item.description || '暂无描述' }}
          </v-card-text>

          <v-divider class="mt-2"></v-divider>

          <div class="d-flex justify-space-between align-center ma-2">
            <div>
              <v-chip prepend-icon="mdi-cursor-default-click-outline" variant="text" size="small">{{ item.clicks }}</v-chip>
<!--              <v-chip :prepend-icon="item.likeStatus ? 'mdi-heart' : 'mdi-heart-outline'" variant="text" size="small">{{ item.likes }}</v-chip>-->
<!--              <v-chip :prepend-icon="item.collectStatus ?'mdi-star' :'mdi-star-outline'" variant="text" size="small">{{ item.collects }}</v-chip>-->
            </div>

            <v-chip
              size="small"
              :color="getStatusColor(item.status)"
              class="ml-2"
            >
              {{ statusLabel(item.status) }}
            </v-chip>
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
const favorites = ref([])

const fetchFavorites = async () => {
  try {
    const res = await server.Get('/workflow/favorites')
    const raw = res.data?.workflows || []
    console.log(raw)

    favorites.value = raw
      .filter(item => item.correspondingWorkflow)
      .map(item => {
        const wf = item.correspondingWorkflow
        return {
          id: wf.id,
          title: wf.name || '未命名',
          image: wf.cover || '/public/images/default_cover.jpg',
          status: wf.status || 'private',
          name: wf.name,
          description: wf.description,
          cover: wf.cover || '',
          likes: wf.likes || 0,
          clicks: wf.clicks || 0,
          likeStatus: wf.likeStatus || false,
          collectStatus: wf.collectStatus || false,
          collects: wf.collects || 0
        }
      })
  } catch (error) {
    ElMessage.error('获取收藏失败')
    console.error('获取收藏失败', error)
  }
}

const navigateToPost = (id) => {
  router.push(`/agents/${id}/introduction`)
}

const statusLabel = status => {
  return {
    private: '私有',
    public: '公开'
  }[status] || status
}

const getStatusColor = (status) => {
  const colors = {
    private: 'grey',
    public: 'green'
  }
  return colors[status] || 'grey'
}

onMounted(fetchFavorites)
</script>

<style scoped>
.favorites-content {
  padding: 20px;
}

.favorite-card {
  transition: transform 0.2s, box-shadow 0.2s;
  cursor: pointer;
}

.favorite-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

.v-chip {
  font-size: 12px;
  height: 24px;
}
</style>
