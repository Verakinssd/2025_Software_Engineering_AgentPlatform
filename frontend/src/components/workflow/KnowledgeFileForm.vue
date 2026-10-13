<template>
	<v-card-title class="d-flex align-center">
		<v-btn icon flat dense :ripple="false" @click="showForm_file = !showForm_file" class="mr-2">
			<v-icon>{{ showForm_file ? 'mdi-chevron-down' : 'mdi-chevron-right' }}</v-icon>
		</v-btn>
		<span class="text-subtitle-1">知识库文件</span>
		<v-tooltip text="选择要检索的知识库文件" :attach="true">
			<template #activator="{ props }">
				<v-icon v-bind="props" size="20" class="ms-2">
					mdi-information-outline
				</v-icon>
			</template>
		</v-tooltip>
	</v-card-title>

	<v-expand-transition v-if="showForm_file">
		<v-card-text v-if="loading" class="text-center py-4">
			<v-progress-circular indeterminate color="primary"></v-progress-circular>
			<div class="mt-2">加载中...</div>
		</v-card-text>
		<v-card-text v-else>
			<el-select 
				v-model="selectedFileId" 
				placeholder="选择知识库文件"
				style="width: 100%;"
				large
				:disabled="readonly"
				@change="updateKnowledgeFileId"
			>
				<el-option
					v-for="file in knowledgeFiles"
					:key="file.ID"
					:label="file.fileName"
					:value="file.ID"
				>
					<div style="display: flex; justify-content: space-between; align-items: center">
						<span>{{ file.fileName }}</span>
						<span style="color: #8492a6; font-size: 13px">{{ formatSize(file.fileSize) }}</span>
					</div>
				</el-option>
			</el-select>
		</v-card-text>
	</v-expand-transition>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import server from '@/api/server'
import { ElMessage } from 'element-plus'
import { selectedNode } from '@/api/workflow'

const props = defineProps({
	readonly: {
		type: Boolean,
		default: false
	}
})

const showForm_file = ref(true)
const knowledgeFiles = ref([])
const selectedFileId = ref(null)
const loading = ref(false)

onMounted(async () => {
	loading.value = true;
  await fetchKnowledgeFiles();
	loading.value = false;

  if (selectedNode.value) {
    if (selectedNode.value.data.knowledgeFileId) {
      selectedFileId.value = selectedNode.value.data.knowledgeFileId;
    }
  }
})

function formatSize(size) {
  if (size < 1024) {
    return size + ' B'
  } else if (size < 1024 * 1024) {
    return (size / 1024).toFixed(2) + ' KB'
  } else {
    return (size / (1024 * 1024)).toFixed(2) + ' MB'
  }
}

async function fetchKnowledgeFiles() {
  try {
    const response = await server.Get('/knowledge/files');
    
    if (response.status === 200) {
      // 提取文件数据
      if (response.data && response.data.files && Array.isArray(response.data.files)) {
        knowledgeFiles.value = response.data.files;
      } else if (Array.isArray(response.data)) {
        knowledgeFiles.value = response.data;
      } else if (response.data && typeof response.data === 'object') {
        const possibleArrays = Object.values(response.data).filter(val => Array.isArray(val));
        if (possibleArrays.length > 0) {
          knowledgeFiles.value = possibleArrays[0];
        } else {
          knowledgeFiles.value = [];
        }
      } else {
        knowledgeFiles.value = [];
      }
    } else {
      ElMessage.error('获取知识库文件失败');
    }
  } catch (error) {
    ElMessage.error('获取知识库文件失败');
  }
}

function updateKnowledgeFileId() {
  if (selectedNode.value && selectedFileId.value) {
    selectedNode.value.data.knowledgeFileId = selectedFileId.value;
  }
}
</script>

<style scoped>
.btn-add {
  opacity: 0.6;
  transition: opacity 0.2s;
}

.btn-add:hover {
  opacity: 1;
}

:deep(.v-btn--icon.v-btn--density-default) {
  width: calc(var(--v-btn-height) + 2px);
  height: calc(var(--v-btn-height) + 2px);
}
</style> 