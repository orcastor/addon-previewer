<template>
  <div class="image__preview">
    <el-image
      style="position:absolute; width:100%; height:100%; top:0; left:0;"
      :src="url"
      :preview-src-list="srcList"
      :initial-index="0"
      fit="contain"
    />
  </div>
</template>

<script lang="ts" setup>
import { defineProps } from 'vue';
import { store } from "@/store";

const props = defineProps<{
    query?:string
}>()

const url = '//' + location.host + '/prvw/api/get?' + props.query + '&token=' + store.token;
const srcList = [ url ]
</script>

<style scoped>
.image__error .image-slot {
  font-size: 30px;
}
.image__error .image-slot .el-icon {
  font-size: 30px;
}
.image__error .el-image {
  width: 100%;
  height: 200px;
}
</style>