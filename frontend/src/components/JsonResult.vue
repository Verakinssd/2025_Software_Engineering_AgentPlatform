<template>
	<div
		v-for="(item, index) in data"
	>
		<v-textarea
			:key="index"
			:label="`${item.value_name}(${item.value_type})`"
			readonly
			v-model="item.value"
			variant="solo"
			rows="1"
      auto-grow
			style="font-family: 'JetBrains Mono'"
		></v-textarea>
	</div>
</template>
<script setup>
const props = defineProps({
	data: {
		required: true,
	}
})
const data = JSON.parse(props.data)
</script>