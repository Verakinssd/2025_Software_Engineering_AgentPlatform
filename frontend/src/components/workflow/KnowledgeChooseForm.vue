<template>
	<v-card-title class="d-flex align-center">
		<v-btn icon flat dense :ripple="false" @click="showRagPrompt = !showRagPrompt" class="mr-2">
				<v-icon>{{ showRagPrompt ? 'mdi-chevron-down' :'mdi-chevron-right' }}</v-icon> 
		</v-btn>
		<span class="text-subtitle-1">查询的文档</span>
		<v-tooltip :attach="true" location="top" text="查询的文档">
			<template #activator="{ props }">
				<v-icon v-bind="props" size="20" class="ms-2">
					mdi-information-outline
				</v-icon>
			</template>
		</v-tooltip>
			<v-spacer />
	</v-card-title>

	<v-expand-transition v-if="showRagPrompt"> 
		<v-container fluid style="padding: 0 16px;">
			<v-textarea
				v-model="selectedNode.data.concatPrompt"
				clear-icon="mdi-close-circle"
				label="查询文档"
				variant="outlined"
				hide-details="true"
			></v-textarea>
		</v-container>
	</v-expand-transition>
</template>
  
<script setup>
import { ref } from 'vue'
import { selectedNode } from '@/api/workflow'
const showRagPrompt = ref(true)

</script>

<style scoped>
:deep(.v-btn--icon.v-btn--density-default) {
	width: calc(var(--v-btn-height) + 2px);
	height: calc(var(--v-btn-height) + 2px);
}
</style>