<template>
  <component :is="comp" :query="query" />
</template>

<script setup lang="ts">
import { ref, shallowRef, watch, defineAsyncComponent, Component } from 'vue';
import qs from 'qs';
import none from '@/comps/None.vue';
import doc from '@/comps/Doc.vue';
import img from '@/comps/Img.vue';
const comp = shallowRef(none);
const query = ref('');

const init = () => {
  query.value = window.location.href.split('#')[0].split('?')[1];
  let q = qs.parse(query.value);
  switch (q.t) {
  case 'pdf':
    comp.value = doc;
    break;
  case 'jpg':
  case 'jpeg':
  case 'png':
  case 'gif':
  case 'bmp':
  case 'svg':
  case 'eps':
  case 'ai':
  case 'psd':
  case 'ico':
  case 'tiff':
  case 'webp':
    comp.value = img;
    break;
  default:
    comp.value = none;
    break;
  }
};

init();

watch(() => window.location.href, (_newValue:any, _oldValue:any) => {
  init();
});

</script>

<style scoped>
</style>
