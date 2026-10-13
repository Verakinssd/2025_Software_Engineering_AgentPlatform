<template>
  <v-container>
    <el-carousel :interval="4000" type="card" height="30vh" class="my-4">
      <el-carousel-item v-for="(item, index) in popularAgents" :key="item">
        <v-card
          height="100%"
          rounded="lg"
          :to="`/agents/${item.id}/introduction`"
        >
          <v-img
            class="align-end text-white"
            :key="index"
            :src="item.cover || '/images/default_cover.jpg'"
            cover
            lazy-src="/images/default_cover.jpg"
          >
            <v-card-title style="backdrop-filter: blur(1px); background: linear-gradient(to top, rgba(0, 0, 0, 0.7), transparent); padding: 10px;">{{ item.name || '未命名' }}</v-card-title>
          </v-img>
        </v-card>
      </el-carousel-item>
    </el-carousel>
    <v-text-field
      style="margin-left: auto; margin-right: auto;"
      width="700px"
      v-model="search"
      class="mt-12"
      label="搜索"
      variant="solo"
      clearable
      color="primary"
      data-guide="search"
      @click:clear="fetchPublicAgents"
    >
      <template #append>
        <v-btn :loading="loading" @click="fetchPublicAgents" icon="mdi-magnify" color="primary"></v-btn>
      </template>
    </v-text-field>
    <v-row class="mt-4" data-guide="agents-grid">
      <v-col :key="index" v-for="(item, index) in agents" :xxl="2" :xl="3" :lg="3" :md="4" :sm="6" :xs="12" class="pa-6">
        <v-card :key="index"
          :to="`/agents/${item.id}/introduction`"
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
          <v-card-text style="height: 60px; overflow: hidden; text-overflow: ellipsis; display: -webkit-box; line-clamp: 3; -webkit-box-orient: vertical;">
            {{ item.description }}
          </v-card-text>
          <v-divider class="mt-2"></v-divider>
          <div class="ma-1">
            <v-chip prepend-icon="mdi-cursor-default-click-outline" variant="text" size="small">{{ item.clicks }}</v-chip>
            <v-chip :prepend-icon="item.likeStatus ? 'mdi-heart' : 'mdi-heart-outline'" variant="text" size="small">{{ item.likes }}</v-chip>
            <v-chip :prepend-icon="item.collectStatus ?'mdi-star' :'mdi-star-outline'" variant="text" size="small">{{ item.collects }}</v-chip>
          </div>
        </v-card>
      </v-col>
    </v-row>
    <div
      style="margin-left: auto; margin-right: auto; margin-top: 20px; text-align: center;">
      <v-progress-circular indeterminate v-if="loading"
        color="primary"
      ></v-progress-circular>
    </div>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { GetPopularAgents, GetPublicAgents } from '@/api/community';
import { ElMessage } from 'element-plus';

const loading = ref(false)
const end = ref(false)

const popularAgents = ref(null)
const agents = ref(null)

const page = ref(1)
const search = ref('')

const fetchPublicAgents = async (reload = true) => {
  loading.value = true
  if (reload) {
    page.value = 1
    end.value = false
  }
  const res = await GetPublicAgents(page.value, 12, search.value)
  if (res.status === 200) {
    if (reload) {
      agents.value = res.data.workflows
    } else {
      agents.value = agents.value.concat(res.data.workflows)
    }
    if (res.data.workflows.length < 12) {
      end.value = true
    }
  } else {
    ElMessage.error('获取公开智能体失败')
  }
  loading.value = false
}

const fetchPopularAgents = async () => {
  const res = await GetPopularAgents(5)
  if (res.status === 200) {
    popularAgents.value = res.data.workflows
  } else {
    ElMessage.error('获取热门智能体失败')
  }
}

onMounted(async () => {
  await fetchPopularAgents()
  await fetchPublicAgents()
})

function isScrolledToBottom() {
    const scrollTop = window.scrollY || document.documentElement.scrollTop;
    const scrollHeight = document.documentElement.scrollHeight;
    const clientHeight = document.documentElement.clientHeight;

    return scrollTop + clientHeight >= scrollHeight - 20;
}

window.addEventListener('scroll', async () => {
    if (isScrolledToBottom() && !loading.value &&!end.value) {
      loading.value = true
      page.value++
      await fetchPublicAgents(false)
    }
});

</script>
<style scoped>
.v-card {
  transition: transform 0.2s, box-shadow 0.2s;
}

.v-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}
</style>
