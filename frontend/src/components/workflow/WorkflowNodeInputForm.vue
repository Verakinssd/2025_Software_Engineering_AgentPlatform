<template>
	<v-card-title class="d-flex align-center">
		<v-btn icon flat dense :ripple="false" @click="showForm_input = !showForm_input" class="mr-2">
			<v-icon>{{ showForm_input ? 'mdi-chevron-down' :'mdi-chevron-right' }}</v-icon>
		</v-btn>
		<span class="text-subtitle-1">输入</span>
		<v-tooltip :attach="true" text="输入">
			<template #activator="{ props }">
				<v-icon v-bind="props" size="20" class="ms-2">
					mdi-information-outline
				</v-icon>
			</template>
		</v-tooltip>
	</v-card-title>

	<v-expand-transition v-if="showForm_input">
		<v-data-table
			:headers="headers_input"
			:items="inputItems"
			hide-default-header
			hide-default-footer
			disable-pagination
			no-data-text="请选择工作流"
			:readonly="readonly"
		>
			<template #item.value_name="{ item }">
				<v-text-field
					readonly
					v-model="item.value_name"
					variant="outlined"
					density="compact"
					label="变量名"
					label-position="top"
					hide-details
					@blur = "validateName(item)"
					style="min-height: 20px"
				/>
			</template>
			<template #item.value_type="{ item }">
				<v-select
				  readonly
					v-model="item.value_type"
					:items="typeOptions"
					variant="outlined"
					density="compact"
					label="变量类型"
					label-position="top"
					@blur = "checkType(item)"
					:menu-props="{ attach: 'body', zIndex: 3000 }"
					hide-details
				/>
			</template>
			<template #item.value="{ item }">
				<v-text-field
					v-if="item.value_type != 'boolean' || item.assigned"
					:model-value="displayText(item)"
					@update:model-value="val => {
						if (!item.assigned && !readonly) {
							item.value = val;
						}
					}"
					:readonly="item.assigned || readonly"
					variant="outlined"
					density="compact"
					:label="item.assigned ? '' : '变量值'"
					label-position="top"
					hide-details
					:clearable="!readonly"
					@click:clear="clearAssignment(item)"
					@blur = "validateValue(item)"
					style="min-height: 20px"
					>
					<template #prepend-inner>
						<el-image
							v-if="item.assigned"
							style="width: 25px; height: 25px;"
							:src=getIcon(item.assignedNodeType) />
					</template>
					<template #append>
						<v-btn icon dense flat class="btn-add" :disabled="readonly">
							<v-icon size="18">
								mdi-database-import
							</v-icon>
							<v-menu activator="parent" :attach="true">
								<v-list>
									<template v-if="predecessorNodes.length!=0">
										<v-list-item
											v-for="node in predecessorNodes"
											link
											color="primary"
											>
											<template v-slot:prepend>
												<el-image style="width: 25px; height: 25px;"
													:src=getIcon(node.type) />
											</template>
											<v-list-item-title v-text="node.data.name" class="ml-3"></v-list-item-title>

											<v-menu :open-on-focus="false" activator="parent" open-on-hover submenu>
												<v-list v-if="nodeAttributes.get(node.id)?.length > 0" style="margin-right: 1px">
													<v-list-item
														v-for="(attr, i) in nodeAttributes.get(node.id)"
														:key="i"
														:value="attr"
														color="primary"
														@click="assignAttribute(item , node , attr)"
													>
														<template v-slot:prepend>
															<span><strong>{{simplifyType(attr.value_type) }}.  </strong></span>
														</template>
														<v-list-item-title v-text="attr.value_name"></v-list-item-title>
													</v-list-item>
												</v-list>
											</v-menu>
										</v-list-item>
									</template>
									<template v-else>
										<v-list-item>
											<v-list-item-title>无前继节点</v-list-item-title>
										</v-list-item>
									</template>
								</v-list>
							</v-menu>
						</v-btn>
					</template>
				</v-text-field>
				<v-select
					v-else
					:v-model="item.value"
					:items="boolOptions"
					variant="outlined"
					density="compact"
					label="变量值"
					label-position="top"
					:menu-props="{ attach: 'body', zIndex: 3000 }"
					hide-details
				>
					<template #append>
						<v-btn icon dense flat class="btn-add" :disabled="readonly">
							<v-icon size="18">
								mdi-database-import
							</v-icon>
							<v-menu activator="parent" :attach="true">
								<v-list>
									<template v-if="predecessorNodes.length!=0">
										<v-list-item
											v-for="node in predecessorNodes"
											link
											color="primary"
											>
											<template v-slot:prepend>
												<el-image style="width: 25px; height: 25px;"
													:src=getIcon(node.type) />
											</template>
											<v-list-item-title v-text="node.data.name" class="ml-3"></v-list-item-title>

											<v-menu :open-on-focus="false" activator="parent" open-on-hover submenu>
												<v-list v-if="nodeAttributes.get(node.id)?.length > 0" style="margin-right: 1px">
													<v-list-item
														v-for="(attr, i) in nodeAttributes.get(node.id)"
														:key="i"
														:value="attr"
														color="primary"
														@click="assignAttribute(item , node , attr)"
													>
														<template v-slot:prepend>
															<span><strong>{{simplifyType(attr.value_type) }}.  </strong></span>
														</template>
														<v-list-item-title v-text="attr.value_name"></v-list-item-title>
													</v-list-item>
												</v-list>
											</v-menu>
										</v-list-item>
									</template>
									<template v-else>
										<v-list-item>
											<v-list-item-title>无前继节点</v-list-item-title>
										</v-list-item>
									</template>
								</v-list>
							</v-menu>
						</v-btn>
					</template>
				</v-select>
			</template>
		</v-data-table>
	</v-expand-transition>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import agentIcon from '@/assets/agent.png'
