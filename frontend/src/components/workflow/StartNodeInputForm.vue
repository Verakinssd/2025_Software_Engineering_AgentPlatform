<template>
	<v-card-title class="d-flex align-center">
		<v-btn icon flat dense :ripple="false" @click="showForm_input = !showForm_input" class="mr-2">
			<v-icon>{{ showForm_input ? 'mdi-chevron-down' : 'mdi-chevron-right' }}</v-icon>
		</v-btn>
		<span class="text-subtitle-1">输入</span>
	</v-card-title>

	<v-expand-transition v-if="showForm_input">
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
					<tr v-for="(input, index) in inputItems" :key="index">
						<td>
							<v-text-field
								v-model="input.value_name"
								variant="outlined"
								density="compact"
								hide-details
								:readonly="readonly"
								@blur="validateName(input)"
								:rules="nameRules"
							></v-text-field>
						</td>
						<td>
							<v-select
								v-model="input.value_type"
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
								@click="removeInput(index)"
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
				@click="addInput"
			>
				<v-icon start>mdi-plus</v-icon>
				添加输入
			</v-btn>
		</v-card-text>
	</v-expand-transition>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { selectedNode } from '@/api/workflow'

defineProps({
	readonly: {
		type: Boolean,
		default: false
	}
})

// ... existing code ...
</script> 