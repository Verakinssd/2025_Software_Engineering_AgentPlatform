<template>
  <v-card-title class="d-flex align-center">
    <v-btn icon flat dense :ripple="false" @click="showWorkflowSelector = !showWorkflowSelector" class="mr-2">
      <v-icon>{{ showWorkflowSelector ? 'mdi-chevron-down' :'mdi-chevron-right' }}</v-icon> 
    </v-btn>
    <span class="text-subtitle-1">选择工作流</span>
    <v-tooltip :attach="true" text="选择要嵌入的工作流">
      <template #activator="{ props }">
        <v-icon v-bind="props" size="20" class="ms-2">
          mdi-information-outline
        </v-icon>
      </template>
    </v-tooltip>
    <!-- 只读模式提示 -->
    <v-spacer />
    <v-chip v-if="readonly && showWorkflowSelector" size="small" color="grey" variant="outlined">
      <v-icon size="14" class="mr-1">mdi-eye</v-icon>
      只读
    </v-chip>
  </v-card-title>

  <v-expand-transition v-if="showWorkflowSelector">
    <v-card class="mx-4 mb-4 pa-4" variant="outlined">
      <!-- 只读模式下显示警告信息 -->
      <v-alert 
        v-if="readonly" 
        type="warning" 
        variant="tonal" 
        class="mb-4"
        density="compact"
      >
        <template v-slot:prepend>
          <v-icon>mdi-alert-circle-outline</v-icon>
        </template>
        <span class="text-caption">该工作流引用的子工作流可能无法访问，仅供查看</span>
      </v-alert>
      
      <v-tabs v-model="activeTab" color="primary" :disabled="readonly">
        <v-tab value="personal">个人工作流</v-tab>
        <v-tab value="public">公开工作流</v-tab>
      </v-tabs>

      <v-window v-model="activeTab">
        <v-window-item value="personal">
          <v-card-text v-if="personalWorkflows.length === 0 && !loading" class="text-center py-4">
            暂无个人工作流
          </v-card-text>
          <v-card-text v-else-if="loading" class="text-center py-4">
            <v-progress-circular indeterminate color="primary"></v-progress-circular>
            <div class="mt-2">加载中...</div>
          </v-card-text>
          <v-list v-else>
            <v-list-item
              v-for="(workflow, index) in personalWorkflows"
              :key="workflow.id || index"
              :value="workflow.id"
              @click="!readonly && selectWorkflow(workflow)"
              :active="selectedWorkflowId === workflow.id"
              :disabled="readonly"
            >
              <template v-slot:prepend>
                <v-icon color="primary" :class="{ 'text-primary': selectedWorkflowId === workflow.id }">
                  {{ selectedWorkflowId === workflow.id ? 'mdi-check-circle' : 'mdi-application-outline' }}
                </v-icon>
              </template>
              <v-list-item-title>{{ workflow.name || '未命名工作流' }}</v-list-item-title>
              <v-list-item-subtitle>{{ workflow.description || '暂无描述' }}</v-list-item-subtitle>
            </v-list-item>
          </v-list>
        </v-window-item>

        <v-window-item value="public">
          <v-card-text v-if="publicWorkflows.length === 0 && !loading" class="text-center py-4">
            暂无公开工作流
          </v-card-text>
          <v-card-text v-else-if="loading" class="text-center py-4">
            <v-progress-circular indeterminate color="primary"></v-progress-circular>
            <div class="mt-2">加载中...</div>
          </v-card-text>
          <v-list v-else>
            <v-list-item
              v-for="(workflow, index) in publicWorkflows"
              :key="workflow.id || index"
              :value="workflow.id"
              @click="!readonly && selectWorkflow(workflow)"
              :active="selectedWorkflowId === workflow.id"
              :disabled="readonly"
            >
              <template v-slot:prepend>
                <v-icon color="primary" :class="{ 'text-primary': selectedWorkflowId === workflow.id }">
                  {{ selectedWorkflowId === workflow.id ? 'mdi-check-circle' : 'mdi-application-outline' }}
                </v-icon>
              </template>
              <v-list-item-title>{{ workflow.name || '未命名工作流' }}</v-list-item-title>
              <v-list-item-subtitle>{{ workflow.description || '暂无描述' }}</v-list-item-subtitle>
            </v-list-item>
          </v-list>
        </v-window-item>
      </v-window>

      <div v-if="selectedWorkflow" class="mt-4 pa-2 rounded bg-grey-lighten-4">
        <div class="d-flex align-center">
          <v-icon color="primary" class="mr-2">mdi-check-circle</v-icon>
          <div>
            <div class="text-subtitle-2">已选择工作流: {{ selectedWorkflow.name }}</div>
            <div class="text-caption">{{ selectedWorkflow.description || '暂无描述' }}</div>
          </div>
        </div>
      </div>
    </v-card>
  </v-expand-transition>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { selectedNode } from '@/api/workflow'
import { getPersonalWorkflows, getPublicWorkflows, getWorkflowIOStructure } from '@/api/workflow'
import { ElMessage } from 'element-plus'

// 定义props
const props = defineProps({
  readonly: {
    type: Boolean,
    default: false
  }
})

const showWorkflowSelector = ref(true)
const activeTab = ref('personal')
const loading = ref(false)
const personalWorkflows = ref<any[]>([])
const publicWorkflows = ref<any[]>([])
const selectedWorkflowId = ref<number | null>(null)
const selectedWorkflow = ref<any>(null)

