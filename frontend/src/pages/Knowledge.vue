<template>
  <v-container>
    <v-row class="mb-4" style="border-bottom: 1px solid #e0e0e0;">
      <v-col cols="12" class="d-flex justify-space-between align-center">
        <h1 class="text-h4" data-guide="knowledge-title">知识库</h1>
        <v-btn color="primary" @click="showUploadDialog = true" data-guide="upload-file">
          <v-icon left>mdi-plus</v-icon>
          上传文件
        </v-btn>
      </v-col>
    </v-row>

    <!-- Loading状态 -->
    <v-row v-if="loading">
      <v-col cols="12" class="text-center py-16">
        <v-progress-circular
          size="64"
          width="6"
          color="primary"
          indeterminate
        ></v-progress-circular>
        <p class="text-h6 grey--text text--lighten-1 mt-4">正在加载知识库...</p>
      </v-col>
    </v-row>

    <!-- 有文件时显示文件列表 -->
    <v-row v-else-if="files.files?.length > 0">
      <v-col
        v-for="file in files.files"
        :key="file.id"
        cols="12"
        sm="6"
        md="4"
        lg="3"
      >
        <v-card hover @click="showFileDetail(file)">
          <div class="card-title-wrapper">
            <v-card-title class="text-subtitle-1 two-line-title">
              {{ file.fileName }}
            </v-card-title>
          </div>
          <v-card-text class="text-caption text-truncate-3">
            {{ truncateText(file.description, 16) }}
          </v-card-text>
          <v-card-actions class="px-2 pb-3 pt-0">
            <v-btn
              icon
              small
              color="primary"
              @click.stop="downloadFile(file.ID)"
            >
              <v-icon small>mdi-download-outline</v-icon>
            </v-btn>
            
            <v-btn
              icon
              small
              color="error"
              @click.stop="deleteFile(file)"
            >
              <v-icon small>mdi-delete-outline</v-icon>
            </v-btn>

            <v-spacer></v-spacer>
            <div class="d-flex align-center">
              <v-chip small class="mr-2">{{ getFileType(file.fileName) }}</v-chip>
              <span class="text-caption grey--text mr-2">{{ formatDate(file.CreatedAt) }}</span>
            </div>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!-- 无文件时显示空状态 -->
    <v-row v-else>
      <v-col cols="12" class="text-center py-16">
        <v-icon size="64" color="grey lighten-1">mdi-file-document-outline</v-icon>
        <p class="text-h6 grey--text text--lighten-1 mt-4">知识库暂无内容</p>
        <p class="text-body-2 grey--text text--lighten-2 mt-2">点击上方"上传文件"按钮添加您的第一个文件</p>
      </v-col>
    </v-row>

    <v-dialog v-model="showUploadDialog" max-width="600">
      <v-card>
        <v-card-title>上传知识文件</v-card-title>
        <v-card-text>
          <v-form @submit.prevent="handleUpload" ref="uploadForm">
            <v-file-input
              v-model="form.file"
              label="选择文件"
              :rules="[
                v => !!v || '请选择文件',
                v => !v || /\.(pdf|docx?|txt|md|zip|pptx|xls|xlsx)$/.test(v.name) || '仅支持PDF/Word/txt/md/zip/pptx/xls/xlsx'
              ]"
              accept=".pdf,.doc,.docx,.txt,.md,.zip,.pptx,.xls,.xlsx"
              prepend-icon="mdi-paperclip"
              required
            ></v-file-input>
            <v-textarea
              v-model="form.description"
              label="简介"
              :rules="[
                 v => !v || v.length <= maxDescriptionLength || `不超过${maxDescriptionLength}字`
              ]"
              rows="3"
              counter
              :maxlength="maxDescriptionLength"
            >
              <template v-slot:counter>
                {{ form.description.length }}/{{ maxDescriptionLength }}
              </template>
            </v-textarea>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="showUploadDialog = false">取消</v-btn>
          <v-btn
            color="primary"
            :loading="uploading"
            @click="handleUpload"
          >
            上传
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

      <v-dialog v-model="showDetailDialog" max-width="900">
        <v-card class="detail-card">
          <v-card-title class="d-flex justify-space-between align-center">
              <span class="text-h5">{{ currentFile?.fileName }}</span>
          </v-card-title>
          <v-divider></v-divider>
          <v-card-text class="pa-0">
            <!-- Loading状态 -->
            <v-container v-if="detailLoading" fluid class="pa-0">
              <v-row no-gutters>
                <v-col cols="12" class="text-center py-16">
                  <v-progress-circular
                    size="48"
                    width="4"
                    color="primary"
                    indeterminate
                  ></v-progress-circular>
                  <p class="text-body-2 grey--text text--lighten-1 mt-4">正在加载文件详情...</p>
                </v-col>
              </v-row>
            </v-container>
            
            <!-- 内容显示 -->
            <v-container v-else fluid class="pa-0">
              <v-row no-gutters>
                <v-col cols="12" md="6" class="pa-4">
                  <div class="mb-4">
                    <p class="text-subtitle-1 font-weight-bold mb-1">文件描述</p>
                    <p>{{ currentFile?.description || '暂无描述' }}</p>
                  </div>

                  <div class="mb-4">
                    <div class="d-flex align-center">
                      <v-chip class="mr-2" color="primary" text-color="white">
                        {{ getFileType(currentFile?.fileName) }}
                      </v-chip>
                      <v-chip color="secondary" text-color="white">
                        {{ formatDate(currentFile?.CreatedAt) }}
                      </v-chip>
                    </div>
                  </div>
                </v-col>

                <v-col cols="12" md="6" class="pa-4 preview-col" style="background: #f5f5f5; min-height: 400px;">
                  <div v-if="currentFile?.content">
                    <p class="text-subtitle-1 font-weight-bold mb-2">文件预览</p>
                    <div v-if="isMarkdownFile(currentFile.fileName)" class="markdown-preview" v-html="renderMarkdown(currentFile.content)"></div>
                    <div v-else-if="isTextFile(currentFile.fileName)" class="text-preview">
                      <pre>{{ currentFile.content }}</pre>
                    </div>
                    <div v-else class="text-preview">
                        <pre>{{ currentFile.content }}</pre>
                    </div>
                  </div>
                  <div v-else class="text-center py-8">
                    <v-icon size="64" color="grey">mdi-file-remove</v-icon>
                    <p class="mt-2 grey--text">无预览内容</p>
                  </div>
                </v-col>
              </v-row>
            </v-container>
          </v-card-text>
          <v-divider></v-divider>
            <v-card-actions class="pa-4">
              <v-btn
                color="primary"
                variant="text"
                @click="downloadFile(currentFile.ID)"
                class="ml-2"
              >
                <v-icon left>mdi-download-outline</v-icon>
                下载文件
              </v-btn>
              
              <v-btn
                color="error"
                variant="text"
                @click="deleteFile(currentFile)"
                class="ml-2"
              >
                <v-icon left>mdi-delete-outline</v-icon>
                删除文件
              </v-btn>
              
              <v-spacer></v-spacer>
              <v-btn
                color="grey-darken-1"
                variant="text"
                @click="showDetailDialog = false; currentFile = null; detailLoading = false"
                class="ml-2"
              >
                关闭
              </v-btn>
            </v-card-actions>
        </v-card>
      </v-dialog>
      
      <!-- 删除确认对话框 -->
      <v-dialog v-model="showDeleteDialog" max-width="400">
        <v-card>
          <v-card-title class="text-h6">
            <v-icon color="warning" class="mr-2">mdi-alert-circle</v-icon>
            确认删除
          </v-card-title>
          <v-card-text>
            <p class="mb-2">确定要删除文件吗？</p>
            <p class="text-h6 font-weight-bold text-primary mb-2">{{ deleteFileData?.fileName }}</p>
            <p class="text-caption text-warning">此操作不可撤销，请谨慎操作。</p>
          </v-card-text>
          <v-card-actions class="px-4 pb-4">
            <v-spacer></v-spacer>
            <v-btn 
              text
              @click="showDeleteDialog = false"
              :disabled="deleting"
            >
              取消
            </v-btn>
            <v-btn
              color="error"
              @click="confirmDelete"
              :loading="deleting"
            >
              删除
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
  </v-container>
