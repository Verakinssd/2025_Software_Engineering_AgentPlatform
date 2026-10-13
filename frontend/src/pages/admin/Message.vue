<template>
  <v-container>
    <v-stepper :items="['选择类型', '编辑内容', '确认发送']" v-model="currentStep">
      <template v-slot:item.1>
        <v-card title="选择公告类型" flat>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <v-radio-group v-model="announcementType">
                  <v-radio label="广播（发送给所有用户）" value="broadcast"></v-radio>
                  <v-radio label="单点（选择特定用户）" value="single"></v-radio>
                </v-radio-group>
              </v-col>

              <v-col cols="12" md="6" v-if="announcementType === 'single'">
                <v-card variant="outlined">
                  <v-card-title>选择接收用户</v-card-title>
                  <v-card-text>
                    <v-text-field
                      v-model="userSearch"
                      label="搜索用户"
                      prepend-inner-icon="mdi-magnify"
                      variant="outlined"
                      clearable
                    ></v-text-field>

                    <v-list lines="two" class="overflow-y-auto" max-height="300">
                      <v-list-item
                        v-for="user in filteredUsers"
                        :key="user.id"
                      >
                        <template v-slot:prepend>
                          <v-checkbox
                            v-model="selectedUsers"
                            :value="user.id"
                          ></v-checkbox>
                          <v-avatar size="36" class="mr-2">
                            <v-img
                              :src="user.avatar || '/public/images/default_avatar.jpg'"
                              :alt="user.username"
                            />
                          </v-avatar>
                        </template>

                        <v-list-item-title>{{ user.username }}</v-list-item-title>
                        <v-list-item-subtitle>{{ user.email }}</v-list-item-subtitle>
                      </v-list-item>
                    </v-list>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </template>

      <template v-slot:item.2>
        <v-card title="编辑公告内容" flat>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="announcementTitle"
                  label="公告标题"
                  variant="outlined"
                  class="mb-4"
                ></v-text-field>

                <v-toolbar density="compact" class="elevation-1">
                  <v-btn icon @click="insertMarkdown('bold')" title="加粗">
                    <v-icon>mdi-format-bold</v-icon>
                  </v-btn>
                  <v-btn icon @click="insertMarkdown('italic')" title="斜体">
                    <v-icon>mdi-format-italic</v-icon>
                  </v-btn>
                  <v-btn icon @click="insertMarkdown('header')" title="标题">
                    <v-icon>mdi-format-header-1</v-icon>
                  </v-btn>
                  <v-btn icon @click="insertMarkdown('header-2')" title="二级标题">
                    <v-icon>mdi-format-header-2</v-icon>
                  </v-btn>
                  <v-divider vertical inset class="mx-2"></v-divider>

                  <v-btn icon @click="insertMarkdown('unorderedList')" title="无序列表">
                    <v-icon>mdi-format-list-bulleted</v-icon>
                  </v-btn>
                  <v-btn icon @click="insertMarkdown('table')" title="表格">
                    <v-icon>mdi-table</v-icon>
                  </v-btn>
                  <v-btn icon @click="insertMarkdown('hr')" title="水平分割线">
                    <v-icon>mdi-minus</v-icon>
                  </v-btn>

                  <v-divider vertical inset class="mx-2"></v-divider>

                  <v-btn icon @click="insertMarkdown('link')" title="链接">
                    <v-icon>mdi-link</v-icon>
                  </v-btn>
                  <v-btn icon @click="insertMarkdown('code')" title="代码块">
                    <v-icon>mdi-code-tags</v-icon>
                  </v-btn>
                  <v-btn icon @click="showImageDialog = true" title="插入图片">
                    <v-icon>mdi-image</v-icon>
                  </v-btn>
                </v-toolbar>

                <v-dialog v-model="showImageDialog" max-width="500">
                  <v-card>
                    <v-card-title>上传图片</v-card-title>
                    <v-card-text>
                      <v-file-input
                        v-model="imageFile"
                        accept="image/*"
                        label="选择图片文件"
                        prepend-icon="mdi-image"
                        variant="outlined"
                        @update:modelValue="handleImageSelect"
                        :loading="isUploading"
                        :disabled="isUploading"
                        :rules="[validateImageFile]"
                      ></v-file-input>

                      <v-progress-linear
                        v-if="uploadProgress > 0 && uploadProgress < 100"
                        :model-value="uploadProgress"
                        color="primary"
                        height="8"
                        striped
                      ></v-progress-linear>

                    </v-card-text>
                    <v-card-actions>
                      <v-spacer></v-spacer>
                      <v-btn @click="cancelImageUpload">取消</v-btn>
                      <v-btn
                        color="primary"
                        @click="uploadImage"
                        :disabled="!imageFile || isUploading"
                        :loading="isUploading"
                      >
                        插入图片
                      </v-btn>
                    </v-card-actions>
                  </v-card>
                </v-dialog>

                <v-textarea
                  v-model="announcementContent"
                  label="公告内容 (支持Markdown格式)"
                  variant="outlined"
                  rows="10"
                  auto-grow
                ></v-textarea>
              </v-col>

              <v-col cols="12" md="6">
                <v-card variant="outlined">
                  <v-card-title>预览</v-card-title>
                  <v-card-text>
                    <div v-if="announcementTitle" class="text-h5 mb-4">{{ announcementTitle }}</div>
                    <v-divider class="mb-4"></v-divider>
                    <div v-html="compiledMarkdown"></div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </template>

      <template v-slot:item.3>
        <v-card title="确认发送" flat>
          <v-card-text>
            <div class="text-h5 mb-4">{{ announcementTitle }}</div>
            <v-divider class="mb-4"></v-divider>

            <div class="mb-6" v-html="compiledMarkdown"></div>

            <div class="text-subtitle-1 mb-2">发送给：</div>
            <div class="mb-6">
              <v-chip v-if="announcementType === 'broadcast'" color="primary" class="mr-2">
                所有用户
              </v-chip>
              <template v-else>
                <v-chip
                  v-for="userId in selectedUsers"
                  :key="userId"
                  color="primary"
                  class="mr-2 mb-2"
                >
                  {{ getUserName(userId) }}
                </v-chip>
              </template>
            </div>

            <v-btn color="primary" @click="sendAnnouncement">确认发送</v-btn>
          </v-card-text>
        </v-card>
      </template>
    </v-stepper>
  </v-container>
