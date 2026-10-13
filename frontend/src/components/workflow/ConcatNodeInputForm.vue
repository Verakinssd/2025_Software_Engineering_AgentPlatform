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
		<v-spacer />
		<v-btn v-if="showForm_input && !readonly" icon @click="addRow_input" flat>
			<v-icon color="primary">mdi-plus</v-icon>
		</v-btn>
	</v-card-title>

	<v-expand-transition v-if="showForm_input">
		<v-data-table
			:headers="headers_input"
			:items="inputItems"
			hide-default-header
			hide-default-footer
			disable-pagination
			no-data-text="暂无输入"
		>
			<template #item.value_name="{ item }">
				<v-text-field
					v-model="item.value_name"
					variant="outlined"
					density="compact"
					label="变量名"
					label-position="top"
					hide-details
					:readonly="readonly"
					@blur = "validateName(item)"
					style="min-height: 20px"
				/>
			</template>
			<template #item.value_type="{ item }">
				<v-select
					v-model="item.value_type"
					:items="typeOptions"
					variant="outlined"
					density="compact"
					label="变量类型"
					label-position="top"
					:readonly="readonly"
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
						<v-btn v-if="!readonly" icon dense flat class="btn-add">
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
					:readonly="readonly"
					:menu-props="{ attach: 'body', zIndex: 3000 }"
					hide-details
				>
					<template #append>
						<v-btn v-if="!readonly" icon dense flat class="btn-add">
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
			<template #item.actions="{ index }">
				<v-icon v-if="!readonly" small class="mr-2" @click="removeRow_input(index)" style="cursor: pointer">
					mdi-minus
				</v-icon>
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
import knowledgeIcon from '@/assets/knowledge.png'
import workflowIcon from '@/assets/workflow.png'
import type { InputItem , OutputItem } from '@/api/workflow'
import { selectedNode,simplifyType, getAttributes ,nodes, edges } from '@/api/workflow'
import { watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useVueFlow } from "@vue-flow/core";

const inputItems = ref<InputItem[]>([])

const showForm_input = ref(true)

const boolOptions = ['true', 'false']

function getIcon(type: any) {
  return type === 'llm' ? agentIcon :
  type === 'myinput' ? inputIcon :
  type === 'myoutput' ? outputIcon :
  type === 'select' ? selectIcon :
  type === 'end' ? endIcon :
  type === 'start' ? startIcon :
  type === 'code' ? codeIcon :
	type === 'concat' ? concatIcon :
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
	item.value = attr.value
	item.assigned = true
	item.assignedNodeId = node.id
	item.assignedNodeType = node.type
	item.assignedNodeName = node.data.name
	item.assignedAttributeName = attr.value_name;
	item.assignedAttributeType = attr.value_type;
}

function clearAssignment (item: any) {
	item.assignedNodeId = null
	item.assignedNodeType = ''
	item.assignedNodeName = ''
	item.assignedAttributeName = ''
	item.assigned = false
	item.value = ''
}

function getPredecessors(nodeId : string) {
	const incoming = edges.value.filter(edge => edge.target === nodeId)
	const sourceIds = incoming.map(edge => edge.source)
	const sourceNodes = nodes.value.filter(node => sourceIds.includes(node.id))
	console.log('sourceNodes' , sourceNodes)
	return sourceNodes
}

const predecessorNodes = computed(() => {
	if (selectedNode.value) {
		return getPredecessors(selectedNode.value.id)
	}
	return []
})

const nodeAttributes = computed(() => {
	const map = new Map<string, ReturnType<typeof getAttributes>>()
	predecessorNodes.value.forEach(node => {
		map.set(node.id, getAttributes(node))
	})
	return map
})

function validateValue(item: any) {
	if (item.assigned) {
		return;
	}
	const v = item.value
	if (!v) {
		ElMessage.error('请输入变量值')
		return
	}
	if (item.value_type === 'integer') {
		if (!/^\d+$/.test(v)) {
			ElMessage.error('请输入整数')
			return
		}
	} else if (item.value_type === 'number') {
		if (!/^\d+(\.\d+)?$/.test(v)) {
			ElMessage.error('请输入数字')
			return
		}
	}
}

watch(
	() => selectedNode.value,
	(node) => {
		inputItems.value = node?.data.input || []
	},
	{ immediate: true }
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

const headers_input = [
	{ text: '变量名', value: 'value_name', width: '25%'},
	{ text: '变量类型', value: 'value_type', width: '23%'},
	{ text: '变量值', value: 'value',width: '47%'},
	{ text: '操作', value: 'actions', sortable: false, width: '5%'},
]

const typeOptions = [
	{title : 'str' , value : 'string' },
	{title : 'int' , value : 'integer' },
	{title : 'num' , value : 'number' },
	{title : 'bool' , value : 'boolean' }
]

function addRow_input() {
	inputItems.value.push({ value_name: '', value_type: 'string', value: '' })
}

function removeRow_input(i: number) {
	inputItems.value.splice(i, 1)
}

// 定义props
const props = defineProps({
	readonly: {
		type: Boolean,
		default: false
	}
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
