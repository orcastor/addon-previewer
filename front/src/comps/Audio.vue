<template v-loading="isLoading" class="app">
  <vue-plyr>
    <audio  controls crossorigin playsinline>
    <source
        :src="source"
        type="audio/mp3"
      />
    </audio >
  </vue-plyr>
</template>

<script setup lang="ts">
import { ref, onMounted, defineProps } from 'vue';
import { store } from "@/store";

const props = defineProps<{
    query?:string
}>()
const isLoading = ref(true);
const source = ref('');

onMounted(() => {
  source.value = `//${location.host}/prvw/api/get?${props.query}&token=${store.token}`;
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
</style>