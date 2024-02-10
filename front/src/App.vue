<template>
  <component :is="comp" :query="query" />
</template>

<script setup lang="ts">
import { ref, shallowRef, watch, defineAsyncComponent, Component } from 'vue';
import qs from 'qs';
import none from '@/comps/None.vue';
import pdf from '@/comps/Pdf.vue';
import img from '@/comps/Img.vue';
import docx from '@/comps/Docx.vue';
import xlsx from '@/comps/Xlsx.vue';
import html from '@/comps/Html.vue';
const comp = shallowRef(none);
const query = ref('');

const init = () => {
  query.value = window.location.href.split('#')[0].split('?')[1];
  let q = qs.parse(query.value);
  switch (q.t) {
  case 'dsp':
  case 'ppt':
  case 'pptx':
  case 'pdf':
    comp.value = pdf;
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
  case 'txt':
  case 'json':
  case 'xml':
  case 'toml':
  case 'yaml':
  case 'config':
  case 'wps':
  case 'docx':
    comp.value = docx;
    break;
  case 'et':
  case 'csv':
  case 'xlsx':
    comp.value = xlsx;
    break;
  case 'dwg':
  case 'dxf':
    comp.value = img;
    break;
  case 'pages':
  case 'numbers':
  case 'key':
    comp.value = html;
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