</template>

<script setup>
import {ref, computed, onMounted, nextTick} from 'vue'
import { marked } from 'marked'
import server from '@/api/server'
import {ElMessage} from "element-plus";

const currentStep = ref(1)
const announcementType = ref('broadcast')
const userSearch = ref('')
const selectedUsers = ref([])
const users = ref([])
const announcementTitle = ref('')
const announcementContent = ref('')
const showImageDialog = ref(false);
const imageUrl = ref('');

const imageFile = ref(null);
const imagePreview = ref('');
const uploadProgress = ref(0);
const isUploading = ref(false);

marked.setOptions({
  breaks: true,
  gfm: true,
  sanitize: false
})

const filteredUsers = computed(() => {
  if (!userSearch.value) return users.value
  const search = userSearch.value.toLowerCase()
  return users.value.filter(user =>
    user.username.toLowerCase().includes(search) ||
    user.email.toLowerCase().includes(search)
)
})

const compiledMarkdown = computed(() => {
  return marked(announcementContent.value || '*请输入公告内容*')
})

const fetchUsers = async () => {
  try {
    const response = await server.Get('/admin/all_user')
    console.log(response.data)
    users.value = response.data || []
  } catch (error) {
    console.error('Error fetching users:', error)
  }
}

const getUserName = (userId) => {
  const user = users.value.find(u => u.id === userId)
  return user ? user.username : ''
}

// const insertImage = () => {
//   if (imageUrl.value) {
//     const markdownImage = `![](${imageUrl.value})`;
//
//     const textarea = document.querySelector('textarea');
//     const start = textarea.selectionStart;
//     announcementContent.value =
//       announcementContent.value.substring(0, start) +
//       markdownImage +
//       announcementContent.value.substring(start);
//
//     imageUrl.value = '';
//     showImageDialog.value = false;
//   }
// };

const validateImageFile = (file) => {
    if (!file) return true

    const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
    const isImage = validTypes.includes(file.type)
    //const isLt5M = file.size / 1024 / 1024 < 5

    if (!isImage) return '请上传JPEG/PNG/GIF/WEBP格式的图片'
    //if (!isLt5M) return '图片大小不能超过5MB'
    return true
}

const handleImageSelect = async (file) => {
    if (!file) {
        imagePreview.value = ''
        return
    }

    const targetWidth = 750

    const image = new Image()
    image.src = URL.createObjectURL(file)

    image.onload = async () => {
        const scale = targetWidth / image.width
        const canvas = document.createElement('canvas')
        canvas.width = targetWidth
        canvas.height = image.height * scale

        const ctx = canvas.getContext('2d')
        ctx.drawImage(image, 0, 0, canvas.width, canvas.height)

        canvas.toBlob((blob) => {
            if (blob) {
                const compressedFile = new File([blob], file.name, { type: blob.type })
                imageFile.value = compressedFile

                const reader = new FileReader()
                reader.onload = (e) => {
                    imagePreview.value = e.target.result
                }
                reader.readAsDataURL(compressedFile)
            }
        }, 'image/jpeg', 1)
    }
}