import selectIcon from '@/assets/choice.png'
import startIcon from '@/assets/start.png'
import endIcon from '@/assets/end.png'
import inputIcon from '@/assets/input.png'
import outputIcon from '@/assets/output.png'
import codeIcon from '@/assets/code.png'
import concatIcon from '@/assets/concat.png'
import workflowIcon from '@/assets/workflow.png'
import knowledgeIcon from '@/assets/knowledge.png'
import type { InputItem , OutputItem } from '@/api/workflow'
import { selectedNode,simplifyType, getAttributes } from '@/api/workflow'
import { watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useVueFlow } from "@vue-flow/core";

const { nodes, edges } = useVueFlow();

const inputItems = ref<InputItem[]>([])

const showForm_input = ref(true)

const boolOptions = ['true', 'false']

const props = defineProps({
	readonly: {
		type: Boolean,
		default: false
	}
})

function getIcon(type: any) {
  return type === 'llm' ? agentIcon :
  type === 'myinput' ? inputIcon :
  type === 'myoutput' ? outputIcon :
  type === 'select' ? selectIcon :
  type === 'end' ? endIcon :
  type === 'start' ? startIcon :
  type === 'code' ? codeIcon :
  type === 'concat' ? concatIcon :
  type === 'workflow' ? workflowIcon :
  type === 'knowledge' ? knowledgeIcon :
  type === 'workflow' ? workflowIcon :
  ''
}

function displayText(item: any) {
	if (item.assigned) {
		return `${item.assignedNodeName}.${item.assignedAttributeName}`;
	} else {
		return item.value;
	}
}

function checkType(item: any) {
	if (item.assigned && item.value_type != item.assignedAttributeType) {
		ElMessage.error('类型不匹配')
		return
	}
}

function assignAttribute (item: any, node: any, attr : any) {
	if (item.value_type != attr.value_type) {
		ElMessage.error('类型不匹配')
		return
	}
	item.assigned = true
	item.assignedNodeId = node.id
	item.assignedNodeName = node.data.name
	item.assignedNodeType = node.type
	item.assignedAttributeName = attr.value_name
	item.assignedAttributeType = attr.value_type
}

function clearAssignment(item: any) {
	item.assigned = false
	item.assignedNodeId = ''
	item.assignedNodeName = ''
	item.assignedNodeType = ''
	item.assignedAttributeName = ''
	item.assignedAttributeType = ''
	item.value = ''
}

function validateName(item: any) {
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

function validateValue(item: any) {
	if (item.value_type === 'integer' && !Number.isInteger(Number(item.value))) {
    item.value = Math.round(Number(item.value) || 0);
    ElMessage.warning('整数类型已四舍五入');
  } else if (item.value_type === 'number' && isNaN(Number(item.value))) {
    item.value = Number(item.value) || 0;
    ElMessage.warning('数值类型已自动调整');
  }
}

watch(
	() => selectedNode.value,
	(node) => {
		inputItems.value = node?.data.input || []
	},
	{ immediate: true, deep: true }
)

watch(
	() => selectedNode.value?.data.input,
	(newInputs) => {
		if (newInputs) {
			inputItems.value = [...newInputs]
		}
	},
	{ deep: true }
)

const nodeAttributes = computed(() => {
	const result = new Map()
	for (const node of nodes.value) {
		if (node.id != selectedNode.value?.id) {
			const attrs = getAttributes(node)
			if (attrs && attrs.length > 0) {
				result.set(node.id, attrs)
			}
		}
	}
	return result
})

const headers_input = [
	{ text: '变量名', value: 'value_name', width: '25%'},
	{ text: '变量类型', value: 'value_type', width: '25%'},
	{ text: '变量值', value: 'value', width: '50%' },
]

const typeOptions = [
	'string' , 'number' , 'integer', 'boolean'
]

function addRow_input() {
	inputItems.value.push({ value_name: '', value_type: 'string', value: '' })
}

function removeRow_input(i: number) {
	inputItems.value.splice(i, 1)
}

const predecessorNodes = computed(() => {
	if (!selectedNode.value) return []

	const predecessors = []

	// 找到所有指向当前节点的边
	for (const edge of edges.value) {
		if (edge.target === selectedNode.value.id) {
			// 找到对应的源节点
			for (const node of nodes.value) {
				if (node.id === edge.source) {
					predecessors.push(node)
					break
				}
			}
		}
	}

	return predecessors
})
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
