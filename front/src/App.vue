<template>
  <component :is="comp"></component>
</template>

<script setup lang="ts">
import { ref, shallowRef, watch, defineAsyncComponent, Component } from 'vue';
import qs from 'qs';
import none from '@/comps/None.vue';
import doc from '@/comps/Doc.vue';
const comp = shallowRef(none);

const init = () => {
  let query = qs.parse(window.location.href.split('#')[0].split('?')[1]);
  console.log(query);
  switch (query.t) {
  case 'pdf':
    comp.value = doc;
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
