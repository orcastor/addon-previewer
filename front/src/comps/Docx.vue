<template v-loading="isLoading" class="app">
  <div class="app-content">
    <div ref="refFile"></div>
  </div>
</template>

<script setup lang="ts">
import docx from 'docx-preview';
import { ref, onMounted, defineProps } from 'vue';
import { store } from "@/store";
import axios from "axios";

const props = defineProps<{
    query?:string
}>()
const isLoading = ref(true);
const refFile = ref();

onMounted(() => {
  load();
});

const load = async () => {
  const res = await axios({
    method: 'get',
    url: '//' + location.host + '/prvw/api/get?' + props.query + '&token=' + store.token,
    responseType: 'blob',
  })
  docx.renderAsync(res.data, refFile.value);
  isLoading.value = false;
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