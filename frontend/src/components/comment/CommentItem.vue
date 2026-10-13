<template>
  <div class="comment-item mb-4">
    <div class="d-flex">
      <v-avatar size="36" class="mr-3">
        <v-img :src="comment.author.avatar" :lazy-src="'/images/default_avatar.jpg'" />
      </v-avatar>
      <div class="flex-grow-1">
        <div class="d-flex align-center mb-1">
          <span class="text-subtitle-2 font-weight-medium">{{ comment.author.name }}</span>
        </div>

        <div
          class="markdown-body mb-2"
          v-html="renderMarkdown(getDisplayContent(comment))"
        />

        <div class="d-flex align-center mb-2">
          <span class="text-caption text-medium-emphasis mr-2">
            {{ formatTime(comment.createdAt) }}
          </span>
          <v-btn variant="text" size="small" class="text-caption" @click="$emit('toggle-reply', comment.id)">
            回复
          </v-btn>
          <v-btn
            v-if="canDelete(comment)"
            variant="text"
            size="small"
            color="error"
            class="text-caption ml-1"
            @click.stop="$emit('confirm-delete', comment.id)"
          >
            删除
          </v-btn>
          <v-btn
            variant="text"
            size="small"
            class="text-caption ml-1"
            :color="comment.isLiked ? 'red' : ''"
            :prepend-icon="comment.isLiked ? 'mdi-heart' : 'mdi-heart-outline'"
            @click.stop="$emit('handle-like', comment)"
          >
            {{ comment.likes }}
          </v-btn>
        </div>

        <div v-if="activeReply === comment.id" class="d-flex mb-4 pl-4">
          <v-textarea
            v-model="replyContents[comment.id]"
            label="写下你的回复..."
            variant="outlined"
            rows="1"
            auto-grow
            hide-details
            class="flex-grow-1 mr-2"
          />
          <v-btn color="white" height="40" class="align-self-center" @click="$emit('submit-reply', comment.id)">
            回复
          </v-btn>
        </div>

        <div v-if="sortedReplies.length" class="pl-4">
          <CommentItem
            v-for="child in displayedReplies"
            :key="child.id"
            :comment="child"
            :is-root="false"
            :render-markdown="renderMarkdown"
            :format-time="formatTime"
            :can-delete="canDelete"
            :active-reply="activeReply"
            :reply-contents="replyContents"
            @toggle-reply="$emit('toggle-reply', $event)"
            @submit-reply="$emit('submit-reply', $event)"
            @confirm-delete="$emit('confirm-delete', $event)"
            @handle-like="$emit('handle-like', $event)"
          />
          
          <!-- 显示更多/收起按钮 -->
          <div v-if="sortedReplies.length > 3" class="mt-2">
            <v-btn 
              variant="text" 
              size="small" 
              color="primary" 
              class="text-caption"
              @click="toggleShowAllReplies"
            >
              {{ showAllReplies ? `收起回复` : `查看更多回复 (${sortedReplies.length - 3})` }}
            </v-btn>
          </div>
        </div>
      </div>
    </div>
    <v-divider v-if="isRoot" class="my-4" />
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import CommentItem from './CommentItem.vue'

defineOptions({ name: 'CommentItem' })

interface Author {
  id: string
  name: string
  avatar: string
}

interface Comment {
  id: string
  content: string
  likes: number
  isLiked?: boolean
  createdAt: Date
  author: Author
  parentId?: string
  replyTo?: string
}

interface CommentWithReplies extends Comment {
  replies?: Comment[]
}

const props = defineProps<{
  comment: CommentWithReplies & { replyToUserName?: string }
  isRoot: boolean
  renderMarkdown: (text: string) => string
  formatTime: (d: Date) => string
  canDelete: (c: CommentWithReplies) => boolean
  activeReply: string
  replyContents: Record<string, string>
}>()

defineEmits<{
  (e: 'toggle-reply', id: string): void
  (e: 'submit-reply', id: string): void
  (e: 'confirm-delete', id: string): void
  (e: 'handle-like', comment: CommentWithReplies): void
}>()

// 控制是否显示所有回复
const showAllReplies = ref(false)

function getDisplayContent(comment: CommentWithReplies & { replyToUserName?: string }): string {
  if (comment.parentId && comment.replyTo && comment.parentId !== comment.replyTo) {
    const username = comment.replyToUserName || '用户'
    return `回复 @${username}：${comment.content}`
  }
  return comment.content
}

const sortedReplies = computed(() => {
  const replies = props.comment.replies || []
  return [...replies].sort((a, b) =>
    new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
  )
})

// 控制显示的回复数量
const displayedReplies = computed(() => {
  const maxDisplay = 3
  if (showAllReplies.value || sortedReplies.value.length <= maxDisplay) {
    return sortedReplies.value
  }
  return sortedReplies.value.slice(0, maxDisplay)
})

// 切换显示状态
const toggleShowAllReplies = () => {
  showAllReplies.value = !showAllReplies.value
}
</script>

<style scoped>
.comment-item {
}
</style>

<style scoped>
.markdown-body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
  font-size: 16px;
  line-height: 1.5;
  word-wrap: break-word;
  white-space: pre-wrap;
  overflow-wrap: break-word;
  word-break: break-word;
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3),
.markdown-body :deep(h4),
.markdown-body :deep(h5),
.markdown-body :deep(h6) {
  margin-top: 24px;
  margin-bottom: 4px;
  font-weight: 600;
  line-height: 1.25;
}

.markdown-body :deep(p) {
  margin-top: 0;
  margin-bottom: 4px;
}

.markdown-body :deep(a) {
  color: #1976D2;
  text-decoration: none;
}

.markdown-body :deep(a:hover) {
  text-decoration: underline;
}

.markdown-body :deep(code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(27, 31, 35, 0.05);
  border-radius: 3px;
}

.markdown-body :deep(pre) {
  word-wrap: normal;
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f6f8fa;
  border-radius: 3px;
}

.markdown-body :deep(blockquote) {
  padding: 0 1em;
  color: #6a737d;
  border-left: 0.25em solid #dfe2e5;
  margin: 0 0 4px 0;
}

.markdown-body :deep(img) {
  max-width: 100%;
  box-sizing: initial;
  background-color: #fff;
}
</style>