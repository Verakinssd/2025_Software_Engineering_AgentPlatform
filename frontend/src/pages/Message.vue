<template>
  <v-card class="notification-container mt-6" flat>
    <header class="d-flex align-center px-4 py-3 bg-white" style="height: 64px">
      <h1 class="text-h4 font-weight-bold" style="font-size: 1.75rem; line-height: 1.2">
        通知中心
      </h1>
      <v-spacer></v-spacer>
      <v-btn
        color="primary"
        variant="text"
        prepend-icon="mdi-check-all"
        @click="markAllAsRead"
        class="mr-2"
      >
        全部已读
      </v-btn>
      <v-btn
        icon="mdi-refresh"
        variant="text"
        size="large"
        @click="fetchNotifications"
      ></v-btn>
    </header>
    <v-divider color="grey-darken-1"></v-divider>

    <div v-if="notifications.length === 0" class="empty-notification py-16 text-center">
      <v-icon size="64" color="grey-lighten-1">mdi-bell-off-outline</v-icon>
      <div class="text-h6 mt-4 text-grey">暂无消息</div>
    </div>

    <v-list v-else class="notification-list" lines="two">
      <v-list-item
        v-for="notification in notifications"
        :key="notification.id"
        @click="showNotification(notification)"
        class="notification-item py-3"
      >
        <template v-slot:prepend>
          <v-avatar
            :color="getNotificationColor(notification)"
            size="48"
          >
            <v-icon
              color="white"
              size="x-large"
            >
              {{ getNotificationIcon(notification) }}
            </v-icon>
          </v-avatar>
        </template>

        <v-list-item-title class="d-flex align-center text-body-1">
          <span class="font-weight-medium">{{ notification.from_user_name }}</span>
          <span class="text-body-2 text-grey ml-3">{{ formatDate(notification.date) }}</span>
        </v-list-item-title>
        <v-list-item-subtitle class="text-wrap text-body-1">{{ notification.title }}</v-list-item-subtitle>

        <template v-slot:append>
          <v-icon color="grey" size="large">mdi-chevron-right</v-icon>
        </template>
      </v-list-item>
    </v-list>

    <v-footer class="pa-4 bg-white" app>
      <div class="d-flex justify-center w-100 gap-4">
        <v-btn
          variant="text"
          :disabled="currentPage <= 1"
          @click="prevPage"
          prepend-icon="mdi-chevron-left"
        >
          上一页
        </v-btn>
        <v-btn
          variant="text"
          :disabled="notifications.length < pageSize"
          @click="nextPage"
          append-icon="mdi-chevron-right"
        >
          下一页
        </v-btn>
      </div>
    </v-footer>

    <v-dialog v-model="dialog" max-width="800">
      <v-card>
        <v-card-title class="text-h6">{{ currentNotification?.title }}</v-card-title>
        <v-divider></v-divider>
        <v-card-text class="markdown-content" v-html="renderedContent"></v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" variant="text" @click="dialog = false">关闭</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script setup>
import { ref, computed } from 'vue'
import { marked } from 'marked'
import server from "@/api/server";
import {ElMessage} from "element-plus";

const notifications = ref([])
const currentPage = ref(1)
const pageSize = ref(8)
const dialog = ref(false)
const currentNotification = ref(null)

const VIOLATION_TYPES = [
    '警告',
    '违规评论',
    '违规智能体',
    '临时封禁',
    '永久封禁'
]

const isViolationNotification = (notification) => {
    return VIOLATION_TYPES.includes(notification.title)
}

const isCommentNotification = (notification) => {
  return notification.title.includes('回复了') || notification.title.includes('评论了');
}

const getNotificationColor = (notification) => {
  if (isViolationNotification(notification)) {
    return notification.is_read ? 'grey-lighten-2' : 'orange-lighten-1';
  }
  if (isCommentNotification(notification)) {
    return notification.is_read ? 'grey-lighten-2' : 'green-darken-1';
  }
  return notification.is_read ? 'grey-lighten-2' : 'blue-darken-2';
}