</template>


<script setup>
import { ref, onMounted } from 'vue'
import server from '@/api/server';
import {marked} from "marked";
import { ElMessage } from 'element-plus'

const maxTitleLength = 13
const maxDescriptionLength = 200

const files = ref([])
const loading = ref(false)

const showUploadDialog = ref(false)
const uploading = ref(false)
const uploadForm = ref()

const showDetailDialog = ref(false)
const currentFile = ref(null)
const detailLoading = ref(false)

const showDeleteDialog = ref(false)
const deleteFileData = ref(null)
const deleting = ref(false)

const form = ref({
    file: null,
    description: ''
})

const fetchFiles = async () => {
    try {
        loading.value = true
        const response = await server.Get('/knowledge/files')
        files.value = response.data
        console.log('获取文件列表:', files.value)
    } catch (error) {
        console.error('获取文件列表失败:', error)
    } finally {
        loading.value = false
    }
}

const truncateText = (text, maxLength) => {
    if (!text) return ''
    return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
}

const showFileDetail = async (file) => {
    try {
        showDetailDialog.value = true
        detailLoading.value = true
        currentFile.value = null
        const response = await server.Get(`/knowledge/file/${file.ID}`)
        currentFile.value = response.data
        console.log(response.data)
    } catch (error) {
        console.error('获取文件详情失败:', error)
        ElMessage.error('获取文件详情失败')
    } finally {
        detailLoading.value = false
    }
}

