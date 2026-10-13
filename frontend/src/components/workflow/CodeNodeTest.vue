<template>
	<v-card class="mt-4 elevation-0">
    <div v-if="isRunning" class="d-flex align-center justify-center">
      <v-progress-circular
        :size="80"
        :width="8"
        color = '#26C6DA'
        class="mt-2"
        indeterminate
      >
      </v-progress-circular>
    </div>

    <v-card-title v-if="!isRunning" class="d-flex align-center" >
      <v-btn icon flat dense :ripple="false" @click="showForm = !showForm" class="mr-2">
        <v-icon>{{ showForm ? 'mdi-chevron-down' :'mdi-chevron-right' }}</v-icon> 
      </v-btn>

      <span class="text-subtitle-1">运行结果</span>
      <v-tooltip :attach="true" text="运行结果">
        <template #activator="{ props }">
          <v-icon v-bind="props" size="20" class="ms-2">
            mdi-information-outline
          </v-icon>
        </template>
      </v-tooltip>
    </v-card-title>

    <v-expand-transition v-if="showForm && !isRunning">
      <v-expansion-panels 
        v-model="panel"
        multiple
        elevation = "0"
      >
        <v-expansion-panel>
          <v-expansion-panel-title style="padding : 0px 20px; min-height: 40px;">输入</v-expansion-panel-title>
          <v-expansion-panel-text>
            <JsonResult
              v-if="!isRunning"
              :data="inputJson"
              language="json"
              theme="vs-dark"
              :height="170"
              :options="{ 
                automaticLayout: true, 
                minimap: { enabled: false },
                readOnly: true
              }"
            />
          </v-expansion-panel-text>
        </v-expansion-panel>

        <v-expansion-panel>
          <v-expansion-panel-title style="padding : 0px 20px; min-height: 40px;">输出</v-expansion-panel-title>
          <v-expansion-panel-text>
            <JsonResult
              v-if="!isRunning"
              :data="outputJson"
              language="json"
              theme="vs-dark"
              :height="170"
              :options="{ 
                automaticLayout: true, 
                minimap: { enabled: false },
                readOnly: true
              }"
            />
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </v-expand-transition>  
  </v-card>
</template>

<script setup lang="ts">
import { ref, watch , computed } from 'vue'
import { selectedNode , showCodeTest , nodes, edges,
   saveWorkflowContent , runWorkflow } from '@/api/workflow'
import ts from 'typescript'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import MonacoEditor from 'monaco-editor-vue3'
import JsonResult from '../JsonResult.vue'

const route = useRoute()
const { agent_id } = route.params
const isRunning = ref(false)
const panel = ref([0, 1])

const inputJson = computed(() => 
  JSON.stringify(selectedNode.value?.data.input ?? [], null, 2)
)

const outputJson = computed(() => 
  JSON.stringify(selectedNode.value?.data.output ?? [], null, 2)
)

async function runCodeNode() {
  isRunning.value = true
  console.log('runJavaScriptCode' , selectedNode.value)

  const node = selectedNode.value

  const names  = node.data.input?.map(i => i.value_name) ?? []
  const values = node.data.input?.map(i => i.value) ?? []

  let mainFn: (...args: any[]) => Promise<any>
  try {
    const jsCode = ts.transpileModule(node.data.codeJavaScript, {
      compilerOptions: { module: ts.ModuleKind.ESNext , target: ts.ScriptTarget.ES2020 } 
    }).outputText
    const factory = new Function(
      ...names,
      jsCode + '\nreturn main;'
    )
    mainFn = factory(...values) 

    if (typeof mainFn !== 'function') {
      ElMessage.error('未能从代码中获取到名为 main 的函数')
      throw new Error('未能从代码中获取到名为 main 的函数')
    }
  } catch (e) {
    ElMessage.error('函数编译错误，请检查代码和输入')
    console.error('函数编译错误:', e)
    isRunning.value = false
    return
  }

  let result: any
  try {
    const params: Record<string, any> = {}
    node.data.input?.forEach(i => {
      params[i.value_name] = i.value
    })
    result = await mainFn({ params })  
  } catch (e) {
    ElMessage.error('函数执行错误，请检查代码')
    console.error('函数执行错误:', e)
    isRunning.value = false
    return
  }

  for (const [key , value] of Object.entries(result)) {
    const output = node.data.output?.find(o => o.value_name === key)
    if (output) {
      output.value = value.toString();
    }
  } 

  isRunning.value = false
}

const saveWorkflow = async () => {
  const content = {
    content: {
      nodes: nodes.value,
      edges: edges.value,
    }
  }
  const res = await saveWorkflowContent(Number(agent_id), content)
  if (res.status === 200) {
    ElMessage.success('保存成功')
  } else {
    ElMessage.error('保存失败')
  }
}

async function runPythonCode() {
	isRunning.value = true;
  console.log('runPythonCode' , selectedNode.value)
	await saveWorkflow();
	const res = await runWorkflow(Number(agent_id), selectedNode.value.id);
	if (res.status === 200) {
    const { workflow } = res.data
    const { content } = workflow
		nodes.value = content.nodes;
		selectedNode.value = nodes.value.find(node => node.id == selectedNode.value.id)
    ElMessage.success('运行成功')
  } else {
    ElMessage.error('运行失败')
  }
	isRunning.value = false;
}

watch(showCodeTest ,
  (val) => {
    if (val) {
      if (selectedNode.value.data.codeType === 'JavaScript') {
        runCodeNode()
      } else {
        runPythonCode()
      }
    }
  },
  { immediate: showCodeTest.value }
)

const showForm = ref(true)

</script>

<style scoped>
:deep(.monaco-editor) {
  border-radius: 8px;
  overflow: hidden;
}
</style>