const getNotificationIcon = (notification) => {
  if (isViolationNotification(notification)) {
    return notification.is_read ? 'mdi-alert' : 'mdi-alert-outline';
  }
  if (isCommentNotification(notification)) {
    return notification.is_read ? 'mdi-comment-text' : 'mdi-comment-text-outline';
  }
  return notification.is_read ? 'mdi-email' : 'mdi-email-outline';
}

const renderedContent = computed(() => {
  return currentNotification.value
    ? marked.parse(currentNotification.value.content || '')
    : ''
})

// const sortedNotifications = computed(() => {
//     return [...notifications.value].sort((a, b) => {
//         if (a.is_read === b.is_read) return 0
//         return a.is_read ? 1 : -1
//     })
// })

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    fetchNotifications()
  }
}

const nextPage = () => {
  currentPage.value++
  fetchNotifications()
}

const markAsRead = async (notificationId) => {
    try {
        await server.Post(`/notification/read/${notificationId}`)
        const notification = notifications.value.find(n => n.id === notificationId)
        if (notification) {
            notification.is_read = true
        }
    } catch (error) {
        console.error('标记已读失败:', error)
    }
}

const markAllAsRead = async () => {
    try {
        await server.Post('/notification/read_all')
        notifications.value.forEach(n => n.is_read = true)
        ElMessage.success('所有通知已标记为已读')
    } catch (error) {
        console.error('标记全部已读失败:', error)
        ElMessage.error('标记全部已读失败')
    }
}

const fetchNotifications = async () => {
  try {
    const response = await server.Get('/notification/list', {
      params: {
        page: currentPage.value,
        per_page: pageSize.value
      }
    })
    notifications.value = response.data.notifications || []
    console.log(response.data.notifications)

    if (notifications.value.length === 0 && currentPage.value > 1) {
      currentPage.value--
      fetchNotifications()
      ElMessage.info("已到达最后一页")
    }
  } catch (error) {
    console.error('获取通知失败:', error)
  }
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString()
}

const showNotification = async (notification) => {
  currentNotification.value = notification
  dialog.value = true

  if (!notification.is_read) {
      await markAsRead(notification.id)
  }
}

fetchNotifications()
</script>

<style scoped>
header {
  position: relative;
  z-index: 2;
  box-shadow: 0 2px 4px -1px rgba(0,0,0,0.1);
}

.notification-container {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  box-shadow: none;
  border-radius: 0;
}

.notification-list {
  padding: 12px 0;
}

.notification-item {
  padding: 16px 32px;
  transition: background-color 0.2s;
}

.notification-item:hover {
  background-color: rgba(0, 0, 0, 0.04);
}

.markdown-content {
  padding: 20px;
  font-size: 1.05rem;
  line-height: 1.8;
}

.markdown-content :deep(h1) {
  font-size: 1.5em;
  margin: 1em 0 0.5em;
}

.markdown-content :deep(h2) {
  font-size: 1.3em;
  margin: 1em 0 0.5em;
}

.markdown-content :deep(p) {
  margin-bottom: 1em;
}

.markdown-content :deep(a) {
  color: #1976d2;
  text-decoration: none;
}

.markdown-content :deep(a:hover) {
  text-decoration: underline;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  padding-left: 2em;
  margin-bottom: 1em;
}

.markdown-content :deep(code) {
  background-color: #f5f5f5;
  padding: 0.2em 0.4em;
  border-radius: 3px;
  font-family: monospace;
}

.markdown-content :deep(pre) {
  background-color: #f5f5f5;
  padding: 1em;
  border-radius: 4px;
  overflow-x: auto;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  padding-left: 2em;
  margin-top: 0;
  margin-bottom: 4px;
}

.markdown-body :deep(li) {
  margin-bottom: 4px;
  line-height: 1.5;
}

.markdown-body :deep(ul ul),
.markdown-body :deep(ol ol),
.markdown-body :deep(ul ol),
.markdown-body :deep(ol ul) {
  padding-left: 1.5em;
}

.markdown-body :deep(ul) {
  list-style-type: disc;
}

.markdown-body :deep(ol) {
  list-style-type: decimal;
}
</style>
