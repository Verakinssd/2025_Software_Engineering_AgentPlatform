<template>
  <transition name="fade">
    <div v-show="AppStore.loading" class="loading-container">
      <svg style="left: 50%;
        top: 50%;
        position: absolute;
        transform: translate(-50%, -50%) matrix(1, 0, 0, 1, 0, 0);" preserveAspectRatio="xMidYMid meet" viewBox="0 0 187.3 93.7" height="300px" width="400px">
        <path d="M93.9,46.4c9.3,9.5,13.8,17.9,23.5,17.9s17.5-7.8,17.5-17.5s-7.8-17.6-17.5-17.5c-9.7,0.1-13.3,7.2-22.1,17.1 				c-8.9,8.8-15.7,17.9-25.4,17.9s-17.5-7.8-17.5-17.5s7.8-17.5,17.5-17.5S86.2,38.6,93.9,46.4z" stroke-miterlimit="10" stroke-linejoin="round" stroke-linecap="round" stroke-width="4" fill="none" id="outline" stroke="#4E4FEB"></path>
        <path d="M93.9,46.4c9.3,9.5,13.8,17.9,23.5,17.9s17.5-7.8,17.5-17.5s-7.8-17.6-17.5-17.5c-9.7,0.1-13.3,7.2-22.1,17.1 				c-8.9,8.8-15.7,17.9-25.4,17.9s-17.5-7.8-17.5-17.5s7.8-17.5,17.5-17.5S86.2,38.6,93.9,46.4z" stroke-miterlimit="10" stroke-linejoin="round" stroke-linecap="round" stroke-width="4" stroke="#4E4FEB" fill="none" opacity="0.05" id="outline-bg"></path>
      </svg>
    </div>
  </transition>

  <v-responsive v-if="!AppStore.loading">
    <v-app>
      <v-navigation-drawer
        v-if="showSidebar && !isAdminRoute"
        :mobile="false"
        :rail-width="64"
        :rail="true"
        expand-on-hover
        permanent
      >
        <v-list nav>
          <v-list-item class="my-2" data-guide="logo">
            <template #prepend>
              <v-avatar size="32">
                <v-icon color="blue" class="ml-1">mdi-infinity</v-icon>
              </v-avatar>
            </template>
            <v-list-item-title class="ml-3 py-2 font-weight-bold text-h6">LingXi</v-list-item-title>
          </v-list-item>
          <v-divider class="my-2"></v-divider>
          <v-list-item to="/" data-guide="home">
            <template #prepend>
              <v-icon class="ml-1">mdi-home</v-icon>
            </template>
            <v-list-item-title>主页</v-list-item-title>
          </v-list-item>
          <!-- v-list-item to="/search">
            <template #prepend>
              <v-icon class="ml-1">mdi-magnify</v-icon>
            </template>
            <v-list-item-title>搜索</v-list-item-title>
          </v-list-item -->
          <!-- <v-list-item to="/workfloweditor">
            <template #prepend>
              <v-icon class="ml-1">mdi-file-tree</v-icon>
            </template>
            <v-list-item-title>工作流</v-list-item-title>
          </v-list-item> -->
          <v-list-item to="/workspace" data-guide="workspace">
            <template #prepend>
              <v-icon class="ml-1">mdi-briefcase-variant</v-icon>
            </template>
            <v-list-item-title>工作空间</v-list-item-title>
          </v-list-item>
