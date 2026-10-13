<template>
	<v-card-title class="d-flex align-center">
		<v-btn icon flat dense :ripple="false" @click="showForm_query = !showForm_query" class="mr-2">
			<v-icon>{{ showForm_query ? 'mdi-chevron-down' : 'mdi-chevron-right' }}</v-icon>
		</v-btn>
		<span class="text-subtitle-1">查询模板</span>
		<v-tooltip text="设置查询模板，可以使用{{变量名}}引用输入" :attach="true">
			<template #activator="{ props }">
				<v-icon v-bind="props" size="20" class="ms-2">
					mdi-information-outline
				</v-icon>
			</template>
		</v-tooltip>
	</v-card-title>

	<v-expand-transition v-if="showForm_query">
		<v-card-text>
			<v-textarea
				v-model="queryTemplate"
				label="查询模板"
				placeholder="例如：请帮我查询有关{{query}}的信息"
				hint="使用{{变量名}}引用输入变量"
				rows="5"
				auto-grow
				variant="outlined"
				density="compact"
				:readonly="readonly"
				@update:model-value="updateQueryTemplate"
			></v-textarea>
		</v-card-text>
	</v-expand-transition>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { selectedNode } from '@/api/workflow'

const showForm_query = ref(true)
const queryTemplate = ref('')
const loading = ref(false)

defineProps({
	readonly: {
		type: Boolean,
		default: false
	}
})

onMounted(async () => {
  // 先获取知识库文件列表
	loading.value = true;
	loading.value = false;

  if (selectedNode.value) {
    if (selectedNode.value.data.queryTemplate) {
      queryTemplate.value = selectedNode.value.data.queryTemplate
    } else {
      queryTemplate.value = '请回答关于{{query}}的问题'
    }
  }
})

function updateQueryTemplate() {
  if (selectedNode.value) {
    selectedNode.value.data.queryTemplate = queryTemplate.value
  }
}
</script>

<style scoped>
:deep(.v-btn--icon.v-btn--density-default) {
  width: calc(var(--v-btn-height) + 2px);
  height: calc(var(--v-btn-height) + 2px);
}
</style>