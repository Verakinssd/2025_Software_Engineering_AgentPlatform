<template>
	<v-card-title class="d-flex align-center">
		<v-btn icon flat dense :ripple="false" @click="showCodeEditor = !showCodeEditor" class="mr-2">
			<v-icon>{{ showCodeEditor ? 'mdi-chevron-down' :'mdi-chevron-right' }}</v-icon> 
		</v-btn>
		<span class="text-subtitle-1">代码</span>
		<v-tooltip :attach="true" text="代码">
			<template #activator="{ props }">
				<v-icon v-bind="props" size="20" class="ms-2">
					mdi-information-outline
				</v-icon>
			</template>
		</v-tooltip>
		<v-switch
			v-model="selectedNode.data.codeType"
			:label="selectedNode.data.codeType"
			false-value="JavaScript"
			true-value="Python"
			hide-details
			class="ml-5"
			color="secondary"
			baseColor="secondary"
			:disabled="readonly"
		>
		</v-switch>
	</v-card-title>

	<v-expand-transition v-if="showCodeEditor"> 
		<v-card-text class="pa-0" style="height: 200px;">
  		<MonacoEditor
				v-if="selectedNode.data.codeType === 'JavaScript'"
				v-model:value="selectedNode.data.codeJavaScript"
				language="javascript"
				theme="vs-dark"
				:height="200"
				:options="{ 
					automaticLayout: true , 
					minimap: { enabled: false },
					readOnly: readonly
				}"
			></MonacoEditor>
			<MonacoEditor
				v-if="selectedNode.data.codeType === 'Python'"
				v-model:value="selectedNode.data.codePython"
				language="python"
				theme="vs-dark"
				:height="200"
				:options="{ 
					automaticLayout: true , 
					minimap: { enabled: false },
					readOnly: readonly
				}"
			></MonacoEditor>
		</v-card-text>
	</v-expand-transition>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { selectedNode } from '@/api/workflow'
import MonacoEditor from 'monaco-editor-vue3'

defineProps({
	readonly: {
		type: Boolean,
		default: false
	}
})

const showCodeEditor = ref(true)

</script>

<style scoped>
:deep(.v-btn--icon.v-btn--density-default) {
  width: calc(var(--v-btn-height) + 2px);
  height: calc(var(--v-btn-height) + 2px);
}

:deep(.monaco-editor) {
  border-radius: 8px;
  overflow: hidden;
}
</style>