<!--          <v-list-item to="/chat">-->
<!--            <template #prepend>-->
<!--              <v-icon class="ml-1">mdi-forum-outline</v-icon>-->
<!--            </template>-->
<!--            <v-list-item-title>智能体</v-list-item-title>-->
<!--          </v-list-item>-->
          <v-list-item to="/knowledge" data-guide="knowledge">
            <template #prepend>
              <v-icon class="ml-1">mdi-book-open-variant</v-icon>
            </template>
            <v-list-item-title>知识库</v-list-item-title>
          </v-list-item>
          <v-list-item to="/profile" data-guide="profile">
            <template #prepend>
              <v-icon class="ml-1">mdi-account-outline</v-icon>
            </template>
            <v-list-item-title>个人中心</v-list-item-title>
          </v-list-item>
        </v-list>

        <template #append>
          <v-list nav>
            <!-- 用户手册菜单项已注释掉
            <v-list-item to="/manual">
              <template #prepend>
                <v-icon class="ml-1">mdi-book-outline</v-icon>
              </template>
              <v-list-item-title>用户手册</v-list-item-title>
            </v-list-item>
            -->
            <v-list-item to="/message">
              <template #prepend>
                <v-icon class="ml-1">mdi-bell-outline</v-icon>
              </template>
              <v-list-item-title>消息通知</v-list-item-title>
            </v-list-item>
            <v-list-item to="/admin/user" v-if="user.role === 'admin'">
              <template #prepend>
                <v-icon class="ml-1">mdi-account-supervisor-outline</v-icon>
              </template>
              <v-list-item-title>管理员视图</v-list-item-title>
            </v-list-item>
            <v-list-item :to="user.loggedIn ? '' : '/login'" @click="handleLogout">
              <template #prepend>
                <v-avatar size="32">
                  <v-img :src="user.avatar !== '' ? user.avatar : '/images/default_avatar.jpg'" ></v-img>
                </v-avatar>
              </template>
              <template #append>
                <v-icon class="ml-1" v-if="user.loggedIn">mdi-logout</v-icon>
              </template>
              <v-list-item-title class="ml-3">{{ user.loggedIn ? user.username : '去登录' }}</v-list-item-title>
            </v-list-item>
          </v-list>
            <!-- <div class="mb-2" style="margin-top: -18px;">
              <el-popover
                placement="right-end"
                :width="220"
                trigger="hover"
                v-model:visible="popoverVisible"
              >
                <template #reference>
                  <v-list nav>
                    <v-list-item>
                      <template #prepend>
                        <v-avatar size="32">
                          <v-img :src="user.avatar !== '' ? user.avatar : '/images/default_avatar.jpg'" ></v-img>
                        </v-avatar>
                      </template>
                      <v-list-item-title class="ml-3">{{ user.username }}</v-list-item-title>
                    </v-list-item>
                  </v-list>
                </template>
                <div v-if="user.loggedIn" class="popover-content">
                  <div class="popover-user-info" @click="() => { closePopover(); router.push('/profile') }">
                    <el-avatar
                      :src="user.avatar !== '' ? user.avatar : '/images/default_avatar.jpg'"
                      class="popover-avatar"
                      :size="56"
                    />
                    <div class="popover-username">{{ user.username }}</div>
                    <div class="popover-email">{{ user.email }}</div>
                  </div>

                  <v-btn v-if="user.role === 'admin'" variant="text" to="/admin/user" color="blue" block>管理员视图</v-btn>
                  <v-divider class="my-2"></v-divider>
                  <v-btn variant="text" @click="handleLogout" color="red" block>退出登录</v-btn>
                </div>

                <div v-else class="popover-content">
                  <v-btn variant="text" to="/login" color="blue" block append-icon="mdi-chevron-right" @click="closePopover">登录</v-btn>
                </div>
              </el-popover>
            </div> -->
        </template>
      </v-navigation-drawer>

      <AdminSidebar v-if="isAdminRoute" />

      <v-main>
        <router-view></router-view>
      </v-main>

      <!-- 新手指引浮动按钮 -->
      <v-btn
        v-if="$route.path === '/' && user.loggedIn"
        icon
        size="large"
        variant="outlined"
        color="primary"
        @click="startUserGuide"
        style="position: fixed; bottom: 80px; right: 20px; z-index: 1000; background: white; border-width: 2px;"
        elevation="4"
        class="guide-btn"
      >
        <v-icon size="28">mdi-help</v-icon>
        <v-tooltip activator="parent" location="left">
          新手指引
        </v-tooltip>
      </v-btn>

      <!-- 新手指引组件 -->
      <UserGuide ref="userGuideRef" />
    </v-app>
  </v-responsive>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { user } from './stores/user'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus';
