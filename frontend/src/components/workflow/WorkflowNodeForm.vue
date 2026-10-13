<template>
  <v-card class="mt-4 elevation-0">
    <WorkflowNodeInputForm :readonly="readonly" />
    <WorkflowNodeSelector :readonly="readonly" />
    <WorkflowNodeOutputForm :readonly="readonly" />
  </v-card>
</template>

<script setup>
defineProps({
  readonly: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
:deep(.v-table > .v-table__wrapper > table > tbody > tr > td){
  padding: 0 3px;
}
</style> 