<template>
  <component :is="comp"></component>
</template>

<script setup lang="ts">
import { shallowRef } from 'vue';
import pdf from '@/comps/Pdf.vue';
const comp = shallowRef(pdf);

</script>

<style scoped>
</style>
