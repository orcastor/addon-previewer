<template v-loading="isLoading" class="app">
  <iframe :src="link" :style=iframeStyle() :onload="isLoading = false" frameborder="0" />
</template>

<script setup lang="ts">
import { ref, onMounted, defineProps } from 'vue';
import { store } from "@/store";

const props = defineProps<{
    query?:string
}>()
const link = ref('');
const isLoading = ref(true);

onMounted(() => {
  load();
});

const iframeStyle = () => {
  return 'width:100%; height:'+(100-5500/document.body.clientHeight).toFixed(2)+'vh;';
}

const load = async () => {
  link.value = '//' + location.host + '/prvw/api/get?' + props.query + '&token=' + store.token;
};
</script>

<style>
.vue-pdf-embed > div {
  margin-bottom: 8px;
  box-shadow: 0 2px 8px 4px rgba(0, 0, 0, 0.1);
}

.app {
  padding: 16px;
  box-shadow: 0 2px 8px 4px rgba(0, 0, 0, 0.1);
  background-color: #555;
  color: #ddd;
}

.app-content {
  padding: 8px;
  width: 100%;
}

.docx-wrapper {
  background: inherit !important;
}
</style>