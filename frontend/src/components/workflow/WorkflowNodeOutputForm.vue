<template>
	<v-card-title class="d-flex align-center">
		<v-btn icon flat dense :ripple="false" @click="showForm_output = !showForm_output" class="mr-2">
			<v-icon>{{ showForm_output ? 'mdi-chevron-down' :'mdi-chevron-right' }}</v-icon>
		</v-btn>
		<span class="text-subtitle-1">输出</span>
		<v-tooltip :attach="true" text="输出">
			<template #activator="{ props }">
				<v-icon v-bind="props" size="20" class="ms-2">
					mdi-information-outline
				</v-icon>
			</template>
		</v-tooltip>
	</v-card-title>

	<v-expand-transition v-if="showForm_output">
		<v-card-text>
			<v-table density="compact">
				<thead>
					<tr>
						<th class="text-left">变量名</th>
						<th class="text-left">变量类型</th>
						<th class="text-left" v-if="!readonly">操作</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="(output, index) in outputItems" :key="index">
						<td>
							<v-text-field
								v-model="output.value_name"
								variant="outlined"
								density="compact"
								hide-details
								:readonly="readonly"
								@blur="validateName(output)"
								:rules="nameRules"
							></v-text-field>
						</td>
						<td>
							<v-select
								v-model="output.value_type"
								:items="typeOptions"
								variant="outlined"
								density="compact"
								hide-details
								:disabled="readonly"
								:menu-props="{ attach: 'body', zIndex: 3000 }"
							></v-select>
						</td>
						<td v-if="!readonly">
							<v-btn
								icon
								size="small"
								variant="text"
								color="error"
								@click="removeOutput(index)"
							>
								<v-icon>mdi-delete</v-icon>
							</v-btn>
						</td>
					</tr>
				</tbody>
			</v-table>

			<v-btn
				v-if="!readonly"
				variant="outlined"
				size="small"
				class="mt-2"
				@click="addOutput"
			>
				<v-icon start>mdi-plus</v-icon>
				添加输出
			</v-btn>
		</v-card-text>
	</v-expand-transition>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { InputItem , OutputItem } from '@/api/workflow'
import { selectedNode } from '@/api/workflow'
import { watch } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
	readonly: {
		type: Boolean,
		default: false
	}
})

const outputItems = ref<OutputItem[]>([])
const showForm_output = ref(true)

watch(
	() => selectedNode.value,
	(node) => {
		outputItems.value = node?.data.output || []
	},
	{ immediate: true, deep: true }
)

watch(
	() => selectedNode.value?.data.output,
	(newOutputs) => {
		if (newOutputs) {
			outputItems.value = [...newOutputs]
		}
	},
	{ deep: true }
)

const validateName = (item: any) => {
	const v = item.value_name
	if (!v) {
		ElMessage.error('请输入变量名')
		return
	}
	if (!/^[a-zA-Z]/.test(v)) {
    //ElMessage.error('变量名只能以字母开头')
    return
  }
  if (!/^[a-zA-Z0-9_]+$/.test(v)){
    //ElMessage.error('变量名只能包含字母、数字、下划线')
    return
  }
}

const headers_output = [
	{ text: '变量名', value: 'value_name', width: '40%'},
	{ text: '变量类型', value: 'value_type', width: '45%'},
	{ text: '操作', value: 'actions', sortable: false, width: '15%'},
]

const typeOptions = [
	'string' , 'number' , 'integer', 'boolean'
]

const nameRules = [
	(v: string) => !!v || '请输入变量名',
	(v: string) => /^[a-zA-Z][a-zA-Z0-9_]*$/.test(v) || '变量名只能包含字母、数字、下划线，且以字母开头',
]
</script>

<style scoped>
:deep(.v-input--horizontal .v-input__append) {
  margin-inline-start: 2px;
}

:deep(.my-menu) {
  box-shadow: 0px 2px 12px rgba(0,0,0,0.1) !important;
}

:deep(.hover-item:hover) {
  background-color: rgba(0,0,0,0.04);
}

:deep(.v-list-item--density-default.v-list-item--one-line) {
  min-height: 20px;
}

:deep(.v-field) {
  --v-field-padding-start: 10px;
}

:deep(.v-btn--icon.v-btn--density-default) {
  width: calc(var(--v-btn-height) + 2px);
  height: calc(var(--v-btn-height) + 2px);
}

:deep(.fixed-width-table) .v-table__wrapper table {
  table-layout: fixed;
}

:deep(.input-field) .v-field__outline {
  border: 1px solid rgba(0, 0, 0, 0.38) !important;
  border-radius: 4px;
}

:deep(.input-field:hover) .v-field__outline {
  border: 1px solid rgba(0, 0, 0, 0.87) !important;
}

:deep(
  .v-data-table table th,
  .v-data-table table td,
  .v-data-table table tr) {
  border: none !important;
}
</style>
