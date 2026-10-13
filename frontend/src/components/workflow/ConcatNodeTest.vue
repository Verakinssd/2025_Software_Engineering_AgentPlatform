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
							:height="150"
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
							:height="150"
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
  
<script setup>
import { ref, watch , computed } from 'vue'
import { useRoute } from 'vue-router'
import { selectedNode , showConcatTest , nodes, edges 
		,saveWorkflowContent ,runWorkflow } from '@/api/workflow'
import { ElMessage } from 'element-plus'
import MonacoEditor from 'monaco-editor-vue3'
import JsonResult from '../JsonResult.vue'

const isRunning = ref(false)
const showForm = ref(true)
const panel = ref([0, 1])
const route = useRoute()
const { agent_id } = route.params

const inputJson = computed(() => 
	JSON.stringify(selectedNode.value?.data.input ?? [], null, 2)
)

const outputJson = computed(() => 
	JSON.stringify(selectedNode.value?.data.output ?? [], null, 2)
)

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

async function runConcatNode() {
	isRunning.value = true;
	console.log('ConcatNodeRun', selectedNode.value)
	await saveWorkflow();
	const res = await runWorkflow(Number(agent_id), selectedNode.value.id);
	if (res.status === 200) {
		const { workflow } = res.data
		const { content } = workflow
		nodes.value = content.nodes;
		selectedNode.value = nodes.value.find(node => node.id == selectedNode.value.id)
		ElMessage.success('运行成功')
	} else {
		ElMessage.error(res.value.data.message)
	}
	isRunning.value = false;
}

watch(showConcatTest ,
	(val) => {
		if (val) {
			runConcatNode();
		}
	},
	{ immediate: showConcatTest.value }
)   
</script>

<style scoped>
:deep(.monaco-editor) {
	border-radius: 8px;
	overflow: hidden;
}
</style>