<template v-loading="isLoading" class="app">
  <div id="luckysheet" />
</template>

<script setup lang="ts">
import LuckyExcel from 'luckyexcel';
import { ref, onMounted, defineProps } from 'vue';
import { store } from "@/store";

import 'luckysheet/dist/css/luckysheet.css';
import 'luckysheet/dist/plugins/plugins.css';
import 'luckysheet/dist/plugins/css/pluginsCss.css';
import 'luckysheet/dist/assets/iconfont/iconfont.css';

function getType(val) {
  return Object.prototype.toString.call(val).slice(8, -1)
}

function isFunction(val) {
  return getType(val) === 'Function'
}

const props = defineProps<{
    query?:string
}>()
const isLoading = ref(true);

const load = async () => {
  const url = `//${location.host}/prvw/api/get?${props.query}&token=${store.token}`;

  LuckyExcel.transformExcelToLuckyByUrl(url, '', (exportJson:any, _luckysheetfile:any) => {
    if (exportJson.sheets == null || exportJson.sheets.length == 0) {
      alert('Failed to read the content of the excel file, currently does not support xls files!')
      return
    }

    isFunction(window?.luckysheet?.destroy) && window.luckysheet.destroy()

    window.luckysheet.create({
      container: 'luckysheet', //luckysheet is the container id
      showtoolbar: false,
      showinfobar: false,
      showstatisticBar: false,
      showstatisticBarConfig:{
        zoom: true,
      },
      data: exportJson.sheets,
      title: exportJson.info.name,
      userInfo: exportJson.info.name.creator,
    })
  });
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

#luckysheet {
  margin: 0px;
  padding: 0px;
  position: absolute;
  width: 100%;
  left: 0px;
  top: 0px;
  bottom: 0px;
}
</style>