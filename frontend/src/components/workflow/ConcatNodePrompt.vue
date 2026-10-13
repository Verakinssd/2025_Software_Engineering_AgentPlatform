<template>
	<v-card-title class="d-flex align-center">
		<v-btn icon flat dense :ripple="false" @click="showConcatPrompt = !showConcatPrompt" class="mr-2">
			<v-icon>{{ showConcatPrompt ? 'mdi-chevron-down' :'mdi-chevron-right' }}</v-icon> 
		</v-btn>
		<span class="text-subtitle-1">拼接规则</span>
		<v-tooltip :attach="true" location="top" text="可以使用{{变量名}}的方式引用输入参数中的变量">
			<template #activator="{ props }">
				<v-icon v-bind="props" size="20" class="ms-2">
					mdi-information-outline
				</v-icon>
			</template>
		</v-tooltip>
		<v-spacer />
	</v-card-title>

	<v-expand-transition v-if="showConcatPrompt"> 
		<v-container fluid style="padding: 0 16px;">
			<v-textarea
				v-model="selectedNode.data.concatPrompt"
				clear-icon="mdi-close-circle"
				label="拼接规则"
				variant="outlined"
				:readonly="readonly"
				hide-details="true"
			></v-textarea>
		</v-container>
	</v-expand-transition>
</template>

<script setup>
import { ref } from 'vue'
import { selectedNode } from '@/api/workflow'

defineProps({
	readonly: {
		type: Boolean,
		default: false
	}
})

const showConcatPrompt = ref(true)

</script>

<style scoped>
:deep(.v-btn--icon.v-btn--density-default) {
	width: calc(var(--v-btn-height) + 2px);
	height: calc(var(--v-btn-height) + 2px);
}
</style>