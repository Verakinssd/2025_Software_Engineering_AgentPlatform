<script lang="ts" setup>
import { Background } from '@vue-flow/background'
</script>

<template>
  <div class="dropzone-background">
    <Background :size="2" :gap="20" pattern-color="#BDBDBD" />
  </div>
</template>


<style>
.dropzone-background {
    position: relative;
    height: 100%;
    width: 100%;
}
</style>