// 初始化时获取工作流列表
onMounted(async () => {
  await loadWorkflows()
  
  // 如果节点已有选择的工作流，则设置选中状态
  if (selectedNode.value?.data.workflowId) {
    selectedWorkflowId.value = selectedNode.value.data.workflowId
    
    // 在个人和公开工作流中查找已选择的工作流
    const foundInPersonal = personalWorkflows.value.find(wf => wf.id === selectedWorkflowId.value)
    const foundInPublic = publicWorkflows.value.find(wf => wf.id === selectedWorkflowId.value)
    selectedWorkflow.value = foundInPersonal || foundInPublic || null
    
    // 设置正确的tab
    if (foundInPersonal) {
      activeTab.value = 'personal'
    } else if (foundInPublic) {
      activeTab.value = 'public'
    }
  }
})

// 监听标签页切换，切换时加载对应的工作流列表
watch(activeTab, async () => {
  await loadWorkflows()
})

// 加载工作流列表
async function loadWorkflows() {
  loading.value = true
  try {
    if (activeTab.value === 'personal') {
      const res = await getPersonalWorkflows()
      if (res.status === 200) {
        console.log('个人工作流数据:', res.data)
        // 检查数据结构，可能返回了嵌套对象
        if (Array.isArray(res.data)) {
          personalWorkflows.value = res.data
        } else if (res.data && Array.isArray(res.data.workflows)) {
          personalWorkflows.value = res.data.workflows
        } else if (res.data && typeof res.data === 'object') {
          // 如果是对象，尝试提取出数组
          const possibleArrays = Object.values(res.data).filter(Array.isArray)
          if (possibleArrays.length > 0) {
            personalWorkflows.value = possibleArrays[0]
          } else {
            personalWorkflows.value = []
            console.error('无法从返回数据中提取工作流数组:', res.data)
          }
        } else {
          personalWorkflows.value = []
        }
      } else {
        ElMessage.error('获取个人工作流失败')
      }
    } else {
      const res = await getPublicWorkflows()
      if (res.status === 200) {
        console.log('公开工作流数据:', res.data)
        // 检查数据结构，可能返回了嵌套对象
        if (Array.isArray(res.data)) {
          publicWorkflows.value = res.data
        } else if (res.data && Array.isArray(res.data.workflows)) {
          publicWorkflows.value = res.data.workflows
        } else if (res.data && typeof res.data === 'object') {
          // 如果是对象，尝试提取出数组
          const possibleArrays = Object.values(res.data).filter(Array.isArray)
          if (possibleArrays.length > 0) {
            publicWorkflows.value = possibleArrays[0]
          } else {
            publicWorkflows.value = []
            console.error('无法从返回数据中提取工作流数组:', res.data)
          }
        } else {
          publicWorkflows.value = []
        }
      } else {
        ElMessage.error('获取公开工作流失败')
      }
    }
  } catch (error) {
    ElMessage.error('加载工作流列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 选择工作流
async function selectWorkflow(workflow: any) {
  selectedWorkflowId.value = workflow.id
  selectedWorkflow.value = workflow
  
  // 更新节点数据
  if (selectedNode.value) {
    selectedNode.value.data.workflowId = workflow.id
    
    // 获取工作流的输入输出结构
    try {
      const res = await getWorkflowIOStructure(workflow.id)
      if (res.status === 200 && res.data) {
        if (res.data.inputs && Array.isArray(res.data.inputs)) {
          selectedNode.value.data.input = res.data.inputs;
        }
        
        if (res.data.outputs && Array.isArray(res.data.outputs)) {
          selectedNode.value.data.output = res.data.outputs;
        }
        
        console.log('已选择工作流，已更新节点数据:', selectedNode.value.data)
        if (selectedNode.value.data.input.length > 0 && selectedNode.value.data.output.length > 0) {
          ElMessage.success(`已选择工作流: ${workflow.name}`)
        } else {
          ElMessage.warning(`已选择工作流，该工作流未配置完整输入输出: ${workflow.name}`)
        }
      } else {
        ElMessage.warning(`已选择工作流，但获取输入输出失败: ${workflow.name}`)
      }
    } catch (error) {
      console.error('获取工作流输入输出结构失败:', error)
      ElMessage.warning(`已选择工作流，但获取输入输出失败: ${workflow.name}`)
    }
  }
}

// 监听节点变化，更新选中的工作流
watch(() => selectedNode.value, (newNode) => {
  if (newNode && newNode.data.workflowId) {
    selectedWorkflowId.value = newNode.data.workflowId
    
    // 在工作流列表中查找匹配的工作流
    const foundInPersonal = personalWorkflows.value.find(wf => wf.id === selectedWorkflowId.value)
    const foundInPublic = publicWorkflows.value.find(wf => wf.id === selectedWorkflowId.value)
    selectedWorkflow.value = foundInPersonal || foundInPublic || null
  } else {
    selectedWorkflowId.value = null
    selectedWorkflow.value = null
  }
}, { immediate: true })
</script>

<style scoped>
.v-list-item--active {
  background-color: rgb(var(--v-theme-primary-lighten-5));
}

:deep(.v-window__container) {
  max-height: 300px;
  overflow-y: auto;
}
</style> 