const cancelImageUpload = () => {
    imageFile.value = null
    imagePreview.value = ''
    uploadProgress.value = 0
    showImageDialog.value = false
}

const uploadImage = async () => {
    if (!imageFile.value) return

    isUploading.value = true
    uploadProgress.value = 0

    try {
        const formData = new FormData()
        formData.append('image', imageFile.value)

        const response = await server.Post('/image/upload', formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            },
            onUploadProgress: (progressEvent) => {
                if (progressEvent.total) {
                    uploadProgress.value = Math.round(
                        (progressEvent.loaded * 100) / progressEvent.total
                    )
                }
            }
        })

        if (!response.data?.url) {
            throw new Error('服务器返回数据异常')
        }

        const markdownImage = `![图片描述](${response.data.url})`
        insertToMarkdown(markdownImage)

        cancelImageUpload()
    } catch (error) {
        console.error('图片上传失败:', error)
        ElMessage.error(`上传失败: ${error.response?.data?.message || error.message}`)
    } finally {
        isUploading.value = false
    }
}

const insertToMarkdown = (content) => {
    const textarea = document.querySelector('textarea')
    const start = textarea.selectionStart
    const end = textarea.selectionEnd

    announcementContent.value =
        announcementContent.value.substring(0, start) +
        content +
        announcementContent.value.substring(end)

    nextTick(() => {
        textarea.focus()
        const newPos = start + content.length
        textarea.setSelectionRange(newPos, newPos)
    })
}

const insertMarkdown = (type, imageUrl='') => {
  const textarea = document.querySelector('textarea')
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = announcementContent.value.substring(start, end)
  let insertText = ''

  switch(type) {
    case 'bold':
      insertText = `**${selectedText}**`
      break
    case 'italic':
      insertText = `*${selectedText}*`
      break
    case 'header':
      insertText = `# ${selectedText}`
      break
    case 'header-2':
      insertText = `## ${selectedText}`
      break
    case 'list':
      insertText = `- ${selectedText}`
      break
    case 'link':
      insertText = `[${selectedText}](url)`
      break
    case 'code':
      insertText = "```\n" + selectedText + "\n```"
      break
    case 'unorderedList':
      insertText = `- ${selectedText}`
      break
    case 'table':
      insertText = `| 标题1 | 标题2 | 标题3 |\n|------|------|------|\n| 内容1 | 内容2 | 内容3 |\n| 内容4 | 内容5 | 内容6 |`
      break
    case 'hr':
      insertText = `\n---\n`
      break
    case 'image':
      insertText = `![图片描述](${imageUrl})`
      break
  }

  announcementContent.value =
    announcementContent.value.substring(0, start) +
    insertText +
    announcementContent.value.substring(end)

  nextTick(() => {
    textarea.focus()
    if (type !== 'image') {
      textarea.setSelectionRange(start, start + insertText.length)
    }
  })
}

const sendAnnouncement = async () => {
  try {
    if (announcementType.value === 'broadcast') {
      const response = await server.Post('/notification/broadcast', {
        title: announcementTitle.value,
        content: announcementContent.value
      });

      console.log('广播发送成功:', response.data);
      ElMessage.success('公告广播发送成功！');
    } else {
      if (selectedUsers.value.length === 0) {
        ElMessage.warning('请至少选择一个用户');
        return;
      }
      const promises = selectedUsers.value.map(async userId => {
        await server.Post('/notification/send', {
          title: announcementTitle.value,
          content: announcementContent.value,
          to_user_id: userId
        });
      });

      await Promise.all(promises);
      console.log('单点通知全部发送成功');
      ElMessage.success(`公告已成功发送给 ${selectedUsers.value.length} 位用户`);
    }

    resetForm();
  } catch (error) {
    console.error('发送失败:', error);
    ElMessage.error(`发送失败: ${error.response?.data?.message || error.message}`);
  }
}

const resetForm = () => {
  currentStep.value = 1
  announcementType.value = 'broadcast'
  selectedUsers.value = []
  announcementTitle.value = ''
  announcementContent.value = ''
  userSearch.value = ''
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.markdown-preview h1, .markdown-preview h2, .markdown-preview h3 {
  margin-top: 1em;
  margin-bottom: 0.5em;
}
.markdown-preview p {
  margin-bottom: 1em;
}
.markdown-preview ul, .markdown-preview ol {
  margin-left: 1.5em;
  margin-bottom: 1em;
}
.markdown-preview code {
  background-color: #f5f5f5;
  padding: 0.2em 0.4em;
  border-radius: 3px;
}
.markdown-preview pre {
  background-color: #f5f5f5;
  padding: 1em;
  border-radius: 4px;
  overflow-x: auto;
}

.markdown-preview img {
  max-width: 40px;
  width: auto;
  height: auto;
  margin: 1em 0;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
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
