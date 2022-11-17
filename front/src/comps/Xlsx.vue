<template v-loading="isLoading" class="app">
  <div class="app-content" id="result" />
</template>

<script setup lang="ts">
import xlsxPreview from 'xlsx-preview';
import { ref, onMounted, defineProps } from 'vue';
import { store } from "@/store";
import axios from "axios";
import qs from 'qs';

const props = defineProps<{
    query?:string
}>()
const isLoading = ref(true);

const load = async () => {
  let q = qs.parse(props.query);
  const res = await axios({
    method: 'get',
    url: '//' + location.host + '/prvw/api/get?' + props.query + '&token=' + store.token,
    responseType: 'arraybuffer',
  });
  const result = await xlsxPreview.xlsx2Html(res.data, { output: 'arraybuffer', format: q.t });
  const url = URL.createObjectURL(new Blob([result], {
    type: 'text/html',
  }));
  document.querySelector('#result').innerHTML =
    `<object class="res-obj" type="text/html" data="${url}"></object>`;
  isLoading.value = false;
}

onMounted(() => {
  load();
});
</script>

<style>
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

.res-obj {
  width: 100%;
  height: 85vh;
  background-color: rgba(255, 255, 255, 0.9);
}
</style>