const handleUpload = async () => {
    console.log(uploadForm)
    const { valid } = await uploadForm.value.validate()
    if (!valid) return

    uploading.value = true

    try {
        const formData = new FormData()
        formData.append('file', form.value.file)
        formData.append('description', form.value.description)

        const response = await server.Post('/knowledge/upload', formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        })
        if(response.status === 200) {
          form.value = {file: null, description: ''}
          showUploadDialog.value = false
          await fetchFiles()
        } else {
          ElMessage.error("文件过大，无法上传")
        }
    } catch (error) {
        ElMessage.error("文件过大，无法上传")
        console.error('上传文件失败:', error)
    } finally {
        uploading.value = false
    }
}

const getFileType = (filename) => {
    if (!filename) return '未知'
    const ext = filename.split('.').pop().toLowerCase()
    switch (ext) {
        case 'pdf': return 'PDF'
        case 'doc': return 'DOC'
        case 'docx': return 'DOCX'
        case 'txt': return 'TXT'
        case 'md': return 'MD'
        case 'zip': return 'ZIP'
        case 'pptx': return 'PPTX'
        case 'xls': return 'XLS'
        case 'xlsx': return 'XLSX'
        default: return ext.toUpperCase()
    }
}

const downloadFile = async (fileId) => {
    try {
        const response = await server.Get(`/knowledge/download/${fileId}`, {
            responseType: 'blob'
        })

        let filename = 'download'
        const contentDisposition = response.headers['content-disposition']

        if (contentDisposition) {
            const utf8Filename = contentDisposition.match(/filename\*=UTF-8''(.+)/)?.[1]
            if (utf8Filename) {
                filename = decodeURIComponent(utf8Filename)
            } else {
                const match = contentDisposition.match(/filename="?(.+)"?/)
                if (match && match[1]) {
                    filename = match[1]
                    if (/[\x80-\xFF]/.test(filename)) {
                        filename = decodeURIComponent(escape(filename))
                    }
                }
            }
        }

        const url = window.URL.createObjectURL(new Blob([response.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', filename)
        document.body.appendChild(link)
        link.click()
        link.remove()

        ElMessage.success('文件下载成功')
    } catch (error) {
        console.error('下载文件失败:', error)
        ElMessage.error('下载文件失败: ' + error.message)
    }
}

const deleteFile = async (file) => {
    try {
        deleteFileData.value = file
        showDeleteDialog.value = true
    } catch (error) {
        console.error('删除文件失败:', error)
        ElMessage.error('删除文件失败: ' + error.message)
    }
}

const confirmDelete = async () => {
    deleting.value = true
    try {
        const response = await server.Delete(`/knowledge/file/${deleteFileData.value.ID}`)
        
        if (response.status === 200) {
            ElMessage.success('文件删除成功')
            // 刷新文件列表
            await fetchFiles()
            showDeleteDialog.value = false
            // 如果详情对话框打开，也关闭它
            if (showDetailDialog.value) {
                showDetailDialog.value = false
            }
        } else {
            ElMessage.error('删除失败')
        }
    } catch (error) {
        console.error('删除文件失败:', error)
        ElMessage.error('删除文件失败: ' + error.message)
    } finally {
        deleting.value = false
    }
}

const formatDate = (dateString) => {
    if (!dateString) return ''
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
    })
}

const isMarkdownFile = (filename) => {
    return filename?.endsWith('.md');
};

const isWordFile = (filename) => {
    return filename?.match(/\.docx?$/i)
}

const isTextFile = (filename) => {
    return filename?.endsWith('.txt');
};

function renderMarkdown(content) {
  return marked(content || '');
}


onMounted(() => {
    fetchFiles()
})
</script>

<style scoped>
.text-truncate-3 {
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    min-height: 72px;
    overflow: hidden;
    line-height: 1.4;
}

.v-card {
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    height: 100%;
    display: flex;
    flex-direction: column;
}

.v-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.1) !important;
}

