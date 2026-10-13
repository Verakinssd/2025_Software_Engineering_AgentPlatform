<template>
	<v-card class="elevation-0">
		<v-card-title class="d-flex align-center">
			<v-btn icon flat dense :ripple="false" @click="showForm_input = !showForm_input" class="mr-2">
				<v-icon>{{ showForm_input ? 'mdi-chevron-down' :'mdi-chevron-right' }}</v-icon> 
			</v-btn>
			<span class="text-subtitle-1">试运行输入</span>
			<v-tooltip :attach="true" text="试运行输入">
				<template #activator="{ props }">
					<v-icon v-bind="props" size="20" class="ms-2">
						mdi-information-outline
					</v-icon>
				</template>
			</v-tooltip>
		</v-card-title>

		<v-expand-transition v-if="showForm_input">
			<v-sheet v-if="finishLoading">
				<v-list density="compact">
					<v-list-item v-for="inputItem in startNode.data.input">
						<div class="text-subtitle-1 text-medium-emphasis">
							<strong class="font-weight-bold">{{ inputItem.value_name }}</strong> 
							{{ simplifyType(inputItem.value_type) }}
						</div>
						<v-text-field
							density="compact"
							v-model="inputItem.value"
							variant="outlined"
							color="primary"
							clearable
							hide-details
							:placeholder = "'请输入' + inputItem.value_name "
						></v-text-field>
					</v-list-item>
				</v-list>

				<div class="d-flex flex-column">
					<v-btn
						class="mt-1"
						color="success"
						@click="runWholeWorkflow"
					>
						试运行
					</v-btn>
				</div>
			</v-sheet>
		</v-expand-transition>

		<v-divider></v-divider>


		<v-card-text v-if="running" class="text-center pt-6">
			<v-progress-circular 
				indeterminate 
				color="#26C6DA"
				class="mt-2"
				:size="60"
				:width="6"
				>
			</v-progress-circular>
			<div class="mt-3">试运行中...</div>
		</v-card-text>

		<template v-if="finishRunning">
			<v-card-title class="d-flex align-center">
				<v-btn icon flat dense :ripple="false" @click="showForm_output = !showForm_output" class="mr-2">
					<v-icon>{{ showForm_output ? 'mdi-chevron-down' :'mdi-chevron-right' }}</v-icon> 
				</v-btn>
				<span class="text-subtitle-1">试运行输出</span>
				<v-tooltip :attach="true" text="试运行输出">
					<template #activator="{ props }">
						<v-icon v-bind="props" size="20" class="ms-2">
							mdi-information-outline
						</v-icon>
					</template>
				</v-tooltip>
			</v-card-title>

			<v-expand-transition v-if="showForm_output">
				<v-list density="compact">
					<v-list-item v-for="outputItem in endNode.data.output">
						<div class="text-subtitle-1 text-medium-emphasis">
							<strong class="font-weight-bold">{{ outputItem.value_name }}</strong> 
							{{ simplifyType(outputItem.value_type) }}
						</div>
						<v-textarea
							:model-value=outputItem.value
							variant="solo"
							color="primary"
							readonly
							hide-details
							rows="1"
							auto-grow
							style="font-family: 'JetBrains Mono'"
						></v-textarea>
					</v-list-item>
				</v-list>
			</v-expand-transition>
		</template>
	</v-card>
</template>

<script setup>
import { ref , onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { nodes , edges, runWorkflow, saveWorkflowContent } from '@/api/workflow'
import { ElMessage } from 'element-plus'

const route = useRoute()
const showForm_input = ref(true)
const showForm_output = ref(true)
const startNode = ref()
const endNode = ref()
const finishLoading = ref(false)
const { agent_id } = route.params 
const running = ref(false)
const finishRunning = ref(false)

onMounted(() => {
	startNode.value = nodes.value.find(node => node.id == 1)
	endNode.value = nodes.value.find(node => node.id == 2)
	finishLoading.value = true
})

const saveWorkflow = async () => {
  const content = {
    content: {
      nodes: nodes.value,
      edges: edges.value,
    }
  }
  const res = await saveWorkflowContent(agent_id, content)
  if (res.status === 200) {
    ElMessage.success('保存成功')
  } else {
    ElMessage.error('保存失败')
  }
}

async function runWholeWorkflow() {
	await saveWorkflow()
	running.value = true;
	finishRunning.value = false;
  const res = await runWorkflow(agent_id)
  if (res.status === 200) {
    const { workflow } = res.data
    const { content } = workflow
    nodes.value = content.nodes
    edges.value = content.edges || []
    ElMessage.success('运行成功')
  } else {
    ElMessage.error('运行失败')
  }
  running.value = false
	finishRunning.value = true
	console.log('endNode' , endNode.value.data.output)
}

function simplifyType(type) {
  if (type == 'string') return 'str';
  if (type == 'number') return 'num';
  if (type == 'boolean') return 'bool';
  if (type == 'integer') return 'int';
}
</script>

<style scoped>
</style>