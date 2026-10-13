<template>
  <div class="profile-container">
    <v-card class="user-card">
      <div class="user-info">
        <v-avatar
          :size="100"
          :image="user.avatar !== '' ? user.avatar : '/images/default_avatar.jpg'"
          class="avatar"
        />
        <div class="user-details">
          <div class="name-section">
            <h2>{{ user.username }}</h2>
            <!-- <v-chip
              :color="user.gender === 'M' ? 'blue' : 'pink'"
              class="gender-tag"
              size="small"
            >
              <v-icon v-if="user.gender === 'M'">
                mdi-gender-male
              </v-icon>
              <v-icon v-else>
                mdi-gender-female
              </v-icon>
            </v-chip> -->
          </div>
          <p class="bio">
            {{ user.profile || "这个人很懒，什么都没留下~" }}
          </p>
          <div class="card-footer">
            <v-divider class="stats-divider" />
            <div class="stats-container">
              <div class="stats-item">
                <span class="stat-number">{{ stats.posts }}</span>
                <span class="stat-label">发布</span>
              </div>
              <div class="stats-item">
                <span class="stat-number">{{ stats.collects }}</span>
                <span class="stat-label">收藏</span>
              </div>
              <div class="stats-item">
                <span class="stat-number">{{ stats.views }}</span>
                <span class="stat-label">点击量</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </v-card>

    <v-tabs v-model="activeTab" class="profile-tabs">
      <v-tab value="posts">创建</v-tab>
      <v-tab value="favorites">收藏</v-tab>
      <!-- <v-tab value="home">工作区</v-tab> -->
      <v-tab value="settings">设置</v-tab>
    </v-tabs>

    <v-window v-model="activeTab" class="mt-4">
      <v-window-item value="home">
        <Home v-if="activeTab === 'home'" />
      </v-window-item>
      <v-window-item value="posts">
        <Posts v-if="activeTab === 'posts'" />
      </v-window-item>
      <v-window-item value="favorites">
        <Favorites v-if="activeTab === 'favorites'" />
      </v-window-item>
      <v-window-item value="settings">
        <Settings v-if="activeTab === 'settings'" />
      </v-window-item>
    </v-window>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useDisplay } from 'vuetify'
import Home from '@/components/profile/Home.vue'
import Posts from '@/components/profile/Posts.vue'
import Favorites from '@/components/profile/Favorites.vue'
import Settings from '@/components/profile/Settings.vue'
import { user } from "@/stores/user"
import server from "@/api/server"
import {ElMessage} from "element-plus";

const { mobile } = useDisplay()
const stats = ref({
  posts: 0,
  collects: 0,
  views: 0,
  likes: 0
})

const activeTab = ref('posts')

const posts = ref([])
const favorites = ref([])

const fetchUserData = async () => {
  try {
    const profileRes = await server.Get('/user/profile')

    const postsRes = await server.Get('/workflow/list')
    console.log(postsRes.data)

    let totalViews = 0
    let totalLikes = 0
    if (postsRes.data.workflows && postsRes.data.workflows.length > 0) {
      postsRes.data.workflows.forEach(post => {
        totalViews += post.clicks || 0
        totalLikes += post.likes || 0
      })
    }

    stats.value = {
      posts: profileRes.data.published_workflow_cnt || 0,
      collects: profileRes.data.favourite_workflow_cnt || 0,
      views: totalViews,
      likes: totalLikes
    }
  } catch (error) {
    ElMessage.error('获取用户数据失败')
    console.error(error)
  }
}

onMounted(() => {
  fetchUserData()
})
</script>

<style scoped>
.profile-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.user-card {
  position: relative;
  min-height: 200px;
  padding-bottom: 60px;
  margin-bottom: 30px;
  border-radius: 12px;
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.08);
}

.user-info {
  display: flex;
  align-items: center;
  padding: 20px 20px 20px 40px;
}

.avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  cursor: pointer;
  margin-right: 20px;
  border: 3px solid #f5f5f5;
}

.user-details {
  flex: 1;
  padding-top: 10px;
  padding-bottom: 10px;
}

.name-section {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.name-section h2 {
  margin: 0;
  font-size: 24px;
  margin-right: 10px;
  font-weight: 600;
}

.gender-tag {
  font-size: 16px;
}

.bio {
  margin: 0;
  color: #666;
  font-size: 14px;
  cursor: pointer;
  line-height: 1.6;
  min-height: 40px;
}

.bio:hover {
  color: #333;
}

.profile-tabs {
  margin-bottom: 20px;
}

.stat-label {
  color: #666;
  font-size: 14px;
}

.card-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 0 20px 15px 20px;
}

.stats-divider {
  margin: 12px 0;
  border-top: 2px solid rgba(0, 0, 0, 0.05);
}

.stats-container {
  display: flex;
  justify-content: space-around;
  align-items: center;
  height: 40px;
}

.stats-item {
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100%;
  text-align: center;
}

.stat-number {
  font-size: 16px;
  font-weight: bold;
  line-height: 1;
}

.stat-label {
  font-size: 12px;
  color: #999;
  line-height: 1.2;
  margin-top: 2px;
}

:deep(.v-tab) {
  font-size: 20px;
  font-weight: 500;
  text-transform: none;
  letter-spacing: normal;
}

:deep(.v-slide-group__content) {
  justify-content: flex-start;
}
</style>
