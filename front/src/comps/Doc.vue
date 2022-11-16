<template v-loading="isLoading" class="app">
  <div class="app-content">
    <vue-pdf-embed :source="pdfSource" @password-requested="passwordPrompt" @rendered="rendered" />
  </div>
</template>

<script setup lang="ts">
import VuePdfEmbed from 'vue-pdf-embed';
import { ref, onMounted, defineProps } from 'vue';
import { store } from "@/store";

const props = defineProps<{
    query?:string
}>()
const isLoading = ref(true);
const pdfSource = ref('');

onMounted(() => {
  pdfSource.value = `//${location.host}/prvw/api/get?${props.query}&token=${store.token}`;
});

const rendered = ()=> {
  isLoading.value = false;
};

const passwordPrompt = (callback:any, retry:any) => {
  callback(prompt(retry ? 'Enter password again' : 'Enter password'));
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
</style>