import { getProfile, logout } from './api/user';
import { onMounted } from 'vue'
import { AppStore } from './stores/app';
import AdminSidebar from './components/admin/AdminSidebar.vue';
import UserGuide from './components/UserGuide.vue';

const router = useRouter()
const route = useRoute()
const userGuideRef = ref<InstanceType<typeof UserGuide> | null>(null)

const isAdminRoute = computed(() => route.path.startsWith('/admin'));
const showSidebar = computed(() =>
  route.path !== '/login' &&
  route.path !== '/register' &&
  !(route.path.startsWith('/agents/') && route.path.endsWith('/workflow')) &&
    !route.path.startsWith('/admin'))
const popoverVisible = ref(false)

onMounted(async () => {
  await getProfile()
  AppStore.loading = false
})

async function handleLogout() {
  if (!user.loggedIn) {
    return;
  }
  const res = await logout()
  if (res.status === 401) {
    //ElMessage.success('退出登录成功')
    await router.push('/login')
    return;
  }
  if (res.status !== 200) {
    ElMessage.error('退出登录失败')
    return
  }
  popoverVisible.value = false
  ElMessage.success('退出登录成功')
  await router.push('/login')
  //console.log("logout")
}

const closePopover = () => {
  popoverVisible.value = false;
};

// 启动新手指引
const startUserGuide = () => {
  // 确保在主页才能开始指引
  if (route.path !== '/') {
    router.push('/').then(() => {
      setTimeout(() => {
        userGuideRef.value?.startGuide()
      }, 500)
    })
  } else {
    userGuideRef.value?.startGuide()
  }
}

</script>

<style scoped>

.popover-content {
  padding: 8px;
}

.popover-user-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 12px;
  cursor: pointer;
  transition: background-color 0.2s;
  border-radius: 8px;
}

.popover-user-info:hover {
  background-color: #f5f5f5;
}

.popover-avatar {
  margin-bottom: 8px;
}

.popover-username {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.popover-email {
  font-size: 13px;
  color: #888;
}

.user-info, .login-tip {
  display: flex;
  align-items: center;
  padding: 8px;
  cursor: pointer;
}

.user-info:hover, .login-tip:hover {
  background-color: #f5f5f5;
}

.nickname {
  margin-left: 8px;
  flex-grow: 1;
}

/* loading 居中 */
.loading-container {
  display: flex;
  justify-content: center;  /* 水平居中 */
  align-items: center;      /* 垂直居中 */
  height: 100vh;            /* 占满整个视口高度 */
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

#outline {
  stroke-dasharray: 2.42777px, 242.77666px;
  stroke-dashoffset: 0;
  -webkit-animation: anim 1.6s linear infinite;
  animation: anim 1.6s linear infinite;
}

@-webkit-keyframes anim {
  12.5% {
    stroke-dasharray: 33.98873px, 242.77666px;
    stroke-dashoffset: -26.70543px;
  }

  43.75% {
    stroke-dasharray: 84.97183px, 242.77666px;
    stroke-dashoffset: -84.97183px;
  }

  100% {
    stroke-dasharray: 2.42777px, 242.77666px;
    stroke-dashoffset: -240.34889px;
  }
}

@keyframes anim {
  12.5% {
    stroke-dasharray: 33.98873px, 242.77666px;
    stroke-dashoffset: -26.70543px;
  }

  43.75% {
    stroke-dasharray: 84.97183px, 242.77666px;
    stroke-dashoffset: -84.97183px;
  }

  100% {
    stroke-dasharray: 2.42777px, 242.77666px;
    stroke-dashoffset: -240.34889px;
  }
}

/* 新手指引按钮样式 */
.guide-btn {
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.guide-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 8px 20px rgba(25, 118, 210, 0.3);
}

</style>