.v-card-text {
    flex-grow: 1;
    padding-bottom: 8px;
}

.detail-card {
  transition: none !important;
  transform: none !important;
  box-shadow: none !important;
}

.action-btn {
    transition: all 0.2s ease;
    opacity: 0.7;
    margin: 0 4px;
}

.action-btn:hover {
    opacity: 1;
    transform: scale(1.15);
}

.empty-state {
    opacity: 0.7;
    transition: opacity 0.3s ease;
}

.empty-state:hover {
    opacity: 0.9;
}

.card-title-wrapper {
    border-bottom: 1px solid #e0e0e0;
    padding: 0 16px;
    margin-bottom: 8px;
}

.two-line-title {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    min-height: 48px;
    line-height: 1.3;
    padding: 12px 0;
    white-space: normal;
    word-break: break-word;
}

.v-chip {
    height: 22px !important;
    font-size: 0.65rem !important;
}

.v-card-actions {
    min-height: 40px !important;
}

.file-type-chip {
    font-size: 0.75rem;
    height: 24px;
}

.file-detail-content {
    max-height: 60vh;
    overflow-y: auto;
    padding-right: 8px;
}

.file-detail-content::-webkit-scrollbar {
    width: 6px;
}

.file-detail-content::-webkit-scrollbar-thumb {
    background-color: rgba(0, 0, 0, 0.2);
    border-radius: 3px;
}

.upload-btn {
    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.upload-btn:hover {
    box-shadow: 0 3px 6px rgba(0, 0, 0, 0.16);
}

.v-card-actions .v-btn--icon {
    opacity: 0.7;
    transition: opacity 0.2s ease;
}

.v-card-actions .v-btn--icon:hover {
    opacity: 1;
}

.v-input :deep(.v-label) {
    font-weight: 500;
    color: rgba(0, 0, 0, 0.8);
}

.date-chip {
    font-size: 0.7rem;
}

.preview-col {
    border-left: 1px solid #e0e0e0;
}

.markdown-preview {
    max-height: 500px;
    overflow-y: auto;
    padding: 8px;
    background: white;
    border-radius: 4px;
}

.markdown-preview .katex {
    font-size: 1.1em;
}

.markdown-preview >>> pre {
    background: #f5f5f5;
    padding: 12px;
    border-radius: 4px;
    overflow-x: auto;
}

.text-preview {
    max-height: 500px;
    overflow-y: auto;
    padding: 8px;
    background: white;
    border-radius: 4px;
}

.text-preview pre {
    white-space: pre-wrap;
    font-family: monospace;
    margin: 0;
}

.unsupported-preview {
    background: #f5f5f5;
    border-radius: 4px;
    height: 100%;
}
</style>

